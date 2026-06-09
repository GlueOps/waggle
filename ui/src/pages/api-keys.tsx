import { useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Copy, Plus } from "lucide-react"

import {
  listApiKeysOptions,
  listApiKeysQueryKey,
  createApiKeyMutation,
  revokeApiKeyMutation,
} from "@/sdk/@tanstack/react-query.gen"
import type { ApiKeyView } from "@/sdk/types.gen"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { FormDialog } from "@/components/crud/form-dialog"
import { RowActions } from "@/components/crud/row-actions"
import { PageShell } from "@/components/crud/page-shell"
import { errMsg } from "@/lib/errors"

function status(k: ApiKeyView): { label: string; variant: "default" | "secondary" | "destructive" } {
  if (k.revoked_at) return { label: "revoked", variant: "destructive" }
  if (k.expires_at && new Date(k.expires_at) < new Date())
    return { label: "expired", variant: "secondary" }
  return { label: "active", variant: "default" }
}

export function ApiKeysPage() {
  const qc = useQueryClient()
  const list = useQuery(listApiKeysOptions())
  const [creating, setCreating] = useState(false)
  const [name, setName] = useState("")
  const [expiresDays, setExpiresDays] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)
  const [newToken, setNewToken] = useState<string | null>(null)

  const invalidate = () => qc.invalidateQueries({ queryKey: listApiKeysQueryKey() })
  const createMut = useMutation(createApiKeyMutation())
  const revokeMut = useMutation(revokeApiKeyMutation())

  async function create() {
    setError(null)
    try {
      const res = await createMut.mutateAsync({
        body: { name, ...(expiresDays > 0 ? { expires_in_days: expiresDays } : {}) },
      })
      await invalidate()
      setCreating(false)
      setNewToken(res?.token ?? null)
    } catch (e) {
      setError(errMsg(e))
    }
  }

  async function revoke(k: ApiKeyView) {
    if (!confirm(`Revoke key "${k.name}"? Anything using it will stop working.`)) return
    try {
      await revokeMut.mutateAsync({ path: { id: k.id } })
      await invalidate()
    } catch (e) {
      setNotice(errMsg(e))
    }
  }

  const items = list.data?.items ?? []

  return (
    <PageShell
      title="API keys"
      description="Organization tokens for automation (e.g. the Terraform provider). Treat them like passwords."
      action={
        <Button
          onClick={() => {
            setName("")
            setExpiresDays(0)
            setError(null)
            setCreating(true)
          }}
        >
          <Plus className="size-4" /> Create key
        </Button>
      }
      notice={notice}
      onDismissNotice={() => setNotice(null)}
      query={list}
      empty={items.length === 0}
      emptyText="No API keys yet."
    >
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Prefix</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Last used</TableHead>
            <TableHead>Expires</TableHead>
            <TableHead className="w-10" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {items.map((k) => {
            const s = status(k)
            return (
              <TableRow key={k.id}>
                <TableCell className="font-medium">{k.name}</TableCell>
                <TableCell className="font-mono text-xs">{k.prefix}…</TableCell>
                <TableCell>
                  <Badge variant={s.variant}>{s.label}</Badge>
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {k.last_used_at ? new Date(k.last_used_at).toLocaleString() : "never"}
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {k.expires_at ? new Date(k.expires_at).toLocaleDateString() : "never"}
                </TableCell>
                <TableCell>
                  {!k.revoked_at && (
                    <RowActions
                      actions={[
                        { label: "Revoke", variant: "destructive", onSelect: () => void revoke(k) },
                      ]}
                    />
                  )}
                </TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>

      <FormDialog
        open={creating}
        onOpenChange={setCreating}
        title="Create API key"
        description="The token is shown once after creation — store it securely."
        onSubmit={() => void create()}
        submitting={createMut.isPending}
        error={error}
        submitLabel="Create"
      >
        <div className="flex flex-col gap-2">
          <Label htmlFor="key-name">Name</Label>
          <Input
            id="key-name"
            placeholder="terraform-prod"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="key-exp">Expires in (days, 0 = never)</Label>
          <Input
            id="key-exp"
            type="number"
            min={0}
            value={expiresDays}
            onChange={(e) => setExpiresDays(Math.max(0, parseInt(e.target.value || "0", 10) || 0))}
          />
        </div>
      </FormDialog>

      <Dialog open={newToken !== null} onOpenChange={(o) => !o && setNewToken(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Copy your API key</DialogTitle>
            <DialogDescription>
              This is the only time the token is shown. Store it now — you can't retrieve it later.
            </DialogDescription>
          </DialogHeader>
          <div className="flex items-center gap-2">
            <code className="bg-muted flex-1 overflow-x-auto rounded-md px-3 py-2 font-mono text-xs">
              {newToken}
            </code>
            <Button
              variant="outline"
              size="icon"
              onClick={() => newToken && navigator.clipboard?.writeText(newToken)}
              aria-label="Copy"
            >
              <Copy className="size-4" />
            </Button>
          </div>
          <DialogFooter>
            <Button onClick={() => setNewToken(null)}>Done</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </PageShell>
  )
}
