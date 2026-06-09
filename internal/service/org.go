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
	"github.com/glueops/waggle/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Membership roles, ordered owner > admin > member.
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)

// ErrForbidden is returned when the caller lacks the role for an operation.
var ErrForbidden = errors.New("forbidden")

// ErrAlreadyMember is returned when adding an account that already belongs to
// the org.
var ErrAlreadyMember = errors.New("already a member")

func roleRank(r string) int {
	switch r {
	case RoleOwner:
		return 3
	case RoleAdmin:
		return 2
	case RoleMember:
		return 1
	default:
		return 0
	}
}

func validRole(r string) bool { return roleRank(r) > 0 }

// OrgService manages organizations and their memberships. Org lifecycle changes
// (create/delete) drive tenant provisioning/teardown via the job queue.
type OrgService struct {
	db       *gorm.DB
	tokens   *TokenService
	enqueuer *jobs.Enqueuer
	sender   EmailSender
}

func NewOrgService(db *gorm.DB, tokens *TokenService, enqueuer *jobs.Enqueuer, sender EmailSender) *OrgService {
	return &OrgService{db: db, tokens: tokens, enqueuer: enqueuer, sender: sender}
}

type OrgView struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	Domain    *string
	Status    string
	Role      string // the calling account's role in this org
	CreatedAt time.Time
}

type MemberView struct {
	UserID        uuid.UUID
	AccountID     uuid.UUID
	DisplayName   string
	Email         string
	Role          string
	IsActive      bool
	AccountActive bool
	Pending       bool // invited but hasn't accepted (account not yet active)
	LastLoginAt   *time.Time
	CreatedAt     time.Time
}

// membership loads the caller's active User row for an org, or ErrForbidden if
// they aren't a member.
func (s *OrgService) membership(ctx context.Context, accountID, orgID uuid.UUID) (*control.User, error) {
	var u control.User
	err := s.db.WithContext(ctx).
		Where("account_id = ? AND organization_id = ? AND is_active = ?", accountID, orgID, true).
		First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrForbidden
	}
	if err != nil {
		return nil, fmt.Errorf("load membership: %w", err)
	}
	return &u, nil
}

// requireRole ensures the caller is a member with at least minRole.
func (s *OrgService) requireRole(ctx context.Context, accountID, orgID uuid.UUID, minRole string) (*control.User, error) {
	u, err := s.membership(ctx, accountID, orgID)
	if err != nil {
		return nil, err
	}
	if roleRank(u.Role) < roleRank(minRole) {
		return nil, ErrForbidden
	}
	return u, nil
}

// CreateOrg creates an organization owned by accountID and enqueues tenant
// provisioning.
func (s *OrgService) CreateOrg(ctx context.Context, accountID uuid.UUID, name string) (*OrgView, error) {
	name = strings.TrimSpace(name)
	if name == "" || accountID == uuid.Nil {
		return nil, ErrInvalidInput
	}
	org := control.Organization{Name: name, Slug: utils.Slugify(name), Status: "pending"}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&org).Error; err != nil {
			return fmt.Errorf("create organization: %w", err)
		}
		user := control.User{AccountID: accountID, OrganizationID: org.ID, Role: RoleOwner, IsActive: true}
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("create membership: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := s.enqueuer.EnqueueProvisionTenant(ctx, org.ID); err != nil {
		log.Printf("OrgService.CreateOrg: enqueue provisioner for %s failed: %v", org.ID, err)
	}
	return &OrgView{ID: org.ID, Name: org.Name, Slug: org.Slug, Domain: org.Domain, Status: org.Status, Role: RoleOwner, CreatedAt: org.CreatedAt}, nil
}

// ListOrgs returns the orgs the account belongs to, with the account's role.
func (s *OrgService) ListOrgs(ctx context.Context, accountID uuid.UUID) ([]OrgView, error) {
	type row struct {
		ID        uuid.UUID
		Name      string
		Slug      string
		Domain    *string
		Status    string
		Role      string
		CreatedAt time.Time
	}
	var rows []row
	if err := s.db.WithContext(ctx).
		Table("organizations").
		Select("organizations.id, organizations.name, organizations.slug, organizations.domain, organizations.status, users.role AS role, organizations.created_at").
		Joins("JOIN users ON users.organization_id = organizations.id").
		Where("users.account_id = ? AND users.is_active = ?", accountID, true).
		Order("organizations.created_at").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list orgs: %w", err)
	}
	out := make([]OrgView, 0, len(rows))
	for _, r := range rows {
		out = append(out, OrgView{ID: r.ID, Name: r.Name, Slug: r.Slug, Domain: r.Domain, Status: r.Status, Role: r.Role, CreatedAt: r.CreatedAt})
	}
	return out, nil
}

func (s *OrgService) GetOrg(ctx context.Context, accountID, orgID uuid.UUID) (*OrgView, error) {
	u, err := s.membership(ctx, accountID, orgID)
	if err != nil {
		return nil, err
	}
	var org control.Organization
	if err := s.db.WithContext(ctx).First(&org, "id = ?", orgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get org: %w", err)
	}
	return &OrgView{ID: org.ID, Name: org.Name, Slug: org.Slug, Domain: org.Domain, Status: org.Status, Role: u.Role, CreatedAt: org.CreatedAt}, nil
}

// UpdateOrg renames an org (admin+).
func (s *OrgService) UpdateOrg(ctx context.Context, accountID, orgID uuid.UUID, name string) (*OrgView, error) {
	caller, err := s.requireRole(ctx, accountID, orgID, RoleAdmin)
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidInput
	}
	var org control.Organization
	if err := s.db.WithContext(ctx).First(&org, "id = ?", orgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get org: %w", err)
	}
	org.Name = name
	if err := s.db.WithContext(ctx).Model(&org).Update("name", name).Error; err != nil {
		return nil, fmt.Errorf("update org: %w", err)
	}
	return &OrgView{ID: org.ID, Name: org.Name, Slug: org.Slug, Domain: org.Domain, Status: org.Status, Role: caller.Role, CreatedAt: org.CreatedAt}, nil
}

// DeleteOrg marks an org destroying and enqueues tenant teardown (owner only).
func (s *OrgService) DeleteOrg(ctx context.Context, accountID, orgID uuid.UUID) error {
	if _, err := s.requireRole(ctx, accountID, orgID, RoleOwner); err != nil {
		return err
	}
	res := s.db.WithContext(ctx).
		Model(&control.Organization{}).
		Where("id = ?", orgID).
		Update("status", "destroying")
	if res.Error != nil {
		return fmt.Errorf("mark org destroying: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	if err := s.enqueuer.EnqueueDestroyTenant(ctx, orgID); err != nil {
		log.Printf("OrgService.DeleteOrg: enqueue destroyer for %s failed: %v", orgID, err)
	}
	return nil
}

// ListMembers returns an org's members (any member may view).
func (s *OrgService) ListMembers(ctx context.Context, accountID, orgID uuid.UUID) ([]MemberView, error) {
	if _, err := s.membership(ctx, accountID, orgID); err != nil {
		return nil, err
	}
	type row struct {
		UserID        uuid.UUID
		AccountID     uuid.UUID
		DisplayName   string
		Email         string
		Role          string
		IsActive      bool
		AccountActive bool
		LastLoginAt   *time.Time
		CreatedAt     time.Time
	}
	var rows []row
	if err := s.db.WithContext(ctx).
		Table("users").
		Select("users.id AS user_id, users.account_id AS account_id, accounts.display_name AS display_name, "+
			"account_emails.email AS email, users.role AS role, users.is_active AS is_active, "+
			"accounts.is_active AS account_active, accounts.last_login_at AS last_login_at, users.created_at AS created_at").
		Joins("JOIN accounts ON accounts.id = users.account_id").
		Joins("LEFT JOIN account_emails ON account_emails.account_id = users.account_id AND account_emails.is_primary = true").
		Where("users.organization_id = ?", orgID).
		Order("users.created_at").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	out := make([]MemberView, 0, len(rows))
	for _, r := range rows {
		out = append(out, MemberView{
			UserID: r.UserID, AccountID: r.AccountID, DisplayName: r.DisplayName, Email: r.Email,
			Role: r.Role, IsActive: r.IsActive, AccountActive: r.AccountActive,
			Pending: !r.AccountActive, LastLoginAt: r.LastLoginAt, CreatedAt: r.CreatedAt,
		})
	}
	return out, nil
}

// AddMemberResult reports the outcome of an invite/add.
type AddMemberResult struct {
	Member  MemberView
	Invited bool // an invite email was sent (pending account)
}

// AddMember adds an account (by email) to an org with a role. Existing active
// accounts become members immediately; unknown or not-yet-activated emails get
// a pending account + invite email. Caller must be admin+ (owner to grant
// owner).
func (s *OrgService) AddMember(ctx context.Context, callerAccountID, orgID uuid.UUID, email, role string) (*AddMemberResult, error) {
	caller, err := s.requireRole(ctx, callerAccountID, orgID, RoleAdmin)
	if err != nil {
		return nil, err
	}
	if role == "" {
		role = RoleMember
	}
	if !validRole(role) {
		return nil, fmt.Errorf("%w: invalid role", ErrInvalidInput)
	}
	if role == RoleOwner && caller.Role != RoleOwner {
		return nil, ErrForbidden
	}
	email = utils.NormalizeEmail(email)
	if email == "" {
		return nil, fmt.Errorf("%w: email required", ErrInvalidInput)
	}

	var org control.Organization
	if err := s.db.WithContext(ctx).First(&org, "id = ?", orgID).Error; err != nil {
		return nil, ErrNotFound
	}

	var (
		result  AddMemberResult
		account control.Account
		invite  bool
	)
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Resolve or create the invitee account.
		var ae control.AccountEmail
		emailErr := tx.Where("email = ?", email).First(&ae).Error
		switch {
		case emailErr == nil:
			if err := tx.First(&account, "id = ?", ae.AccountID).Error; err != nil {
				return fmt.Errorf("load account: %w", err)
			}
		case errors.Is(emailErr, gorm.ErrRecordNotFound):
			account = control.Account{DisplayName: email, IsActive: false}
			if err := tx.Create(&account).Error; err != nil {
				return fmt.Errorf("create pending account: %w", err)
			}
			// The Account.IsActive `default:true` tag makes GORM drop the
			// zero-value false on insert, so the row is created active; force it
			// inactive so the invitee must accept (set a password) first.
			if err := tx.Model(&account).Update("is_active", false).Error; err != nil {
				return fmt.Errorf("deactivate pending account: %w", err)
			}
			account.IsActive = false
			if err := tx.Create(&control.AccountEmail{AccountID: account.ID, Email: email, IsPrimary: true}).Error; err != nil {
				return fmt.Errorf("create account email: %w", err)
			}
		default:
			return fmt.Errorf("lookup email: %w", emailErr)
		}

		// Reject duplicate membership.
		var existing control.User
		dup := tx.Where("account_id = ? AND organization_id = ?", account.ID, orgID).First(&existing).Error
		if dup == nil {
			return ErrAlreadyMember
		}
		if !errors.Is(dup, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check membership: %w", dup)
		}

		member := control.User{AccountID: account.ID, OrganizationID: orgID, DisplayName: account.DisplayName, Role: role, IsActive: true}
		if err := tx.Create(&member).Error; err != nil {
			return fmt.Errorf("create membership: %w", err)
		}
		invite = !account.IsActive
		result.Member = MemberView{
			UserID: member.ID, AccountID: account.ID, DisplayName: account.DisplayName, Email: email,
			Role: role, IsActive: true, AccountActive: account.IsActive, Pending: !account.IsActive,
			CreatedAt: member.CreatedAt,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Pending invitees need to set a password; email them an accept link.
	if invite {
		if token, _, ierr := s.tokens.IssueInvite(account.ID, orgID); ierr == nil && s.sender != nil {
			if serr := s.sender.SendInvite(ctx, email, org.Name, token); serr != nil {
				log.Printf("OrgService.AddMember: send invite to %s failed: %v", email, serr)
			}
		}
	}
	result.Invited = invite
	return &result, nil
}

// UpdateMember changes a member's role (admin+; owner required to touch owners).
func (s *OrgService) UpdateMember(ctx context.Context, callerAccountID, orgID, userID uuid.UUID, role string) (*MemberView, error) {
	caller, err := s.requireRole(ctx, callerAccountID, orgID, RoleAdmin)
	if err != nil {
		return nil, err
	}
	if !validRole(role) {
		return nil, fmt.Errorf("%w: invalid role", ErrInvalidInput)
	}
	var target control.User
	if err := s.db.WithContext(ctx).Where("id = ? AND organization_id = ?", userID, orgID).First(&target).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("load member: %w", err)
	}
	// Only owners may promote to or demote from owner.
	if (role == RoleOwner || target.Role == RoleOwner) && caller.Role != RoleOwner {
		return nil, ErrForbidden
	}
	// Keep at least one owner.
	if target.Role == RoleOwner && role != RoleOwner {
		if err := s.guardLastOwner(ctx, orgID, target.ID); err != nil {
			return nil, err
		}
	}
	if err := s.db.WithContext(ctx).Model(&target).Update("role", role).Error; err != nil {
		return nil, fmt.Errorf("update member: %w", err)
	}
	target.Role = role
	return &MemberView{UserID: target.ID, AccountID: target.AccountID, DisplayName: target.DisplayName, Role: target.Role, IsActive: target.IsActive}, nil
}

// RemoveMember removes a membership (admin+; owner required to remove owners;
// never the last owner).
func (s *OrgService) RemoveMember(ctx context.Context, callerAccountID, orgID, userID uuid.UUID) error {
	caller, err := s.requireRole(ctx, callerAccountID, orgID, RoleAdmin)
	if err != nil {
		return err
	}
	var target control.User
	if err := s.db.WithContext(ctx).Where("id = ? AND organization_id = ?", userID, orgID).First(&target).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("load member: %w", err)
	}
	if target.Role == RoleOwner {
		if caller.Role != RoleOwner {
			return ErrForbidden
		}
		if err := s.guardLastOwner(ctx, orgID, target.ID); err != nil {
			return err
		}
	}
	if err := s.db.WithContext(ctx).Delete(&control.User{}, "id = ?", target.ID).Error; err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return nil
}

// guardLastOwner returns ErrInvalidInput if removing/demoting excludeUserID
// would leave the org with no owner.
func (s *OrgService) guardLastOwner(ctx context.Context, orgID, excludeUserID uuid.UUID) error {
	var owners int64
	if err := s.db.WithContext(ctx).
		Model(&control.User{}).
		Where("organization_id = ? AND role = ? AND id <> ?", orgID, RoleOwner, excludeUserID).
		Count(&owners).Error; err != nil {
		return fmt.Errorf("count owners: %w", err)
	}
	if owners == 0 {
		return fmt.Errorf("%w: organization must keep at least one owner", ErrInvalidInput)
	}
	return nil
}
