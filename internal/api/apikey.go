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

type apiKeyView struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix" doc:"Leading characters of the key, for identification."`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" format:"date-time"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" format:"date-time"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" format:"date-time"`
	CreatedAt  time.Time  `json:"created_at" format:"date-time"`
}

func toAPIKeyViewJSON(v service.APIKeyView) apiKeyView {
	return apiKeyView{
		ID:         v.ID,
		Name:       v.Name,
		Prefix:     v.Prefix,
		LastUsedAt: v.LastUsedAt,
		ExpiresAt:  v.ExpiresAt,
		RevokedAt:  v.RevokedAt,
		CreatedAt:  v.CreatedAt,
	}
}

// ---- I/O ----

type createAPIKeyInput struct {
	Body struct {
		Name string `json:"name" required:"true" maxLength:"255" doc:"Human label for the key."`
		// ExpiresInDays is optional; 0/omitted means the key never expires.
		ExpiresInDays int `json:"expires_in_days,omitempty" minimum:"0" maximum:"3650" doc:"Days until the key expires; omit or 0 for no expiry."`
	}
}

type createAPIKeyOutput struct {
	Body struct {
		// Token is the plaintext key, returned ONCE. Store it now; it cannot be
		// retrieved again.
		Token string     `json:"token" doc:"Plaintext API key — shown once, store it securely."`
		Key   apiKeyView `json:"key"`
	}
}

type listAPIKeysOutput struct {
	Body struct {
		Items []apiKeyView `json:"items"`
	}
}

type apiKeyIDInput struct {
	ID uuid.UUID `path:"id"`
}

func (s *Server) registerAPIKeys(keys *service.APIKeyService, tokens *service.TokenService) {
	if keys == nil || tokens == nil {
		return
	}
	// Management is human-only: minting/revoking keys requires a user JWT, never
	// an API key (an API key cannot bootstrap more keys).
	requireAuth := RequireAuth(s.API, tokens)
	secured := func(op huma.Operation) huma.Operation {
		op.Tags = []string{"api-keys"}
		op.Security = []map[string][]string{{"bearer": {}}}
		op.Middlewares = huma.Middlewares{requireAuth}
		return op
	}

	huma.Register(s.API, secured(huma.Operation{
		OperationID:   "create-api-key",
		Method:        http.MethodPost,
		Path:          "/api-keys",
		Summary:       "Mint an organization API key for automation (e.g. Terraform). The plaintext token is returned once.",
		DefaultStatus: http.StatusCreated,
	}), func(ctx context.Context, in *createAPIKeyInput) (*createAPIKeyOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		var createdBy *uuid.UUID
		if claims, ok := ClaimsFromContext(ctx); ok && claims.AccountID != uuid.Nil {
			id := claims.AccountID
			createdBy = &id
		}
		var expiresAt *time.Time
		if in.Body.ExpiresInDays > 0 {
			t := time.Now().UTC().AddDate(0, 0, in.Body.ExpiresInDays)
			expiresAt = &t
		}
		issued, err := keys.Issue(ctx, orgID, in.Body.Name, createdBy, expiresAt)
		if err != nil {
			return nil, mapAPIKeyError(err)
		}
		out := &createAPIKeyOutput{}
		out.Body.Token = issued.Token
		out.Body.Key = toAPIKeyViewJSON(issued.View)
		return out, nil
	})

	huma.Register(s.API, secured(huma.Operation{
		OperationID: "list-api-keys",
		Method:      http.MethodGet,
		Path:        "/api-keys",
		Summary:     "List the organization's API keys (secrets are never returned).",
	}), func(ctx context.Context, _ *struct{}) (*listAPIKeysOutput, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		views, err := keys.List(ctx, orgID)
		if err != nil {
			return nil, mapAPIKeyError(err)
		}
		out := &listAPIKeysOutput{}
		out.Body.Items = make([]apiKeyView, 0, len(views))
		for _, v := range views {
			out.Body.Items = append(out.Body.Items, toAPIKeyViewJSON(v))
		}
		return out, nil
	})

	huma.Register(s.API, secured(huma.Operation{
		OperationID:   "revoke-api-key",
		Method:        http.MethodDelete,
		Path:          "/api-keys/{id}",
		Summary:       "Revoke an organization API key. Idempotent from the caller's view.",
		DefaultStatus: http.StatusNoContent,
	}), func(ctx context.Context, in *apiKeyIDInput) (*struct{}, error) {
		orgID, err := orgFromContext(ctx)
		if err != nil {
			return nil, err
		}
		if err := keys.Revoke(ctx, orgID, in.ID); err != nil {
			return nil, mapAPIKeyError(err)
		}
		return nil, nil
	})
}

func mapAPIKeyError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return huma.Error422UnprocessableEntity(err.Error())
	case errors.Is(err, service.ErrNotFound):
		return huma.Error404NotFound("api key not found")
	default:
		return huma.Error500InternalServerError("internal error")
	}
}
