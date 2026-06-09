import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClientProvider } from '@tanstack/react-query'
import './index.css'
import '@/lib/api' // side effect: configure the SDK client (baseUrl + bearer auth)
import App from './App.tsx'
import { queryClient } from './lib/query'
import { ThemeProvider } from '@/components/theme/theme-provider'
import { AuthProvider } from '@/lib/auth'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <AuthProvider>
          <App />
        </AuthProvider>
      </QueryClientProvider>
    </ThemeProvider>
  </StrictMode>,
)
