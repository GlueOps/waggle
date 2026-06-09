import { useEffect, useState } from "react"

import { authVerifyEmail } from "@/sdk/sdk.gen"
import { Sidebar, type NavKey } from "@/components/layout/sidebar"
import { Topbar } from "@/components/layout/topbar"
import { DashboardPage } from "@/pages/dashboard"
import { DatacentersPage } from "@/pages/datacenters"
import { SlotsPage } from "@/pages/slots"
import { HypervisorsPage } from "@/pages/hypervisors"
import { PoolsPage } from "@/pages/pools"
import { PlacementsPage } from "@/pages/placements"
import { OrganizationsPage } from "@/pages/organizations"
import { ApiKeysPage } from "@/pages/api-keys"
import { LoginPage } from "@/pages/login"
import { AcceptInvitePage } from "@/pages/accept-invite"
import { useAuth } from "@/lib/auth"

const TITLES: Record<NavKey, string> = {
  dashboard: "Dashboard",
  datacenters: "Datacenters",
  slots: "Slots",
  hypervisors: "Hypervisors",
  pools: "Pools",
  placements: "Placements",
  organizations: "Organizations",
  "api-keys": "API keys",
}

function queryParam(name: string): string | null {
  return new URLSearchParams(window.location.search).get(name)
}

function App() {
  const { session } = useAuth()
  const [active, setActive] = useState<NavKey>("dashboard")
  const [verifyMsg, setVerifyMsg] = useState<string | null>(null)

  // Consume a ?verify=<token> email link once on load (independent of session).
  useEffect(() => {
    const token = queryParam("verify")
    if (!token) return
    authVerifyEmail({ body: { token } })
      .then(({ error }) =>
        setVerifyMsg(error ? "Verification link is invalid or expired." : "Email verified — you can sign in.")
      )
      .catch(() => setVerifyMsg("Verification link is invalid or expired."))
      .finally(() => window.history.replaceState(null, "", window.location.pathname))
  }, [])

  // An invite link (?invite=...) takes precedence when not signed in.
  const invite = queryParam("invite")
  if (invite && !session) {
    return <AcceptInvitePage token={invite} />
  }
  if (!session) {
    return <LoginPage notice={verifyMsg} />
  }

  return (
    <div className="bg-background text-foreground flex min-h-svh">
      <Sidebar active={active} onNavigate={setActive} />
      <div className="flex min-w-0 flex-1 flex-col">
        <Topbar title={TITLES[active]} />
        <main className="flex-1">
          {active === "dashboard" && <DashboardPage />}
          {active === "datacenters" && <DatacentersPage />}
          {active === "slots" && <SlotsPage />}
          {active === "hypervisors" && <HypervisorsPage />}
          {active === "pools" && <PoolsPage />}
          {active === "placements" && <PlacementsPage />}
          {active === "organizations" && <OrganizationsPage />}
          {active === "api-keys" && <ApiKeysPage />}
        </main>
      </div>
    </div>
  )
}

export default App
