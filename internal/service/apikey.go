package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/glueops/waggle/internal/models/control"
	"github.com/glueops/waggle/internal/repo"
	"github.com/glueops/waggle/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// apiKeyPrefix marks Waggle organization API keys. It lets the auth middleware
// distinguish an API key from a JWT on the same Authorization header, and lets
// users eyeball-identify keys in logs/UI. Plaintext form: "wgl_<43 b64url>".
const apiKeyPrefix = "wgl_"

// apiKeyDisplayLen is how many leading characters of the plaintext token are
// stored as OrgAPIKey.Prefix for display (e.g. "wgl_Ab12cdEf"). Enough to
// recognise a key without revealing enough to be useful to an attacker.
const apiKeyDisplayLen = 12

// APIKeyService mints, authenticates, lists, and revokes organization-scoped
// API keys. Only the SHA-256 hash of a token is persisted; the plaintext is
// returned once at Issue time and never recoverable afterwards.
type APIKeyService struct {
	keys *repo.OrgAPIKeyRepoGorm
}

func NewAPIKeyService(keys *repo.OrgAPIKeyRepoGorm) *APIKeyService {
	return &APIKeyService{keys: keys}
}

// APIKeyView is the safe, secret-free representation of a key for listing.
type APIKeyView struct {
	ID         uuid.UUID
	Name       string
	Prefix     string
	LastUsedAt *time.Time
	ExpiresAt  *time.Time
	RevokedAt  *time.Time
	CreatedAt  time.Time
}

// IssuedAPIKey is returned once from Issue: the plaintext Token must be shown
// to the caller immediately and cannot be retrieved later.
type IssuedAPIKey struct {
	Token string
	View  APIKeyView
}

// Issue creates a new API key for orgID. createdBy is optional (nil for
// platform-minted keys). expiresAt is optional (nil = never expires).
func (s *APIKeyService) Issue(ctx context.Context, orgID uuid.UUID, name string, createdBy *uuid.UUID, expiresAt *time.Time) (*IssuedAPIKey, error) {
	if s == nil || s.keys == nil {
		return nil, errors.New("api key service: not configured")
	}
	if orgID == uuid.Nil {
		return nil, fmt.Errorf("%w: organization is required", ErrInvalidInput)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if expiresAt != nil && !expiresAt.After(time.Now().UTC()) {
		return nil, fmt.Errorf("%w: expires_at must be in the future", ErrInvalidInput)
	}

	token, err := newAPIKeyToken()
	if err != nil {
		return nil, fmt.Errorf("generate api key: %w", err)
	}

	key := &control.OrgAPIKey{
		OrganizationID:     orgID,
		Name:               name,
		TokenHash:          hashAPIKey(token),
		Prefix:             token[:apiKeyDisplayLen],
		CreatedByAccountID: createdBy,
		ExpiresAt:          expiresAt,
	}
	if err := s.keys.Create(ctx, key); err != nil {
		return nil, err
	}
	return &IssuedAPIKey{Token: token, View: toAPIKeyView(key)}, nil
}

// AuthenticatedKey identifies the org and key behind a successful API-key auth.
type AuthenticatedKey struct {
	KeyID          uuid.UUID
	OrganizationID uuid.UUID
}

// LooksLikeAPIKey reports whether raw is a Waggle API key rather than a JWT, so
// the auth middleware can route it to Authenticate instead of token parsing.
func LooksLikeAPIKey(raw string) bool {
	return strings.HasPrefix(raw, apiKeyPrefix)
}

// Authenticate validates a plaintext API key and returns its org context. It
// updates LastUsedAt best-effort. Returns ErrBadCredentials when the key is
// unknown, revoked, or expired.
func (s *APIKeyService) Authenticate(ctx context.Context, token string) (*AuthenticatedKey, error) {
	if s == nil || s.keys == nil {
		return nil, errors.New("api key service: not configured")
	}
	token = strings.TrimSpace(token)
	if !LooksLikeAPIKey(token) {
		return nil, ErrBadCredentials
	}
	key, err := s.keys.FindActiveByHash(ctx, hashAPIKey(token))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBadCredentials
		}
		return nil, err
	}
	// Best-effort; don't fail the request if the timestamp update fails.
	_ = s.keys.TouchLastUsed(ctx, key.ID)
	return &AuthenticatedKey{KeyID: key.ID, OrganizationID: key.OrganizationID}, nil
}

// List returns all keys for an organization, secret-free.
func (s *APIKeyService) List(ctx context.Context, orgID uuid.UUID) ([]APIKeyView, error) {
	if s == nil || s.keys == nil {
		return nil, errors.New("api key service: not configured")
	}
	keys, err := s.keys.ListByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}
	views := make([]APIKeyView, 0, len(keys))
	for i := range keys {
		views = append(views, toAPIKeyView(&keys[i]))
	}
	return views, nil
}

// Revoke revokes a key owned by orgID. Returns ErrNotFound when no active key
// with that id belongs to the org (already-revoked or wrong-org).
func (s *APIKeyService) Revoke(ctx context.Context, orgID, keyID uuid.UUID) error {
	if s == nil || s.keys == nil {
		return errors.New("api key service: not configured")
	}
	n, err := s.keys.Revoke(ctx, orgID, keyID)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func toAPIKeyView(k *control.OrgAPIKey) APIKeyView {
	return APIKeyView{
		ID:         k.ID,
		Name:       k.Name,
		Prefix:     k.Prefix,
		LastUsedAt: k.LastUsedAt,
		ExpiresAt:  k.ExpiresAt,
		RevokedAt:  k.RevokedAt,
		CreatedAt:  k.CreatedAt,
	}
}

// newAPIKeyToken returns a fresh plaintext key: prefix + 32 bytes of URL-safe
// base64 randomness (256 bits of entropy).
func newAPIKeyToken() (string, error) {
	b, err := utils.RandomBytes(32)
	if err != nil {
		return "", err
	}
	return apiKeyPrefix + base64.RawURLEncoding.EncodeToString(b), nil
}

// hashAPIKey returns the hex SHA-256 of a plaintext token. The token's 256 bits
// of entropy make a plain hash safe (no brute-force surface), so we index it
// directly for O(1) lookup — mirroring HashRefreshToken's approach.
func hashAPIKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
