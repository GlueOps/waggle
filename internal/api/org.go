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

// ---- views ----

type orgFullView struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Domain    *string   `json:"domain,omitempty"`
	Status    string    `json:"status"`
	Role      string    `json:"role" doc:"The calling account's role in this organization."`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
}

func toOrgFullView(o *service.OrgView) orgFullView {
	return orgFullView{
		ID: o.ID, Name: o.Name, Slug: o.Slug, Domain: o.Domain,
		Status: o.Status, Role: o.Role, CreatedAt: o.CreatedAt,
	}
}

type memberJSONView struct {
	UserID      uuid.UUID  `json:"user_id"`
	AccountID   uuid.UUID  `json:"account_id"`
	DisplayName string     `json:"display_name"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	IsActive    bool       `json:"is_active"`
	Pending     bool       `json:"pending" doc:"Invited but hasn't accepted (no password set yet)."`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" format:"date-time"`
	CreatedAt   time.Time  `json:"created_at" format:"date-time"`
}

func toMemberJSONView(m service.MemberView) memberJSONView {
	return memberJSONView{
		UserID: m.UserID, AccountID: m.AccountID, DisplayName: m.DisplayName, Email: m.Email,
		Role: m.Role, IsActive: m.IsActive, Pending: m.Pending, LastLoginAt: m.LastLoginAt, CreatedAt: m.CreatedAt,
	}
}

// ---- I/O ----

type createOrgInput struct {
	Body struct {
		Name string `json:"name" required:"true" maxLength:"255"`
	}
}
type orgOutput struct{ Body orgFullView }
type orgListOutput struct {
	Body struct {
		Items []orgFullView `json:"items"`
	}
}
type orgIDInput struct {
	ID uuid.UUID `path:"id"`
}
type updateOrgInput struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Name string `json:"name" required:"true" maxLength:"255"`
	}
}

type memberListOutput struct {
	Body struct {
		Items []memberJSONView `json:"items"`
	}
}
type addMemberInput struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Email string `json:"email" format:"email" required:"true" maxLength:"255"`
		Role  string `json:"role,omitempty" enum:"owner,admin,member"`
	}
}
type addMemberOutput struct {
	Body struct {
		Member  memberJSONView `json:"member"`
		Invited bool           `json:"invited" doc:"An invite email was sent (pending account)."`
	}
}
type memberIDInput struct {
	ID     uuid.UUID `path:"id"`
	UserID uuid.UUID `path:"userId"`
}
type updateMemberInput struct {
	ID     uuid.UUID `path:"id"`
	UserID uuid.UUID `path:"userId"`
	Body   struct {
		Role string `json:"role" required:"true" enum:"owner,admin,member"`
	}
}
type memberOutput struct{ Body memberJSONView }

func (s *Server) registerOrgs(orgs *service.OrgService, tokens *service.TokenService) {
	if orgs == nil || tokens == nil {
		return
	}
	requireAuth := RequireAuth(s.API, tokens) // JWT-only: org admin is a human action
	secured := func(op huma.Operation) huma.Operation {
		op.Tags = []string{"organizations"}
		op.Security = []map[string][]string{{"bearer": {}}}
		op.Middlewares = huma.Middlewares{requireAuth}
		return op
	}

	huma.Register(s.API, secured(huma.Operation{
		OperationID:   "create-org",
		Method:        http.MethodPost,
		Path:          "/organizations",
		Summary:       "Create an organization (you become its owner) and enqueue tenant provisioning.",
		DefaultStatus: http.StatusCreated,
	}), func(ctx context.Context, in *createOrgInput) (*orgOutput, error) {
		acc, err := accountFromContext(ctx)
		if err != nil {
			return nil, err
		}
		o, err := orgs.CreateOrg(ctx, acc, in.Body.Name)
		if err != nil {
			return nil, mapOrgError(err)
		}
		return &orgOutput{Body: toOrgFullView(o)}, nil
	})

	huma.Register(s.API, secured(huma.Operation{
		OperationID: "list-orgs",
		Method:      http.MethodGet,
		Path:        "/organizations",
		Summary:     "List the organizations the caller belongs to (with their role).",
	}), func(ctx context.Context, _ *struct{}) (*orgListOutput, error) {
		acc, err := accountFromContext(ctx)
		if err != nil {
			return nil, err
		}
		list, err := orgs.ListOrgs(ctx, acc)
		if err != nil {
			return nil, mapOrgError(err)
		}
		out := &orgListOutput{}
		out.Body.Items = make([]orgFullView, 0, len(list))
		for i := range list {
			out.Body.Items = append(out.Body.Items, toOrgFullView(&list[i]))
		}
		return out, nil
	})

	huma.Register(s.API, secured(huma.Operation{
		OperationID: "get-org",
		Method:      http.MethodGet,
		Path:        "/organizations/{id}",
		Summary:     "Get an organization the caller belongs to.",
	}), func(ctx context.Context, in *orgIDInput) (*orgOutput, error) {
		acc, err := accountFromContext(ctx)
		if err != nil {
			return nil, err
		}
		o, err := orgs.GetOrg(ctx, acc, in.ID)
		if err != nil {
			return nil, mapOrgError(err)
		}
		return &orgOutput{Body: toOrgFullView(o)}, nil
	})

	huma.Register(s.API, secured(huma.Operation{
		OperationID: "update-org",
		Method:      http.MethodPatch,
		Path:        "/organizations/{id}",
		Summary:     "Rename an organization (admin or owner).",
	}), func(ctx context.Context, in *updateOrgInput) (*orgOutput, error) {
		acc, err := accountFromContext(ctx)
		if err != nil {
			return nil, err
		}
		o, err := orgs.UpdateOrg(ctx, acc, in.ID, in.Body.Name)
		if err != nil {
			return nil, mapOrgError(err)
		}
		return &orgOutput{Body: toOrgFullView(o)}, nil
	})

	huma.Register(s.API, secured(huma.Operation{
		OperationID:   "delete-org",
		Method:        http.MethodDelete,
		Path:          "/organizations/{id}",
		Summary:       "Delete an organization and enqueue tenant teardown (owner only).",
		DefaultStatus: http.StatusNoContent,
	}), func(ctx context.Context, in *orgIDInput) (*struct{}, error) {
		acc, err := accountFromContext(ctx)
		if err != nil {
			return nil, err
		}
		if err := orgs.DeleteOrg(ctx, acc, in.ID); err != nil {
			return nil, mapOrgError(err)
		}
		return nil, nil
	})

	huma.Register(s.API, secured(huma.Operation{
		OperationID: "list-members",
		Method:      http.MethodGet,
		Path:        "/organizations/{id}/members",
		Summary:     "List an organization's members.",
	}), func(ctx context.Context, in *orgIDInput) (*memberListOutput, error) {
		acc, err := accountFromContext(ctx)
		if err != nil {
			return nil, err
		}
		members, err := orgs.ListMembers(ctx, acc, in.ID)
		if err != nil {
			return nil, mapOrgError(err)
		}
		out := &memberListOutput{}
		out.Body.Items = make([]memberJSONView, 0, len(members))
		for _, m := range members {
			out.Body.Items = append(out.Body.Items, toMemberJSONView(m))
		}
		return out, nil
	})

	huma.Register(s.API, secured(huma.Operation{
		OperationID:   "add-member",
		Method:        http.MethodPost,
		Path:          "/organizations/{id}/members",
		Summary:       "Add or invite a member by email (admin+; owner required to grant owner). Unknown emails get an invite link.",
		DefaultStatus: http.StatusCreated,
	}), func(ctx context.Context, in *addMemberInput) (*addMemberOutput, error) {
		acc, err := accountFromContext(ctx)
		if err != nil {
			return nil, err
		}
		res, err := orgs.AddMember(ctx, acc, in.ID, in.Body.Email, in.Body.Role)
		if err != nil {
			return nil, mapOrgError(err)
		}
		out := &addMemberOutput{}
		out.Body.Member = toMemberJSONView(res.Member)
		out.Body.Invited = res.Invited
		return out, nil
	})

	huma.Register(s.API, secured(huma.Operation{
		OperationID: "update-member",
		Method:      http.MethodPatch,
		Path:        "/organizations/{id}/members/{userId}",
		Summary:     "Change a member's role (admin+; owner required to touch owners).",
	}), func(ctx context.Context, in *updateMemberInput) (*memberOutput, error) {
		acc, err := accountFromContext(ctx)
		if err != nil {
			return nil, err
		}
		m, err := orgs.UpdateMember(ctx, acc, in.ID, in.UserID, in.Body.Role)
		if err != nil {
			return nil, mapOrgError(err)
		}
		return &memberOutput{Body: toMemberJSONView(*m)}, nil
	})

	huma.Register(s.API, secured(huma.Operation{
		OperationID:   "remove-member",
		Method:        http.MethodDelete,
		Path:          "/organizations/{id}/members/{userId}",
		Summary:       "Remove a member (admin+; owner required to remove owners; never the last owner).",
		DefaultStatus: http.StatusNoContent,
	}), func(ctx context.Context, in *memberIDInput) (*struct{}, error) {
		acc, err := accountFromContext(ctx)
		if err != nil {
			return nil, err
		}
		if err := orgs.RemoveMember(ctx, acc, in.ID, in.UserID); err != nil {
			return nil, mapOrgError(err)
		}
		return nil, nil
	})
}

// accountFromContext returns the authenticated account id from JWT claims.
func accountFromContext(ctx context.Context) (uuid.UUID, error) {
	claims, ok := ClaimsFromContext(ctx)
	if !ok || claims.AccountID == uuid.Nil {
		return uuid.Nil, huma.Error401Unauthorized("missing account context")
	}
	return claims.AccountID, nil
}

func mapOrgError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return huma.Error422UnprocessableEntity(err.Error())
	case errors.Is(err, service.ErrForbidden):
		return huma.Error403Forbidden("you don't have permission for this action")
	case errors.Is(err, service.ErrAlreadyMember):
		return huma.Error409Conflict("that account is already a member")
	case errors.Is(err, service.ErrNotFound):
		return huma.Error404NotFound("not found")
	default:
		return huma.Error500InternalServerError("internal error")
	}
}
