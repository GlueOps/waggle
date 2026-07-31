import { writeFileSync } from 'node:fs'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig, type Plugin } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// dist/.gitkeep is tracked in git so that the `//go:embed all:dist` directive in
// ui/embed.go stays valid on a fresh clone — without it `go build ./...` and
// `go test ./...` fail before the frontend has ever been built. Vite empties
// outDir on every build, which deletes the placeholder and leaves it showing up
// as a spurious deletion in `git status`, so put it back afterwards.
function keepEmbedPlaceholder(): Plugin {
  return {
    name: 'waggle:keep-embed-placeholder',
    apply: 'build',
    closeBundle() {
      writeFileSync(
        fileURLToPath(new URL('./dist/.gitkeep', import.meta.url)),
        'Placeholder that keeps `//go:embed all:dist` in ui/embed.go valid.\n' +
          'Restored automatically by the waggle:keep-embed-placeholder Vite plugin.\n',
      )
    },
  }
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss(), keepEmbedPlaceholder()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    // In FRONTEND_MODE=proxy the Go server reverse-proxies here, so API calls
    // are same-origin from the browser's perspective. This proxy is only for
    // running `vite` standalone (visiting :5173 directly) during UI-only work.
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
