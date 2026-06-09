import { useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Plus } from "lucide-react"

import {
  listDatacentersOptions,
  listDatacentersQueryKey,
  createDatacenterMutation,
  updateDatacenterMutation,
  deleteDatacenterMutation,
  discoverHypervisorsMutation,
  listHypervisorsQueryKey,
} from "@/sdk/@tanstack/react-query.gen"
import type { DatacenterView } from "@/sdk/types.gen"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
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

type Editing = { mode: "create" } | { mode: "edit"; dc: DatacenterView } | null

export function DatacentersPage() {
  const qc = useQueryClient()
  const list = useQuery(listDatacentersOptions())
  const [editing, setEditing] = useState<Editing>(null)
  const [name, setName] = useState("")
  const [url, setUrl] = useState("")
  const [token, setToken] = useState("")
  const [insecure, setInsecure] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)

  const invalidate = () =>
    qc.invalidateQueries({ queryKey: listDatacentersQueryKey() })

  const createMut = useMutation(createDatacenterMutation())
  const updateMut = useMutation(updateDatacenterMutation())
  const deleteMut = useMutation(deleteDatacenterMutation())
  const discoverMut = useMutation(discoverHypervisorsMutation())

  function openCreate() {
    setName("")
    setUrl("")
    setToken("")
    setInsecure(false)
    setError(null)
    setEditing({ mode: "create" })
  }
  function openEdit(dc: DatacenterView) {
    setName(dc.name)
    setUrl(dc.url)
    setToken("")
    setInsecure(dc.insecure_skip_verify)
    setError(null)
    setEditing({ mode: "edit", dc })
  }

  async function save() {
    setError(null)
    const body = { name, url, insecure_skip_verify: insecure, ...(token ? { token } : {}) }
    try {
      if (editing?.mode === "edit") {
        await updateMut.mutateAsync({ body, path: { id: editing.dc.id } })
      } else {
        await createMut.mutateAsync({ body })
      }
      await invalidate()
      setEditing(null)
    } catch (e) {
      setError(errMsg(e))
    }
  }

  async function remove(dc: DatacenterView) {
    if (!confirm(`Delete datacenter "${dc.name}"?`)) return
    try {
      await deleteMut.mutateAsync({ path: { id: dc.id } })
      await invalidate()
    } catch (e) {
      setNotice(errMsg(e))
    }
  }

  async function discover(dc: DatacenterView) {
    setNotice(null)
    try {
      const res = await discoverMut.mutateAsync({
        body: { async: false },
        path: { id: dc.id },
      })
      await qc.invalidateQueries({ queryKey: listHypervisorsQueryKey() })
      setNotice(`Discovered ${res?.items?.length ?? 0} hypervisor(s) from ${dc.name}.`)
    } catch (e) {
      setNotice(errMsg(e))
    }
  }

  const items = list.data?.items ?? []

  return (
    <PageShell
      title="Datacenters"
      description="Proxmox clusters this tenant can place workloads onto."
      action={
        <Button onClick={openCreate}>
          <Plus className="size-4" /> Add datacenter
        </Button>
      }
      notice={notice}
      onDismissNotice={() => setNotice(null)}
      query={list}
      empty={items.length === 0}
      emptyText="No datacenters yet. Add one to begin."
    >
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>URL</TableHead>
            <TableHead>Token</TableHead>
            <TableHead className="w-10" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {items.map((dc) => (
            <TableRow key={dc.id}>
              <TableCell className="font-medium">{dc.name}</TableCell>
              <TableCell className="text-muted-foreground">
                {dc.url}
                {dc.insecure_skip_verify && (
                  <Badge variant="outline" className="ml-2">insecure TLS</Badge>
                )}
              </TableCell>
              <TableCell>
                {dc.has_token ? (
                  <Badge variant="secondary">configured</Badge>
                ) : (
                  <Badge variant="outline">none</Badge>
                )}
              </TableCell>
              <TableCell>
                <RowActions
                  actions={[
                    { label: "Edit", onSelect: () => openEdit(dc) },
                    {
                      label: discoverMut.isPending ? "Discovering…" : "Discover hypervisors",
                      onSelect: () => void discover(dc),
                    },
                    { label: "Delete", variant: "destructive", onSelect: () => void remove(dc) },
                  ]}
                />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <FormDialog
        open={editing !== null}
        onOpenChange={(o) => !o && setEditing(null)}
        title={editing?.mode === "edit" ? "Edit datacenter" : "Add datacenter"}
        description="Point Waggle at a Proxmox cluster. The API token is encrypted at rest and used only for read-only discovery."
        onSubmit={() => void save()}
        submitting={createMut.isPending || updateMut.isPending}
        error={error}
      >
        <div className="flex flex-col gap-2">
          <Label htmlFor="dc-name">Name</Label>
          <Input id="dc-name" value={name} onChange={(e) => setName(e.target.value)} required />
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="dc-url">Proxmox URL</Label>
          <Input
            id="dc-url"
            placeholder="https://pve.example.com:8006"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            required
          />
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="dc-token">
            API token {editing?.mode === "edit" && "(leave blank to keep current)"}
          </Label>
          <Input
            id="dc-token"
            type="password"
            placeholder="user@realm!tokenid=secret"
            value={token}
            onChange={(e) => setToken(e.target.value)}
          />
        </div>
        <div className="flex items-center justify-between gap-4 rounded-md border p-3">
          <div>
            <Label htmlFor="dc-insecure">Allow self-signed TLS</Label>
            <p className="text-muted-foreground text-xs">
              Skip certificate verification (common for homelab Proxmox).
            </p>
          </div>
          <Switch id="dc-insecure" checked={insecure} onCheckedChange={setInsecure} />
        </div>
      </FormDialog>
    </PageShell>
  )
}
