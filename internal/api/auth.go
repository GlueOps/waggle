package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/glueops/waggle/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type orgView struct {
	ID     uuid.UUID `json:"id"`
	Slug   string    `json:"slug"`
	Status string    `json:"status"`
}

type membershipView struct {
	OrganizationID   uuid.UUID `json:"organization_id"`
	OrganizationName string    `json:"organization_name"`
	OrganizationSlug string    `json:"organization_slug"`
}

type authTokens struct {
	AccessToken      string    `json:"access_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at" format:"date-time"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at" format:"date-time"`
}

type signupInput struct {
	Body struct {
		Email            string `json:"email" format:"email" required:"true" maxLength:"255"`
		Password         string `json:"password" required:"true" minLength:"8" maxLength:"128"`
		DisplayName      string `json:"display_name,omitempty" maxLength:"255"`
		OrganizationName string `json:"organization_name" required:"true" maxLength:"255"`
	}
}

type signupOutput struct {
	Body struct {
		AccountID        uuid.UUID `json:"account_id"`
		Organization     orgView   `json:"organization"`
		AccessToken      string    `json:"access_token"`
		AccessExpiresAt  time.Time `json:"access_expires_at" format:"date-time"`
		RefreshToken     string    `json:"refresh_token"`
		RefreshExpiresAt time.Time `json:"refresh_expires_at" format:"date-time"`
	}
}

type loginInput struct {
	Body struct {
		Email          string     `json:"email" format:"email" required:"true" maxLength:"255"`
		Password       string     `json:"password" required:"true" maxLength:"128"`
		OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	}
}

type loginOutput struct {
	Body struct {
		AccountID    *uuid.UUID       `json:"account_id,omitempty"`
		Organization *orgView         `json:"organization,omitempty"`
		Tokens       *authTokens      `json:"tokens,omitempty"`
		Memberships  []membershipView `json:"memberships,omitempty"`
	}
}

type refreshInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token" required:"true"`
	}
}

type refreshOutput struct {
	Body struct {
		AccessToken      string    `json:"access_token"`
		AccessExpiresAt  time.Time `json:"access_expires_at" format:"date-time"`
		RefreshToken     string    `json:"refresh_token"`
		RefreshExpiresAt time.Time `json:"refresh_expires_at" format:"date-time"`
	}
}

type logoutInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token" required:"true"`
	}
}

type verifyEmailInput struct {
	Body struct {
		Token string `json:"token" required:"true"`
	}
}

type accountEmailView struct {
	ID         uuid.UUID  `json:"id"`
	Email      string     `json:"email"`
	IsPrimary  bool       `json:"is_primary"`
	VerifiedAt *time.Time `json:"verified_at,omitempty" format:"date-time"`
}

type meOutput struct {
	Body struct {
		AccountID    uuid.UUID          `json:"account_id"`
		DisplayName  string             `json:"display_name"`
		LastLoginAt  *time.Time         `json:"last_login_at,omitempty" format:"date-time"`
		Emails       []accountEmailView `json:"emails"`
		Memberships  []membershipView   `json:"memberships"`
		CurrentOrg   *orgView           `json:"current_organization,omitempty"`
	}
}

func (s *Server) registerAuth(auth *service.AuthService, tokens *service.TokenService) {
	if auth == nil || tokens == nil {
		return
	}
	requireAuth := RequireAuth(s.API, tokens)

	huma.Register(s.API, huma.Operation{
		OperationID: "auth-signup",
		Method:      http.MethodPost,
		Path:        "/auth/signup",
		Summary:     "Create an account, organization, and first user; enqueue tenant provisioning.",
		Tags:        []string{"auth"},
	}, func(ctx context.Context, in *signupInput) (*signupOutput, error) {
		res, err := auth.Signup(ctx, service.SignupInput{
			Email:            in.Body.Email,
			Password:         in.Body.Password,
			DisplayName:      in.Body.DisplayName,
			OrganizationName: in.Body.OrganizationName,
			UserAgent:        UserAgent(ctx),
			IPAddress:        ClientIP(ctx),
		})
		if err != nil {
			return nil, mapAuthError(err)
		}
		out := &signupOutput{}
		out.Body.AccountID = res.AccountID
		out.Body.Organization = orgView{ID: res.OrganizationID, Slug: res.OrganizationSlug, Status: res.OrgStatus}
		out.Body.AccessToken = res.AccessToken
		out.Body.AccessExpiresAt = res.AccessExpiresAt
		out.Body.RefreshToken = res.RefreshToken
		out.Body.RefreshExpiresAt = res.RefreshExpiresAt
		return out, nil
	})

	huma.Register(s.API, huma.Operation{
		OperationID: "auth-login",
		Method:      http.MethodPost,
		Path:        "/auth/login",
		Summary:     "Exchange credentials for an access + refresh token. Returns membership list when no organization_id is given and multiple memberships exist.",
		Tags:        []string{"auth"},
	}, func(ctx context.Context, in *loginInput) (*loginOutput, error) {
		res, err := auth.Login(ctx, service.LoginInput{
			Email:          in.Body.Email,
			Password:       in.Body.Password,
			OrganizationID: in.Body.OrganizationID,
			UserAgent:      UserAgent(ctx),
			IPAddress:      ClientIP(ctx),
		})
		out := &loginOutput{}
		if err != nil && errors.Is(err, service.ErrAmbiguousOrg) && res != nil {
			out.Body.Memberships = toMembershipViews(res.Memberships)
			return out, nil
		}
		if err != nil {
			return nil, mapAuthError(err)
		}
		acctID := res.Auth.AccountID
		out.Body.AccountID = &acctID
		out.Body.Organization = &orgView{ID: res.Auth.OrganizationID, Slug: res.Auth.OrganizationSlug, Status: res.Auth.OrgStatus}
		out.Body.Tokens = &authTokens{
			AccessToken:      res.Auth.AccessToken,
			AccessExpiresAt:  res.Auth.AccessExpiresAt,
			RefreshToken:     res.Auth.RefreshToken,
			RefreshExpiresAt: res.Auth.RefreshExpiresAt,
		}
		return out, nil
	})

	huma.Register(s.API, huma.Operation{
		OperationID: "auth-refresh",
		Method:      http.MethodPost,
		Path:        "/auth/refresh",
		Summary:     "Rotate a refresh token, returning a new access + refresh pair.",
		Tags:        []string{"auth"},
	}, func(ctx context.Context, in *refreshInput) (*refreshOutput, error) {
		res, err := auth.Refresh(ctx, service.RefreshInput{
			RefreshToken: in.Body.RefreshToken,
			UserAgent:    UserAgent(ctx),
			IPAddress:    ClientIP(ctx),
		})
		if err != nil {
			return nil, mapAuthError(err)
		}
		out := &refreshOutput{}
		out.Body.AccessToken = res.AccessToken
		out.Body.AccessExpiresAt = res.AccessExpiresAt
		out.Body.RefreshToken = res.RefreshToken
		out.Body.RefreshExpiresAt = res.RefreshExpiresAt
		return out, nil
	})

	huma.Register(s.API, huma.Operation{
		OperationID:   "auth-logout",
		Method:        http.MethodPost,
		Path:          "/auth/logout",
		Summary:       "Revoke the supplied refresh token's session. Idempotent.",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *logoutInput) (*struct{}, error) {
		if err := auth.Logout(ctx, in.Body.RefreshToken); err != nil {
			return nil, mapAuthError(err)
		}
		return nil, nil
	})

	huma.Register(s.API, huma.Operation{
		OperationID:   "auth-verify-email",
		Method:        http.MethodPost,
		Path:          "/auth/verify-email",
		Summary:       "Consume a verification token to mark an email address verified. Idempotent.",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *verifyEmailInput) (*struct{}, error) {
		if err := auth.VerifyEmail(ctx, in.Body.Token); err != nil {
			return nil, mapAuthError(err)
		}
		return nil, nil
	})

	huma.Register(s.API, huma.Operation{
		OperationID: "auth-me",
		Method:      http.MethodGet,
		Path:        "/auth/me",
		Summary:     "Return the authenticated account, its verified-or-pending emails, and its org memberships.",
		Tags:        []string{"auth"},
		Security:    []map[string][]string{{"bearer": {}}},
		Middlewares: huma.Middlewares{requireAuth},
	}, func(ctx context.Context, _ *struct{}) (*meOutput, error) {
		claims, ok := ClaimsFromContext(ctx)
		if !ok {
			return nil, huma.Error401Unauthorized("missing claims")
		}
		res, err := auth.Me(ctx, claims.AccountID, claims.OrganizationID)
		if err != nil {
			return nil, mapAuthError(err)
		}
		out := &meOutput{}
		out.Body.AccountID = res.AccountID
		out.Body.DisplayName = res.DisplayName
		out.Body.LastLoginAt = res.LastLoginAt
		out.Body.Emails = make([]accountEmailView, 0, len(res.Emails))
		for _, e := range res.Emails {
			out.Body.Emails = append(out.Body.Emails, accountEmailView{
				ID: e.ID, Email: e.Email, IsPrimary: e.IsPrimary, VerifiedAt: e.VerifiedAt,
			})
		}
		out.Body.Memberships = toMembershipViews(res.Memberships)
		if res.CurrentOrgID != uuid.Nil {
			out.Body.CurrentOrg = &orgView{ID: res.CurrentOrgID, Slug: res.CurrentOrgSlug, Status: res.CurrentOrgStatus}
		}
		return out, nil
	})

	registerSwitchAndInvite(s, auth, requireAuth)
}

type switchOrgInput struct {
	Body struct {
		OrganizationID uuid.UUID `json:"organization_id" required:"true" format:"uuid"`
	}
}

type acceptInviteInput struct {
	Body struct {
		Token       string `json:"token" required:"true"`
		Password    string `json:"password,omitempty" minLength:"8" maxLength:"128"`
		DisplayName string `json:"display_name,omitempty" maxLength:"255"`
	}
}

func registerSwitchAndInvite(s *Server, auth *service.AuthService, requireAuth func(huma.Context, func(huma.Context))) {
	huma.Register(s.API, huma.Operation{
		OperationID: "auth-switch-org",
		Method:      http.MethodPost,
		Path:        "/auth/switch",
		Summary:     "Issue a new token pair scoped to another organization the account belongs to.",
		Tags:        []string{"auth"},
		Security:    []map[string][]string{{"bearer": {}}},
		Middlewares: huma.Middlewares{requireAuth},
	}, func(ctx context.Context, in *switchOrgInput) (*loginOutput, error) {
		claims, ok := ClaimsFromContext(ctx)
		if !ok {
			return nil, huma.Error401Unauthorized("missing claims")
		}
		res, err := auth.SwitchOrg(ctx, claims.AccountID, in.Body.OrganizationID, UserAgent(ctx), ClientIP(ctx))
		if err != nil {
			return nil, mapAuthError(err)
		}
		return authResultOutput(res), nil
	})

	huma.Register(s.API, huma.Operation{
		OperationID: "auth-accept-invite",
		Method:      http.MethodPost,
		Path:        "/auth/accept-invite",
		Summary:     "Accept an organization invite: set a password (if new) and sign in to the org.",
		Tags:        []string{"auth"},
	}, func(ctx context.Context, in *acceptInviteInput) (*loginOutput, error) {
		res, err := auth.AcceptInvite(ctx, service.AcceptInviteInput{
			Token:       in.Body.Token,
			Password:    in.Body.Password,
			DisplayName: in.Body.DisplayName,
			UserAgent:   UserAgent(ctx),
			IPAddress:   ClientIP(ctx),
		})
		if err != nil {
			return nil, mapAuthError(err)
		}
		return authResultOutput(res), nil
	})
}

// authResultOutput shapes an AuthResult into the shared login output body.
func authResultOutput(res *service.AuthResult) *loginOutput {
	out := &loginOutput{}
	acctID := res.AccountID
	out.Body.AccountID = &acctID
	out.Body.Organization = &orgView{ID: res.OrganizationID, Slug: res.OrganizationSlug, Status: res.OrgStatus}
	out.Body.Tokens = &authTokens{
		AccessToken:      res.AccessToken,
		AccessExpiresAt:  res.AccessExpiresAt,
		RefreshToken:     res.RefreshToken,
		RefreshExpiresAt: res.RefreshExpiresAt,
	}
	return out
}

func toMembershipViews(ms []service.Membership) []membershipView {
	out := make([]membershipView, 0, len(ms))
	for _, m := range ms {
		out = append(out, membershipView{
			OrganizationID:   m.OrganizationID,
			OrganizationName: m.OrganizationName,
			OrganizationSlug: m.OrganizationSlug,
		})
	}
	return out
}

func mapAuthError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return huma.Error422UnprocessableEntity(err.Error())
	case errors.Is(err, service.ErrBadCredentials):
		return huma.Error401Unauthorized(err.Error())
	case errors.Is(err, service.ErrSessionInvalid):
		return huma.Error401Unauthorized(err.Error())
	case errors.Is(err, service.ErrForbidden):
		return huma.Error403Forbidden("you don't have permission for this action")
	case errors.Is(err, service.ErrEmailExists), errors.Is(err, service.ErrDomainClaimed), errors.Is(err, service.ErrAmbiguousOrg):
		return huma.Error409Conflict(err.Error())
	default:
		return huma.Error500InternalServerError("internal error")
	}
}

