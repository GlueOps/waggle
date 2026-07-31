package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/glueops/waggle/internal/jobs"
	"github.com/glueops/waggle/internal/models/control"
	"github.com/glueops/waggle/internal/repo"
	"github.com/glueops/waggle/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrEmailExists    = errors.New("email already registered")
	ErrDomainClaimed  = errors.New("organization domain already claimed")
	ErrBadCredentials = errors.New("invalid credentials")
	ErrSessionInvalid = errors.New("session invalid or expired")
	ErrAmbiguousOrg   = errors.New("multiple organization memberships; choose one")
)

const minPasswordLen = 8

// maxPasswordLen is bcrypt's hard input limit. Passwords longer than this make
// utils.HashPassword fail, so they must be rejected as invalid input rather
// than surfacing as a 500 from the hash call.
const maxPasswordLen = 72

// validatePassword enforces the length bounds both password entry points share.
func validatePassword(pw string) error {
	if len(pw) < minPasswordLen {
		return fmt.Errorf("%w: password must be at least %d characters", ErrInvalidInput, minPasswordLen)
	}
	if len(pw) > maxPasswordLen {
		return fmt.Errorf("%w: password must be at most %d characters", ErrInvalidInput, maxPasswordLen)
	}
	return nil
}

type AuthService struct {
	db       *gorm.DB
	tokens   *TokenService
	accounts *repo.AccountRepoGorm
	emails   *repo.AccountEmailRepoGorm
	users    *repo.UserRepoGorm
	sessions *repo.TokenSessionRepoGorm
	audit    *repo.AuthAuditRepoGorm
	enqueuer *jobs.Enqueuer
	sender   EmailSender
}

func NewAuthService(
	db *gorm.DB,
	tokens *TokenService,
	accounts *repo.AccountRepoGorm,
	emails *repo.AccountEmailRepoGorm,
	users *repo.UserRepoGorm,
	sessions *repo.TokenSessionRepoGorm,
	audit *repo.AuthAuditRepoGorm,
	enqueuer *jobs.Enqueuer,
	sender EmailSender,
) *AuthService {
	return &AuthService{
		db:       db,
		tokens:   tokens,
		accounts: accounts,
		emails:   emails,
		users:    users,
		sessions: sessions,
		audit:    audit,
		enqueuer: enqueuer,
		sender:   sender,
	}
}

type SignupInput struct {
	Email            string
	Password         string
	DisplayName      string
	OrganizationName string
	UserAgent        string
	IPAddress        string
}

type LoginInput struct {
	Email          string
	Password       string
	OrganizationID *uuid.UUID
	UserAgent      string
	IPAddress      string
}

type RefreshInput struct {
	RefreshToken string
	UserAgent    string
	IPAddress    string
}

type Membership struct {
	OrganizationID   uuid.UUID
	OrganizationName string
	OrganizationSlug string
}

type AuthResult struct {
	AccountID        uuid.UUID
	OrganizationID   uuid.UUID
	OrganizationSlug string
	OrgStatus        string
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
}

type LoginResult struct {
	Auth        *AuthResult
	Memberships []Membership
}

func (s *AuthService) Signup(ctx context.Context, in SignupInput) (*AuthResult, error) {
	email := utils.NormalizeEmail(in.Email)
	if email == "" || in.OrganizationName == "" {
		return nil, ErrInvalidInput
	}
	if err := validatePassword(in.Password); err != nil {
		return nil, err
	}
	domain, err := utils.ExtractDomain(email)
	if err != nil {
		return nil, ErrInvalidInput
	}

	pwdHash, err := utils.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	refresh, err := s.tokens.GenerateRefresh()
	if err != nil {
		return nil, err
	}

	var (
		account      control.Account
		accountEmail control.AccountEmail
		org          control.Organization
		user         control.User
		session      control.TokenSession
	)

	domainPtr := &domain
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		account = control.Account{
			DisplayName:  in.DisplayName,
			PasswordHash: pwdHash,
			IsActive:     true,
		}
		if err := tx.Create(&account).Error; err != nil {
			return fmt.Errorf("create account: %w", err)
		}

		accountEmail = control.AccountEmail{
			AccountID: account.ID,
			Email:     email,
			IsPrimary: true,
		}
		if err := tx.Create(&accountEmail).Error; err != nil {
			return fmt.Errorf("create account email: %w", err)
		}

		org = control.Organization{
			Name:   in.OrganizationName,
			Slug:   utils.Slugify(in.OrganizationName),
			Domain: domainPtr,
			Status: "pending",
		}
		if err := tx.Create(&org).Error; err != nil {
			return fmt.Errorf("create organization: %w", err)
		}

		user = control.User{
			AccountID:      account.ID,
			OrganizationID: org.ID,
			DisplayName:    in.DisplayName,
			Role:           RoleOwner,
			IsActive:       true,
		}
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		session = control.TokenSession{
			AccountID:        account.ID,
			OrganizationID:   org.ID,
			RefreshTokenHash: refresh.Hashed,
			UserAgent:        in.UserAgent,
			IPAddress:        in.IPAddress,
			ExpiresAt:        refresh.ExpiresAt,
		}
		if err := tx.Create(&session).Error; err != nil {
			return fmt.Errorf("create session: %w", err)
		}

		return nil
	})
	if err != nil {
		if specific := classifyUniqueViolation(err); specific != nil {
			return nil, specific
		}
		return nil, err
	}

	if err := s.enqueuer.EnqueueProvisionTenant(ctx, org.ID); err != nil {
		// Signup succeeded but the provisioner job didn't enqueue. Don't fail the
		// signup — operator can re-enqueue later. Log loudly so it's visible.
		log.Printf("auth.Signup: failed to enqueue tenant provisioner for org %s: %v", org.ID, err)
	}

	if s.sender != nil {
		if vt, _, vErr := s.tokens.IssueEmailVerification(accountEmail.ID); vErr == nil {
			if sErr := s.sender.SendVerification(ctx, accountEmail.Email, vt); sErr != nil {
				log.Printf("auth.Signup: send verification for %s failed: %v", accountEmail.Email, sErr)
			}
		} else {
			log.Printf("auth.Signup: issue verification token for %s failed: %v", accountEmail.Email, vErr)
		}
	}

	access, accessExp, err := s.tokens.IssueAccess(account.ID, org.ID, session.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		AccountID:        account.ID,
		OrganizationID:   org.ID,
		OrganizationSlug: org.Slug,
		OrgStatus:        org.Status,
		AccessToken:      access,
		AccessExpiresAt:  accessExp,
		RefreshToken:     refresh.Plain,
		RefreshExpiresAt: refresh.ExpiresAt,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (*LoginResult, error) {
	email := utils.NormalizeEmail(in.Email)
	if email == "" || in.Password == "" {
		return nil, ErrInvalidInput
	}

	account, err := s.accounts.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrBadCredentials
	}
	if !account.IsActive || !utils.VerifyPassword(in.Password, account.PasswordHash) {
		return nil, ErrBadCredentials
	}

	q := s.db.WithContext(ctx).
		Model(&control.User{}).
		Where("account_id = ?", account.ID).
		Where("is_active = ?", true)
	if in.OrganizationID != nil {
		q = q.Where("organization_id = ?", *in.OrganizationID)
	}

	var users []control.User
	if err := q.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("find memberships: %w", err)
	}
	if len(users) == 0 {
		return nil, ErrBadCredentials
	}
	if in.OrganizationID == nil && len(users) > 1 {
		memberships, err := s.fetchMemberships(ctx, users)
		if err != nil {
			return nil, err
		}
		return &LoginResult{Memberships: memberships}, ErrAmbiguousOrg
	}

	user := users[0]
	now := time.Now().UTC()
	_ = s.db.WithContext(ctx).Model(&control.User{}).Where("id = ?", user.ID).Update("last_login_at", &now).Error
	_ = s.db.WithContext(ctx).Model(&control.Account{}).Where("id = ?", account.ID).Update("last_login_at", &now).Error

	refresh, err := s.tokens.GenerateRefresh()
	if err != nil {
		return nil, err
	}
	session := &control.TokenSession{
		AccountID:        account.ID,
		OrganizationID:   user.OrganizationID,
		RefreshTokenHash: refresh.Hashed,
		UserAgent:        in.UserAgent,
		IPAddress:        in.IPAddress,
		ExpiresAt:        refresh.ExpiresAt,
	}
	if err := s.sessions.Create(ctx, session); err != nil {
		return nil, err
	}

	access, accessExp, err := s.tokens.IssueAccess(account.ID, user.OrganizationID, session.ID)
	if err != nil {
		return nil, err
	}

	var orgSlug, orgStatus string
	var fetched control.Organization
	if err := s.db.WithContext(ctx).First(&fetched, "id = ?", user.OrganizationID).Error; err == nil {
		orgSlug = fetched.Slug
		orgStatus = fetched.Status
	}

	return &LoginResult{
		Auth: &AuthResult{
			AccountID:        account.ID,
			OrganizationID:   user.OrganizationID,
			OrganizationSlug: orgSlug,
			OrgStatus:        orgStatus,
			AccessToken:      access,
			AccessExpiresAt:  accessExp,
			RefreshToken:     refresh.Plain,
			RefreshExpiresAt: refresh.ExpiresAt,
		},
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, in RefreshInput) (*AuthResult, error) {
	if in.RefreshToken == "" {
		return nil, ErrSessionInvalid
	}
	hashed := HashRefreshToken(in.RefreshToken)
	session, err := s.sessions.FindByRefreshHash(ctx, hashed)
	if err != nil {
		return nil, ErrSessionInvalid
	}

	if err := s.sessions.Revoke(ctx, session.ID); err != nil {
		return nil, err
	}

	newRefresh, err := s.tokens.GenerateRefresh()
	if err != nil {
		return nil, err
	}
	newSession := &control.TokenSession{
		AccountID:        session.AccountID,
		OrganizationID:   session.OrganizationID,
		RefreshTokenHash: newRefresh.Hashed,
		UserAgent:        in.UserAgent,
		IPAddress:        in.IPAddress,
		ExpiresAt:        newRefresh.ExpiresAt,
	}
	if err := s.sessions.Create(ctx, newSession); err != nil {
		return nil, err
	}

	access, accessExp, err := s.tokens.IssueAccess(session.AccountID, session.OrganizationID, newSession.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		AccountID:        session.AccountID,
		OrganizationID:   session.OrganizationID,
		AccessToken:      access,
		AccessExpiresAt:  accessExp,
		RefreshToken:     newRefresh.Plain,
		RefreshExpiresAt: newRefresh.ExpiresAt,
	}, nil
}

// VerifyEmail consumes a verification token. Idempotent — verifying an
// already-verified address is a no-op success.
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidInput
	}
	aemID, err := s.tokens.VerifyEmailVerification(token)
	if err != nil {
		return ErrSessionInvalid
	}
	var ae control.AccountEmail
	if err := s.db.WithContext(ctx).First(&ae, "id = ?", aemID).Error; err != nil {
		return ErrSessionInvalid
	}
	if ae.VerifiedAt != nil {
		return nil
	}
	now := time.Now().UTC()
	return s.db.WithContext(ctx).
		Model(&control.AccountEmail{}).
		Where("id = ?", aemID).
		Update("verified_at", &now).Error
}

type AccountEmailView struct {
	ID         uuid.UUID
	Email      string
	IsPrimary  bool
	VerifiedAt *time.Time
}

type MeResult struct {
	AccountID        uuid.UUID
	DisplayName      string
	LastLoginAt      *time.Time
	Emails           []AccountEmailView
	Memberships      []Membership
	CurrentOrgID     uuid.UUID
	CurrentOrgSlug   string
	CurrentOrgStatus string
}

// Me returns the current account, its emails, its memberships, and (if the
// caller's session has an org context) the current organization summary.
func (s *AuthService) Me(ctx context.Context, accountID, currentOrgID uuid.UUID) (*MeResult, error) {
	if accountID == uuid.Nil {
		return nil, ErrSessionInvalid
	}

	var account control.Account
	if err := s.db.WithContext(ctx).First(&account, "id = ?", accountID).Error; err != nil {
		return nil, ErrSessionInvalid
	}
	if !account.IsActive {
		return nil, ErrBadCredentials
	}

	var emails []control.AccountEmail
	if err := s.db.WithContext(ctx).Where("account_id = ?", accountID).Find(&emails).Error; err != nil {
		return nil, fmt.Errorf("load emails: %w", err)
	}
	emailViews := make([]AccountEmailView, 0, len(emails))
	for _, e := range emails {
		emailViews = append(emailViews, AccountEmailView{
			ID:         e.ID,
			Email:      e.Email,
			IsPrimary:  e.IsPrimary,
			VerifiedAt: e.VerifiedAt,
		})
	}

	var users []control.User
	if err := s.db.WithContext(ctx).Where("account_id = ?", accountID).Where("is_active = ?", true).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("load memberships: %w", err)
	}
	memberships, err := s.fetchMemberships(ctx, users)
	if err != nil {
		return nil, err
	}

	res := &MeResult{
		AccountID:   account.ID,
		DisplayName: account.DisplayName,
		LastLoginAt: account.LastLoginAt,
		Emails:      emailViews,
		Memberships: memberships,
	}

	if currentOrgID != uuid.Nil {
		var org control.Organization
		if err := s.db.WithContext(ctx).First(&org, "id = ?", currentOrgID).Error; err == nil {
			res.CurrentOrgID = org.ID
			res.CurrentOrgSlug = org.Slug
			res.CurrentOrgStatus = org.Status
		}
	}

	return res, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	hashed := HashRefreshToken(refreshToken)
	session, err := s.sessions.FindByRefreshHash(ctx, hashed)
	if err != nil {
		return nil // idempotent — unknown/expired token is already "logged out"
	}
	return s.sessions.Revoke(ctx, session.ID)
}

// issueOrgSession creates a session and access/refresh pair for an account in
// an org (shared by SwitchOrg and AcceptInvite).
func (s *AuthService) issueOrgSession(ctx context.Context, accountID, orgID uuid.UUID, userAgent, ip string) (*AuthResult, error) {
	refresh, err := s.tokens.GenerateRefresh()
	if err != nil {
		return nil, err
	}
	session := &control.TokenSession{
		AccountID:        accountID,
		OrganizationID:   orgID,
		RefreshTokenHash: refresh.Hashed,
		UserAgent:        userAgent,
		IPAddress:        ip,
		ExpiresAt:        refresh.ExpiresAt,
	}
	if err := s.sessions.Create(ctx, session); err != nil {
		return nil, err
	}
	access, accessExp, err := s.tokens.IssueAccess(accountID, orgID, session.ID)
	if err != nil {
		return nil, err
	}
	var org control.Organization
	_ = s.db.WithContext(ctx).First(&org, "id = ?", orgID).Error
	return &AuthResult{
		AccountID:        accountID,
		OrganizationID:   orgID,
		OrganizationSlug: org.Slug,
		OrgStatus:        org.Status,
		AccessToken:      access,
		AccessExpiresAt:  accessExp,
		RefreshToken:     refresh.Plain,
		RefreshExpiresAt: refresh.ExpiresAt,
	}, nil
}

// SwitchOrg issues a fresh token pair scoped to a different org the account is
// an active member of. Returns ErrForbidden if not a member.
func (s *AuthService) SwitchOrg(ctx context.Context, accountID, orgID uuid.UUID, userAgent, ip string) (*AuthResult, error) {
	if accountID == uuid.Nil || orgID == uuid.Nil {
		return nil, ErrInvalidInput
	}
	var account control.Account
	if err := s.db.WithContext(ctx).First(&account, "id = ?", accountID).Error; err != nil {
		return nil, ErrSessionInvalid
	}
	if !account.IsActive {
		return nil, ErrBadCredentials
	}
	var user control.User
	if err := s.db.WithContext(ctx).
		Where("account_id = ? AND organization_id = ? AND is_active = ?", accountID, orgID, true).
		First(&user).Error; err != nil {
		return nil, ErrForbidden
	}
	return s.issueOrgSession(ctx, accountID, orgID, userAgent, ip)
}

type AcceptInviteInput struct {
	Token       string
	Password    string
	DisplayName string
	UserAgent   string
	IPAddress   string
}

// AcceptInvite consumes an org-invite token: for a pending account it sets the
// password and activates the account + its primary email; for an already-active
// account it just confirms the membership. Either way it logs the user into the
// invited org.
func (s *AuthService) AcceptInvite(ctx context.Context, in AcceptInviteInput) (*AuthResult, error) {
	accountID, orgID, err := s.tokens.VerifyInvite(in.Token)
	if err != nil {
		return nil, ErrSessionInvalid
	}
	var account control.Account
	if err := s.db.WithContext(ctx).First(&account, "id = ?", accountID).Error; err != nil {
		return nil, ErrSessionInvalid
	}

	now := time.Now().UTC()
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if !account.IsActive {
			if verr := validatePassword(in.Password); verr != nil {
				return verr
			}
			hash, herr := utils.HashPassword(in.Password)
			if herr != nil {
				return herr
			}
			updates := map[string]any{"is_active": true, "password_hash": hash}
			if strings.TrimSpace(in.DisplayName) != "" {
				updates["display_name"] = in.DisplayName
			}
			if err := tx.Model(&control.Account{}).Where("id = ?", accountID).Updates(updates).Error; err != nil {
				return fmt.Errorf("activate account: %w", err)
			}
			if err := tx.Model(&control.AccountEmail{}).
				Where("account_id = ? AND is_primary = ?", accountID, true).
				Update("verified_at", &now).Error; err != nil {
				return fmt.Errorf("verify email: %w", err)
			}
		}
		// Membership was created at invite time; ensure it's active.
		res := tx.Model(&control.User{}).
			Where("account_id = ? AND organization_id = ?", accountID, orgID).
			Update("is_active", true)
		if res.Error != nil {
			return fmt.Errorf("activate membership: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return ErrSessionInvalid // invite revoked (membership removed)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.issueOrgSession(ctx, accountID, orgID, in.UserAgent, in.IPAddress)
}

func (s *AuthService) fetchMemberships(ctx context.Context, users []control.User) ([]Membership, error) {
	ids := make([]uuid.UUID, 0, len(users))
	for _, u := range users {
		ids = append(ids, u.OrganizationID)
	}
	var orgs []control.Organization
	if err := s.db.WithContext(ctx).Where("id IN ?", ids).Find(&orgs).Error; err != nil {
		return nil, fmt.Errorf("fetch orgs: %w", err)
	}
	out := make([]Membership, 0, len(orgs))
	for _, o := range orgs {
		out = append(out, Membership{
			OrganizationID:   o.ID,
			OrganizationName: o.Name,
			OrganizationSlug: o.Slug,
		})
	}
	return out, nil
}

func classifyUniqueViolation(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return nil
	}
	switch pgErr.ConstraintName {
	case "idx_account_emails_email":
		return ErrEmailExists
	case "idx_organizations_domain":
		return ErrDomainClaimed
	}
	return nil
}
