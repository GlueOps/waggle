package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/glueops/waggle/internal/models/control"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenSessionRepoGorm struct {
	db *gorm.DB
}

func NewTokenSessionRepoGorm(db *gorm.DB) *TokenSessionRepoGorm {
	return &TokenSessionRepoGorm{db: db}
}

func (r *TokenSessionRepoGorm) Create(ctx context.Context, s *control.TokenSession) error {
	if r == nil || r.db == nil {
		return errors.New("token session repo: db is nil")
	}
	if s == nil {
		return errors.New("token session repo: session is nil")
	}
	if err := r.db.WithContext(ctx).Create(s).Error; err != nil {
		return fmt.Errorf("create token session: %w", err)
	}
	return nil
}

func (r *TokenSessionRepoGorm) FindByRefreshHash(ctx context.Context, hash string) (*control.TokenSession, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("token session repo: db is nil")
	}
	var s control.TokenSession
	err := r.db.WithContext(ctx).
		Where("refresh_token_hash = ?", hash).
		Where("revoked_at IS NULL").
		Where("expires_at > ?", time.Now().UTC()).
		First(&s).Error
	if err != nil {
		return nil, fmt.Errorf("find token session: %w", err)
	}
	return &s, nil
}

func (r *TokenSessionRepoGorm) Revoke(ctx context.Context, id uuid.UUID) error {
	if r == nil || r.db == nil {
		return errors.New("token session repo: db is nil")
	}
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).
		Model(&control.TokenSession{}).
		Where("id = ?", id).
		Where("revoked_at IS NULL").
		Update("revoked_at", &now)
	if res.Error != nil {
		return fmt.Errorf("revoke token session: %w", res.Error)
	}
	return nil
}

type AuthAuditRepoGorm struct {
	db *gorm.DB
}

func NewAuthAuditRepoGorm(db *gorm.DB) *AuthAuditRepoGorm {
	return &AuthAuditRepoGorm{db: db}
}

func (r *AuthAuditRepoGorm) Record(ctx context.Context, e *control.AuthAuditEvent) error {
	if r == nil || r.db == nil {
		return errors.New("auth audit repo: db is nil")
	}
	if e == nil {
		return errors.New("auth audit repo: event is nil")
	}
	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return fmt.Errorf("record audit event: %w", err)
	}
	return nil
}

type UserRepoGorm struct {
	db *gorm.DB
}

func NewUserRepoGorm(db *gorm.DB) UserRepoGorm {
	return UserRepoGorm{db: db}
}

func (r UserRepoGorm) FindByEmail(ctx context.Context, orgID uuid.UUID, email string) (*control.User, error) {
	if r.db == nil {
		return nil, errors.New("user repo: db is nil")
	}
	var u control.User
	err := r.db.WithContext(ctx).
		Joins("JOIN account_emails ON account_emails.account_id = users.account_id").
		Where("users.organization_id = ?", orgID).
		Where("account_emails.email = ?", email).
		First(&u).Error
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &u, nil
}

func (r UserRepoGorm) Create(ctx context.Context, u *control.User) error {
	if r.db == nil {
		return errors.New("user repo: db is nil")
	}
	if u == nil {
		return errors.New("user repo: user is nil")
	}
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r UserRepoGorm) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	if r.db == nil {
		return errors.New("user repo: db is nil")
	}
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).
		Model(&control.User{}).
		Where("id = ?", id).
		Update("last_login_at", &now)
	if res.Error != nil {
		return fmt.Errorf("update last login: %w", res.Error)
	}
	return nil
}

type AccountRepoGorm struct {
	db *gorm.DB
}

func NewAccountRepoGorm(db *gorm.DB) *AccountRepoGorm {
	return &AccountRepoGorm{db: db}
}

func (r *AccountRepoGorm) FindByEmail(ctx context.Context, email string) (*control.Account, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("account repo: db is nil")
	}
	var a control.Account
	err := r.db.WithContext(ctx).
		Joins("JOIN account_emails ON account_emails.account_id = accounts.id").
		Where("account_emails.email = ?", email).
		First(&a).Error
	if err != nil {
		return nil, fmt.Errorf("find account by email: %w", err)
	}
	return &a, nil
}

func (r *AccountRepoGorm) Create(ctx context.Context, a *control.Account) error {
	if r == nil || r.db == nil {
		return errors.New("account repo: db is nil")
	}
	if a == nil {
		return errors.New("account repo: account is nil")
	}
	if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
		return fmt.Errorf("create account: %w", err)
	}
	return nil
}

func (r *AccountRepoGorm) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	if r == nil || r.db == nil {
		return errors.New("account repo: db is nil")
	}
	res := r.db.WithContext(ctx).
		Model(&control.Account{}).
		Where("id = ?", id).
		Update("updated_at", time.Now().UTC())
	if res.Error != nil {
		return fmt.Errorf("touch account: %w", res.Error)
	}
	return nil
}

type PlatformPolicyRepoGorm struct {
	db *gorm.DB
}

func NewPlatformPolicyRepoGorm(db *gorm.DB) *PlatformPolicyRepoGorm {
	return &PlatformPolicyRepoGorm{db: db}
}

type OrgPolicyRepoGorm struct {
	db *gorm.DB
}

func NewOrgPolicyRepoGorm(db *gorm.DB) *OrgPolicyRepoGorm {
	return &OrgPolicyRepoGorm{db: db}
}

type UserPasskeyRepoGorm struct {
	db *gorm.DB
}

func NewUserPasskeyRepoGorm(db *gorm.DB) *UserPasskeyRepoGorm {
	return &UserPasskeyRepoGorm{db: db}
}

type AccountEmailRepoGorm struct {
	db *gorm.DB
}

func NewAccountEmailRepoGorm(db *gorm.DB) *AccountEmailRepoGorm {
	return &AccountEmailRepoGorm{db: db}
}

func (r *AccountEmailRepoGorm) FindByEmail(ctx context.Context, email string) (*control.AccountEmail, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("account email repo: db is nil")
	}
	var ae control.AccountEmail
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&ae).Error; err != nil {
		return nil, fmt.Errorf("find account email: %w", err)
	}
	return &ae, nil
}

func (r *AccountEmailRepoGorm) Create(ctx context.Context, ae *control.AccountEmail) error {
	if r == nil || r.db == nil {
		return errors.New("account email repo: db is nil")
	}
	if ae == nil {
		return errors.New("account email repo: nil entity")
	}
	if err := r.db.WithContext(ctx).Create(ae).Error; err != nil {
		return fmt.Errorf("create account email: %w", err)
	}
	return nil
}

type OrgAPIKeyRepoGorm struct {
	db *gorm.DB
}

func NewOrgAPIKeyRepoGorm(db *gorm.DB) *OrgAPIKeyRepoGorm {
	return &OrgAPIKeyRepoGorm{db: db}
}

func (r *OrgAPIKeyRepoGorm) Create(ctx context.Context, k *control.OrgAPIKey) error {
	if r == nil || r.db == nil {
		return errors.New("org api key repo: db is nil")
	}
	if k == nil {
		return errors.New("org api key repo: key is nil")
	}
	if err := r.db.WithContext(ctx).Create(k).Error; err != nil {
		return fmt.Errorf("create org api key: %w", err)
	}
	return nil
}

// FindActiveByHash returns the key matching tokenHash that is neither revoked
// nor expired. The unique index on token_hash makes this an O(1) lookup.
func (r *OrgAPIKeyRepoGorm) FindActiveByHash(ctx context.Context, tokenHash string) (*control.OrgAPIKey, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("org api key repo: db is nil")
	}
	var k control.OrgAPIKey
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		Where("revoked_at IS NULL").
		Where("expires_at IS NULL OR expires_at > ?", time.Now().UTC()).
		First(&k).Error
	if err != nil {
		return nil, fmt.Errorf("find org api key: %w", err)
	}
	return &k, nil
}

// ListByOrg returns all keys (active or not) for an organization, newest first,
// so callers can show revoked/expired keys in management UIs.
func (r *OrgAPIKeyRepoGorm) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]control.OrgAPIKey, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("org api key repo: db is nil")
	}
	var keys []control.OrgAPIKey
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Find(&keys).Error
	if err != nil {
		return nil, fmt.Errorf("list org api keys: %w", err)
	}
	return keys, nil
}

// Revoke marks a key revoked, scoped to its owning org so one tenant can't
// revoke another's key. Returns the number of rows affected.
func (r *OrgAPIKeyRepoGorm) Revoke(ctx context.Context, orgID, id uuid.UUID) (int64, error) {
	if r == nil || r.db == nil {
		return 0, errors.New("org api key repo: db is nil")
	}
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).
		Model(&control.OrgAPIKey{}).
		Where("id = ?", id).
		Where("organization_id = ?", orgID).
		Where("revoked_at IS NULL").
		Update("revoked_at", &now)
	if res.Error != nil {
		return 0, fmt.Errorf("revoke org api key: %w", res.Error)
	}
	return res.RowsAffected, nil
}

// TouchLastUsed records the most recent authentication for a key. Best-effort:
// callers typically ignore the error so auth latency isn't coupled to it.
func (r *OrgAPIKeyRepoGorm) TouchLastUsed(ctx context.Context, id uuid.UUID) error {
	if r == nil || r.db == nil {
		return errors.New("org api key repo: db is nil")
	}
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).
		Model(&control.OrgAPIKey{}).
		Where("id = ?", id).
		Update("last_used_at", &now)
	if res.Error != nil {
		return fmt.Errorf("touch org api key: %w", res.Error)
	}
	return nil
}
