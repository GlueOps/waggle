// Package ui embeds the built frontend (Vite output in dist/) into the binary
// so a single Waggle executable can serve its own web UI. The dist/ directory
// is produced by `yarn build` (or `just ui`); a placeholder keeps the embed
// directive valid before the first build, in which case Dist reports ok=false.
package ui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distFS embed.FS

// Dist returns the built frontend rooted at the dist directory. ok is false
// when no real build is present (only the placeholder), letting callers fall
// back to "none" behaviour instead of serving a blank page.
func Dist() (assets fs.FS, ok bool) {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil, false
	}
	if _, err := fs.Stat(sub, "index.html"); err != nil {
		return nil, false
	}
	return sub, true
}
