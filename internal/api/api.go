package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/glueops/waggle/internal/app"
	"github.com/glueops/waggle/internal/buildinfo"
	"github.com/glueops/waggle/internal/config"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Router http.Handler
	API    huma.API
	cfg    *config.Config
}

func Build(cfg config.Config, deps *app.Deps) (*Server, error) {
	if deps == nil {
		return nil, fmt.Errorf("api.Build: deps is nil")
	}

	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "/api/v1"
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(securityHeaders(cfg))

	humaCfg := huma.DefaultConfig("Waggle", "0.1.0")
	humaCfg.OpenAPI.Info.Description = "Waggle API."
	// Relative server URL (just the base path, no host). huma derives the
	// API prefix from this for $schema links AND uses the request host as the
	// link base; an absolute host+path here double-counts the prefix on
	// localhost (e.g. /api/v1/api/v1/schemas/...). Relative keeps it correct
	// across environments and avoids baking a host into the served spec.
	humaCfg.OpenAPI.Servers = []*huma.Server{{URL: basePath}}
	// Disable huma's built-in Stoplight docs; we serve RapiDoc at /docs instead
	// (see rapiDocHandler). OpenAPIPath stays so the spec is still published.
	humaCfg.DocsPath = ""
	humaCfg.OpenAPIPath = "/openapi"
	// Declare the "bearer" security scheme that protected operations reference
	// via Operation.Security (e.g. auth.go). Without this components entry the
	// emitted spec carries dangling security requirements, so the docs
	// "Authorize" button and generated SDK clients have no scheme to bind to.
	humaCfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	var api huma.API
	r.Route(basePath, func(sub chi.Router) {
		sub.Use(requestMeta)
		api = humachi.New(sub, humaCfg)
		sub.Get("/docs", rapiDocHandler(basePath))
	})

	s := &Server{
		Router: r,
		API:    api,
		cfg:    &cfg,
	}

	s.registerHealth()
	s.registerVersion()
	if deps.Auth != nil && deps.Tokens != nil {
		s.registerAuth(deps.Auth, deps.Tokens)
	}
	if deps.Fleet != nil && deps.Tokens != nil {
		// deps.Jobs is the concrete jobs.Enqueuer (or nil during spec gen); the
		// type assertion yields a usable enqueuer only when async discovery can
		// actually run.
		enq, _ := deps.Jobs.(discoveryEnqueuer)
		s.registerTenant(deps.Fleet, deps.Tokens, deps.APIKeys, enq)
		s.registerPool(deps.Fleet, deps.Tokens, deps.APIKeys)
	}
	if deps.APIKeys != nil && deps.Tokens != nil {
		s.registerAPIKeys(deps.APIKeys, deps.Tokens)
	}
	if deps.Orgs != nil && deps.Tokens != nil {
		s.registerOrgs(deps.Orgs, deps.Tokens)
	}

	// Serve the web UI on every path the API didn't claim (NotFound). Mounted
	// last so all API routes are registered first. See frontend.go.
	mountFrontend(r, cfg)

	return s, nil
}

func (s *Server) PublicURLs() map[string]string {
	base := strings.TrimRight(s.cfg.BaseURL, "/")
	bp := s.cfg.BasePath
	if bp == "" {
		bp = "/api/v1"
	}
	return map[string]string{
		"web":     base,
		"docs":    base + bp + "/docs",
		"openapi": base + bp + "/openapi.json",
	}
}

type healthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok"`
	}
}

func (s *Server) registerHealth() {
	huma.Register(s.API, huma.Operation{
		OperationID: "health",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health check",
		Tags:        []string{"system"},
	}, func(_ context.Context, _ *struct{}) (*healthOutput, error) {
		out := &healthOutput{}
		out.Body.Status = "ok"
		return out, nil
	})
}

type versionOutput struct {
	Body struct {
		Version string `json:"version" example:"v0.1.15"`
		Commit  string `json:"commit" example:"abc1234"`
		Date    string `json:"date" example:"2026-06-25T01:17:56Z"`
	}
}

func (s *Server) registerVersion() {
	huma.Register(s.API, huma.Operation{
		OperationID: "version",
		Method:      http.MethodGet,
		Path:        "/version",
		Summary:     "Server version",
		Tags:        []string{"system"},
	}, func(_ context.Context, _ *struct{}) (*versionOutput, error) {
		out := &versionOutput{}
		out.Body.Version = buildinfo.Version
		out.Body.Commit = buildinfo.Commit
		out.Body.Date = buildinfo.Date
		return out, nil
	})
}
