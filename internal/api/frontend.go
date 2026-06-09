package api

import (
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/glueops/waggle/internal/config"
	"github.com/glueops/waggle/ui"

	"github.com/go-chi/chi/v5"
)

// webCSP relaxes the strict apiCSP (default-src 'none') for the SPA shell so the
// browser may load the hashed JS/CSS bundles and call the same-origin API.
// 'unsafe-inline' on style-src covers runtime style injection by UI libraries;
// scripts stay 'self'-only since Vite emits external module files, not inline.
const webCSP = "default-src 'self'; " +
	"script-src 'self'; " +
	"style-src 'self' 'unsafe-inline'; " +
	"img-src 'self' data:; " +
	"font-src 'self' data:; " +
	"connect-src 'self'; " +
	"frame-ancestors 'none'; " +
	"base-uri 'self'"

// devCSP is permissive enough for the Vite dev server proxied in FRONTEND_MODE
// =proxy: HMR needs inline scripts/eval and a websocket back to the dev server.
// Never used for embedded production assets.
const devCSP = "default-src 'self' 'unsafe-inline' 'unsafe-eval' data: blob: ws: wss:"

const defaultViteDevURL = "http://localhost:5173"

// mountFrontend wires the web UI onto the router's NotFound handler so it owns
// every path the API does not. Resolution by cfg.FrontendMode:
//
//	embed  (default) – serve the assets baked into the binary via go:embed
//	proxy            – reverse-proxy to the Vite dev server (cfg.ViteDevURL)
//	none             – leave the default 404 in place (API-only)
func mountFrontend(r chi.Router, cfg config.Config) {
	mode := cfg.FrontendMode
	if mode == "" {
		mode = config.FrontendModeEmbed
	}

	switch mode {
	case config.FrontendModeNone:
		return

	case config.FrontendModeProxy:
		target := cfg.ViteDevURL
		if strings.TrimSpace(target) == "" {
			target = defaultViteDevURL
		}
		u, err := url.Parse(target)
		if err != nil {
			log.Printf("frontend: invalid vite_dev_url %q: %v (frontend disabled)", target, err)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(u)
		r.NotFound(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Security-Policy", devCSP)
			proxy.ServeHTTP(w, req)
		})
		log.Printf("frontend: proxying to %s", target)

	case config.FrontendModeEmbed:
		assets, ok := ui.Dist()
		if !ok {
			log.Printf("frontend: no embedded build found (run `just ui`); serving API only")
			return
		}
		r.NotFound(spaHandler(assets))

	default:
		log.Printf("frontend: unknown mode %q; serving API only", mode)
	}
}

// spaHandler serves static files from the embedded build and falls back to
// index.html for any unknown path so client-side routing works on deep links.
func spaHandler(assets fs.FS) http.HandlerFunc {
	fileServer := http.FileServer(http.FS(assets))
	index, err := fs.ReadFile(assets, "index.html")
	if err != nil {
		log.Printf("frontend: cannot read embedded index.html: %v", err)
	}

	writeIndex := func(w http.ResponseWriter) {
		w.Header().Set("Content-Security-Policy", webCSP)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write(index)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		upath := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
		if upath == "" {
			writeIndex(w)
			return
		}

		f, ferr := assets.Open(upath)
		if ferr != nil {
			// Unknown path → hand it to the SPA router.
			writeIndex(w)
			return
		}
		stat, serr := f.Stat()
		_ = f.Close()
		if serr != nil || stat.IsDir() {
			writeIndex(w)
			return
		}

		w.Header().Set("Content-Security-Policy", webCSP)
		// Vite emits content-hashed filenames under assets/, safe to cache hard.
		if strings.HasPrefix(upath, "assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		fileServer.ServeHTTP(w, r)
	}
}
