package service

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/glueops/waggle/internal/repo"
	"github.com/glueops/waggle/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct {
	secret     []byte
	issuer     string
	audience   string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenService(secret, issuer string, accessTTL, refreshTTL time.Duration, audience string) (*TokenService, error) {
	if secret == "" {
		return nil, errors.New("token service: secret is required")
	}
	if audience == "" {
		audience = "waggle"
	}
	return &TokenService{
		secret:     []byte(secret),
		issuer:     issuer,
		audience:   audience,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

// Claims are the JWT payload waggle issues. AccountID is duplicated in Subject
// for convenience; OrganizationID is the currently-selected org context for
// this session (zero UUID means "no org context selected").
type Claims struct {
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"org_id,omitempty"`
	SessionID      uuid.UUID `json:"sid"`
	jwt.RegisteredClaims
}

// RefreshTokenPair pairs the plaintext refresh token (returned to the
// client) with the hash stored on TokenSession.RefreshTokenHash.
type RefreshTokenPair struct {
	Plain     string
	Hashed    string
	ExpiresAt time.Time
}

// GenerateRefresh produces a new opaque refresh token. Returns both the
// plaintext (for the client) and the hash (for storage). The caller is
// responsible for creating a TokenSession with these values.
func (ts *TokenService) GenerateRefresh() (*RefreshTokenPair, error) {
	plain, err := newRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh: %w", err)
	}
	return &RefreshTokenPair{
		Plain:     plain,
		Hashed:    HashRefreshToken(plain),
		ExpiresAt: time.Now().UTC().Add(ts.refreshTTL),
	}, nil
}

// IssueAccess signs an access JWT for an existing TokenSession.
func (ts *TokenService) IssueAccess(accountID, orgID, sessionID uuid.UUID) (string, time.Time, error) {
	if accountID == uuid.Nil || sessionID == uuid.Nil {
		return "", time.Time{}, errors.New("issue access: account_id and session_id are required")
	}
	now := time.Now().UTC()
	exp := now.Add(ts.accessTTL)
	claims := Claims{
		AccountID:      accountID,
		OrganizationID: orgID,
		SessionID:      sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ts.issuer,
			Subject:   accountID.String(),
			Audience:  jwt.ClaimStrings{ts.audience},
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(ts.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}
	return signed, exp, nil
}

// Verify parses and validates an access JWT, returning its claims.
func (ts *TokenService) Verify(tokenStr string) (*Claims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(ts.issuer),
		jwt.WithAudience(ts.audience),
		jwt.WithExpirationRequired(),
	)
	tok, err := parser.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return ts.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

const emailVerifyAudience = "waggle-email-verify"

// EmailVerifyTTL is how long a verification link stays valid. 24h is
// standard for email verification.
const EmailVerifyTTL = 24 * time.Hour

// EmailVerificationClaims are the JWT payload for verification links.
// Uses a separate audience so an access JWT can't masquerade as a verify
// token (and vice-versa).
type EmailVerificationClaims struct {
	AccountEmailID uuid.UUID `json:"aem"`
	jwt.RegisteredClaims
}

// IssueEmailVerification signs a verification JWT that, when later
// presented to VerifyEmailVerification, identifies an AccountEmail row.
func (ts *TokenService) IssueEmailVerification(accountEmailID uuid.UUID) (string, time.Time, error) {
	if accountEmailID == uuid.Nil {
		return "", time.Time{}, errors.New("issue email verify: account_email_id required")
	}
	now := time.Now().UTC()
	exp := now.Add(EmailVerifyTTL)
	claims := EmailVerificationClaims{
		AccountEmailID: accountEmailID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ts.issuer,
			Audience:  jwt.ClaimStrings{emailVerifyAudience},
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(ts.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign email-verify token: %w", err)
	}
	return signed, exp, nil
}

// VerifyEmailVerification parses a verification JWT and returns the
// AccountEmail.ID it identifies.
func (ts *TokenService) VerifyEmailVerification(tokenStr string) (uuid.UUID, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(ts.issuer),
		jwt.WithAudience(emailVerifyAudience),
		jwt.WithExpirationRequired(),
	)
	tok, err := parser.ParseWithClaims(tokenStr, &EmailVerificationClaims{}, func(t *jwt.Token) (any, error) {
		return ts.secret, nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("verify email-verify token: %w", err)
	}
	c, ok := tok.Claims.(*EmailVerificationClaims)
	if !ok || !tok.Valid {
		return uuid.Nil, errors.New("invalid email-verify claims")
	}
	return c.AccountEmailID, nil
}

const inviteAudience = "waggle-invite"

// InviteTTL is how long an org invite link stays valid.
const InviteTTL = 7 * 24 * time.Hour

// InviteClaims are the JWT payload for org-invite links. The separate audience
// stops an access/verify token from being replayed as an invite (and vice
// versa). AccountID is the pending (or existing) invitee account.
type InviteClaims struct {
	AccountID      uuid.UUID `json:"acc"`
	OrganizationID uuid.UUID `json:"org"`
	jwt.RegisteredClaims
}

// IssueInvite signs an invite JWT binding an invitee account to an org.
func (ts *TokenService) IssueInvite(accountID, orgID uuid.UUID) (string, time.Time, error) {
	if accountID == uuid.Nil || orgID == uuid.Nil {
		return "", time.Time{}, errors.New("issue invite: account_id and org_id required")
	}
	now := time.Now().UTC()
	exp := now.Add(InviteTTL)
	claims := InviteClaims{
		AccountID:      accountID,
		OrganizationID: orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ts.issuer,
			Audience:  jwt.ClaimStrings{inviteAudience},
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(ts.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign invite token: %w", err)
	}
	return signed, exp, nil
}

// VerifyInvite parses an invite JWT and returns the invitee account + org.
func (ts *TokenService) VerifyInvite(tokenStr string) (uuid.UUID, uuid.UUID, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(ts.issuer),
		jwt.WithAudience(inviteAudience),
		jwt.WithExpirationRequired(),
	)
	tok, err := parser.ParseWithClaims(tokenStr, &InviteClaims{}, func(t *jwt.Token) (any, error) {
		return ts.secret, nil
	})
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("verify invite token: %w", err)
	}
	c, ok := tok.Claims.(*InviteClaims)
	if !ok || !tok.Valid {
		return uuid.Nil, uuid.Nil, errors.New("invalid invite claims")
	}
	return c.AccountID, c.OrganizationID, nil
}

// HashRefreshToken returns the storage representation of an opaque refresh
// token. The same token always hashes to the same value, so lookup is a
// simple equality check against TokenSession.RefreshTokenHash.
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.StdEncoding.EncodeToString(sum[:])
}

func newRefreshToken() (string, error) {
	b, err := utils.RandomBytes(32)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

type TokenSessionService struct {
	tokens   *TokenService
	sessions *repo.TokenSessionRepoGorm
	audit    *repo.AuthAuditRepoGorm
	users    *repo.UserRepoGorm
}

func NewTokenSessionService(
	tokens *TokenService,
	sessions *repo.TokenSessionRepoGorm,
	audit *repo.AuthAuditRepoGorm,
	users *repo.UserRepoGorm,
) *TokenSessionService {
	return &TokenSessionService{
		tokens:   tokens,
		sessions: sessions,
		audit:    audit,
		users:    users,
	}
}
