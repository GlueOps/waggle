package api

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/glueops/waggle/internal/config"
	"github.com/glueops/waggle/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

// apiCSP is the strict default applied to every response. API responses are
// JSON and load no sub-resources, so locking everything down is safe. The
// docs route (HTML) overrides this with docsCSP (see docs.go).
const apiCSP = "default-src 'none'; frame-ancestors 'none'; base-uri 'none'"

// securityHeaders applies baseline hardening headers to all responses.
// HSTS is only sent when the connection is TLS or the app runs in production
// (common when a TLS-terminating proxy speaks plain HTTP to the app).
func securityHeaders(cfg config.Config) func(http.Handler) http.Handler {
	prod := strings.EqualFold(cfg.Env, "production")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			h.Set("Content-Security-Policy", apiCSP)
			if prod || r.TLS != nil {
				h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
			next.ServeHTTP(w, r)
		})
	}
}

type ctxKey int

const (
	ctxKeyClaims ctxKey = iota
	ctxKeyClientIP
	ctxKeyUserAgent
	ctxKeyPrincipal
)

// Principal kinds recorded by the auth middleware so handlers/audit can tell a
// human session apart from an automation API key.
type Principal string

const (
	PrincipalUser   Principal = "user"
	PrincipalAPIKey Principal = "api_key"
)

// requestMeta is a chi middleware that records the client's TCP peer and
// User-Agent on the request context. Public auth endpoints (signup/login)
// need these for TokenSession + audit even when no bearer is present.
func requestMeta(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.RemoteAddr
		if h, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			host = h
		}
		ctx := context.WithValue(r.Context(), ctxKeyClientIP, host)
		ctx = context.WithValue(ctx, ctxKeyUserAgent, r.Header.Get("User-Agent"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth returns a huma middleware that validates a user Bearer JWT and
// stores the parsed Claims in the request context. Operations opt in via
// Operation.Middlewares — public ops (signup/login/refresh/logout/verify-
// email/health) do not include this middleware. Use this for account-scoped
// endpoints that only a human session may reach (e.g. /auth/me, API-key mgmt).
func RequireAuth(api huma.API, tokens *service.TokenService) func(huma.Context, func(huma.Context)) {
	return RequireBearer(api, tokens, nil)
}

// RequireBearer validates the Authorization: Bearer credential, accepting
// EITHER a user JWT or — when keys is non-nil — an organization API key
// (distinguished by its "wgl_" prefix). Both resolve to an org context via the
// stored Claims, so tenant endpoints work identically for the web app and the
// Terraform provider. The resolved Principal is recorded for audit.
func RequireBearer(api huma.API, tokens *service.TokenService, keys *service.APIKeyService) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		header := ctx.Header("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing bearer token")
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))

		// API key path: synthesize org-only Claims so orgFromContext and the
		// fleet handlers behave the same as for a user session.
		if keys != nil && service.LooksLikeAPIKey(raw) {
			authed, err := keys.Authenticate(ctx.Context(), raw)
			if err != nil {
				_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid or expired api key")
				return
			}
			c := huma.WithValue(ctx, ctxKeyClaims, &service.Claims{OrganizationID: authed.OrganizationID})
			c = huma.WithValue(c, ctxKeyPrincipal, PrincipalAPIKey)
			next(c)
			return
		}

		claims, err := tokens.Verify(raw)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		c := huma.WithValue(ctx, ctxKeyClaims, claims)
		c = huma.WithValue(c, ctxKeyPrincipal, PrincipalUser)
		next(c)
	}
}

// PrincipalFromContext returns the authenticated principal kind recorded by the
// auth middleware, or false if the request was not authenticated.
func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(ctxKeyPrincipal).(Principal)
	return p, ok
}

// ClaimsFromContext retrieves the verified JWT claims placed by RequireAuth.
// Returns false on a public operation that did not run the middleware.
func ClaimsFromContext(ctx context.Context) (*service.Claims, bool) {
	c, ok := ctx.Value(ctxKeyClaims).(*service.Claims)
	return c, ok
}

// ClientIP returns the TCP peer IP recorded by requestMeta. Falls back to
// empty string if requestMeta did not run (shouldn't happen in production).
func ClientIP(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyClientIP).(string); ok {
		return v
	}
	return ""
}

// UserAgent returns the User-Agent header recorded by requestMeta.
func UserAgent(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyUserAgent).(string); ok {
		return v
	}
	return ""
}
