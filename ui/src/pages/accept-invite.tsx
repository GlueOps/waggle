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

export function AcceptInvitePage({ token }: { token: string }) {
  const { acceptInvite } = useAuth()
  const [password, setPassword] = useState("")
  const [displayName, setDisplayName] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function submit() {
    setBusy(true)
    setError(null)
    const res = await acceptInvite(token, password, displayName || undefined)
    setBusy(false)
    if (!res.ok) {
      setError(res.error)
      return
    }
    // Drop the invite token from the URL; the app re-renders into the session.
    window.history.replaceState(null, "", window.location.pathname)
  }

  return (
    <div className="bg-background flex min-h-svh items-center justify-center p-4">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Accept your invitation</CardTitle>
          <CardDescription>Set a password to join the organization.</CardDescription>
        </CardHeader>
        <CardContent>
          <form
            className="flex flex-col gap-4"
            onSubmit={(e) => {
              e.preventDefault()
              void submit()
            }}
          >
            <div className="flex flex-col gap-2">
              <Label htmlFor="inv-name">Display name (optional)</Label>
              <Input id="inv-name" value={displayName} onChange={(e) => setDisplayName(e.target.value)} />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="inv-pass">Password</Label>
              <Input
                id="inv-pass"
                type="password"
                minLength={8}
                autoComplete="new-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            {error && <p className="text-destructive text-sm">{error}</p>}
            <Button type="submit" disabled={busy}>
              {busy ? "Joining…" : "Join organization"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
