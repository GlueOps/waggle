import { createContext, useContext, useState } from "react"

import { authLogin, authLogout, authSwitchOrg, authAcceptInvite } from "@/sdk/sdk.gen"
import { getAccessToken, setAccessToken } from "@/lib/api"
import { queryClient } from "@/lib/query"

type OrgInfo = { id: string; slug: string; status: string }

type Membership = {
  organization_id: string
  organization_name: string
  organization_slug: string
}

type Session = {
  accessToken: string
  refreshToken?: string
  org?: OrgInfo
}

const ORG_KEY = "waggle-org"
const REFRESH_KEY = "waggle-refresh-token"

type LoginResult =
  | { ok: true }
  | { ok: false; ambiguous: Membership[] }
  | { ok: false; error: string }

type AuthState = {
  session: Session | null
  login: (email: string, password: string, organizationId?: string) => Promise<LoginResult>
  logout: () => Promise<void>
  switchOrg: (organizationId: string) => Promise<{ ok: true } | { ok: false; error: string }>
  acceptInvite: (
    token: string,
    password: string,
    displayName?: string
  ) => Promise<{ ok: true } | { ok: false; error: string }>
}

const AuthContext = createContext<AuthState>({
  session: null,
  login: async () => ({ ok: false, error: "not ready" }),
  logout: async () => {},
  switchOrg: async () => ({ ok: false, error: "not ready" }),
  acceptInvite: async () => ({ ok: false, error: "not ready" }),
})

function loadInitialSession(): Session | null {
  const token = getAccessToken()
  if (!token) return null
  const orgRaw = localStorage.getItem(ORG_KEY)
  return {
    accessToken: token,
    refreshToken: localStorage.getItem(REFRESH_KEY) ?? undefined,
    org: orgRaw ? (JSON.parse(orgRaw) as OrgInfo) : undefined,
  }
}

// Shape of the shared login/switch/accept response body.
type AuthBody = {
  account_id?: string
  organization?: { id: string; slug: string; status: string }
  tokens?: { access_token: string; refresh_token?: string }
  memberships?: Membership[]
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [session, setSession] = useState<Session | null>(loadInitialSession)

  // Persist an auth response and become that session. Clears cached tenant
  // data so queries refetch under the new org context.
  function applyAuth(data: AuthBody) {
    const token = data.tokens?.access_token
    if (!token) return false
    const org = data.organization
      ? { id: data.organization.id, slug: data.organization.slug, status: data.organization.status }
      : undefined
    setAccessToken(token)
    if (data.tokens?.refresh_token) localStorage.setItem(REFRESH_KEY, data.tokens.refresh_token)
    if (org) localStorage.setItem(ORG_KEY, JSON.stringify(org))
    queryClient.clear()
    setSession({ accessToken: token, refreshToken: data.tokens?.refresh_token, org })
    return true
  }

  const login: AuthState["login"] = async (email, password, organizationId) => {
    const { data, error } = await authLogin({ body: { email, password, organization_id: organizationId } })
    if (error) return { ok: false, error: extractError(error) }
    if (data?.memberships && data.memberships.length > 0 && !data.tokens) {
      return { ok: false, ambiguous: data.memberships }
    }
    return applyAuth(data as AuthBody) ? { ok: true } : { ok: false, error: "login failed: no token" }
  }

  const switchOrg: AuthState["switchOrg"] = async (organizationId) => {
    const { data, error } = await authSwitchOrg({ body: { organization_id: organizationId } })
    if (error) return { ok: false, error: extractError(error) }
    return applyAuth(data as AuthBody) ? { ok: true } : { ok: false, error: "switch failed" }
  }

  const acceptInvite: AuthState["acceptInvite"] = async (token, password, displayName) => {
    const { data, error } = await authAcceptInvite({ body: { token, password, display_name: displayName } })
    if (error) return { ok: false, error: extractError(error) }
    return applyAuth(data as AuthBody) ? { ok: true } : { ok: false, error: "accept failed" }
  }

  const logout: AuthState["logout"] = async () => {
    const refresh = session?.refreshToken
    if (refresh) {
      try {
        await authLogout({ body: { refresh_token: refresh } })
      } catch {
        /* best-effort */
      }
    }
    setAccessToken(null)
    localStorage.removeItem(REFRESH_KEY)
    localStorage.removeItem(ORG_KEY)
    queryClient.clear()
    setSession(null)
  }

  return (
    <AuthContext.Provider value={{ session, login, logout, switchOrg, acceptInvite }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}

function extractError(error: unknown): string {
  if (error && typeof error === "object") {
    const e = error as { detail?: string; title?: string; message?: string }
    return e.detail || e.title || e.message || "request failed"
  }
  return "request failed"
}
