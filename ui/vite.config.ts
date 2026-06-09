import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
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
