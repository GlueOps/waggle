import { useState } from "react"

import { useAuth } from "@/lib/auth"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

type Membership = {
  organization_id: string
  organization_name: string
}

export function LoginPage({ notice }: { notice?: string | null }) {
  const { login } = useAuth()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [orgChoices, setOrgChoices] = useState<Membership[] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function submit(organizationId?: string) {
    setBusy(true)
    setError(null)
    const res = await login(email, password, organizationId)
    setBusy(false)
    if (res.ok) return
    if ("ambiguous" in res) {
      setOrgChoices(res.ambiguous)
      return
    }
    setError(res.error)
  }

  return (
    <div className="bg-background flex min-h-svh items-center justify-center p-4">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <div className="bg-primary text-primary-foreground mb-2 flex size-9 items-center justify-center rounded-md font-semibold">
            W
          </div>
          <CardTitle>Sign in to Waggle</CardTitle>
          <CardDescription>Placement oracle &amp; ledger for Proxmox</CardDescription>
        </CardHeader>
        <CardContent>
          {notice && (
            <p className="bg-muted text-foreground mb-4 rounded-md border px-3 py-2 text-sm">
              {notice}
            </p>
          )}
          {orgChoices ? (
            <div className="flex flex-col gap-3">
              <p className="text-muted-foreground text-sm">
                This account belongs to multiple organizations. Choose one:
              </p>
              {orgChoices.map((m) => (
                <Button
                  key={m.organization_id}
                  variant="outline"
                  disabled={busy}
                  onClick={() => submit(m.organization_id)}
                >
                  {m.organization_name}
                </Button>
              ))}
              <Button variant="ghost" onClick={() => setOrgChoices(null)}>
                Back
              </Button>
            </div>
          ) : (
            <form
              className="flex flex-col gap-4"
              onSubmit={(e) => {
                e.preventDefault()
                void submit()
              }}
            >
              <div className="flex flex-col gap-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  autoComplete="username"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </div>
              <div className="flex flex-col gap-2">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  autoComplete="current-password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>
              {error && <p className="text-destructive text-sm">{error}</p>}
              <Button type="submit" disabled={busy}>
                {busy ? "Signing in…" : "Sign in"}
              </Button>
            </form>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
