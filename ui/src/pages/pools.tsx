import { useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Plus } from "lucide-react"

import {
  listPoolsOptions,
  listPoolsQueryKey,
  createPoolMutation,
  resizePoolMutation,
  deletePoolMutation,
  listDatacentersOptions,
  listSlotsOptions,
  listPoolPlacementsOptions,
} from "@/sdk/@tanstack/react-query.gen"
import type { PoolView } from "@/sdk/types.gen"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { FormDialog } from "@/components/crud/form-dialog"
import { RowActions } from "@/components/crud/row-actions"
import { PageShell } from "@/components/crud/page-shell"
import { errMsg } from "@/lib/errors"
import { cn } from "@/lib/utils"

const selectCls =
  "border-input flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-sm shadow-xs outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]"

export function PoolsPage() {
  const qc = useQueryClient()
  const list = useQuery(listPoolsOptions())
  const datacenters = useQuery(listDatacentersOptions())
  const slots = useQuery(listSlotsOptions())

  const [creating, setCreating] = useState(false)
  const [resizing, setResizing] = useState<PoolView | null>(null)
  const [viewing, setViewing] = useState<PoolView | null>(null)
  const [form, setForm] = useState({ datacenter_id: "", slot_id: "", name: "", desired_count: 1 })
  const [desired, setDesired] = useState(1)
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)

  const invalidate = () => qc.invalidateQueries({ queryKey: listPoolsQueryKey() })
  const createMut = useMutation(createPoolMutation())
  const resizeMut = useMutation(resizePoolMutation())
  const deleteMut = useMutation(deletePoolMutation())

  const placementsQ = useQuery({
    ...listPoolPlacementsOptions({ path: { id: viewing?.id ?? "" } }),
    enabled: !!viewing,
  })
  const placements = placementsQ.data?.items ?? []

  const dcs = datacenters.data?.items ?? []
  const slotList = slots.data?.items ?? []
  const items: PoolView[] = list.data?.items ?? []

  function openCreate() {
    setForm({
      datacenter_id: dcs[0]?.id ?? "",
      slot_id: slotList[0]?.id ?? "",
      name: "",
      desired_count: 1,
    })
    setError(null)
    setCreating(true)
  }

  async function create() {
    setError(null)
    try {
      await createMut.mutateAsync({ body: form })
      await invalidate()
      setCreating(false)
    } catch (e) {
      setError(errMsg(e))
    }
  }

  async function resize() {
    if (!resizing) return
    setError(null)
    try {
      await resizeMut.mutateAsync({ body: { desired_count: desired }, path: { id: resizing.id } })
      await invalidate()
      setResizing(null)
    } catch (e) {
      setError(errMsg(e))
    }
  }

  async function remove(p: PoolView) {
    if (!confirm(`Delete pool "${p.name}" and release its placements?`)) return
    try {
      await deleteMut.mutateAsync({ path: { id: p.id } })
      await invalidate()
    } catch (e) {
      setNotice(errMsg(e))
    }
  }

  const dcName = (id: string) => dcs.find((d) => d.id === id)?.name ?? id.slice(0, 8)
  const slotName = (id: string) => slotList.find((s) => s.id === id)?.name ?? id.slice(0, 8)

  return (
    <PageShell
      title="Pools"
      description="Declarative groups of VMs placed across hypervisors with anti-affinity spread."
      action={
        <Button onClick={openCreate} disabled={dcs.length === 0 || slotList.length === 0}>
          <Plus className="size-4" /> Add pool
        </Button>
      }
      notice={notice}
      onDismissNotice={() => setNotice(null)}
      query={list}
      empty={items.length === 0}
      emptyText="No pools yet. Create one to place VMs."
    >
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Datacenter</TableHead>
            <TableHead>Slot</TableHead>
            <TableHead>Desired</TableHead>
            <TableHead className="w-10" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {items.map((p) => (
            <TableRow key={p.id}>
              <TableCell className="font-medium">{p.name}</TableCell>
              <TableCell className="text-muted-foreground">{dcName(p.datacenter_id)}</TableCell>
              <TableCell className="text-muted-foreground">{slotName(p.slot_id)}</TableCell>
              <TableCell>{p.desired_count}</TableCell>
              <TableCell>
                <RowActions
                  actions={[
                    { label: "View placements", onSelect: () => setViewing(p) },
                    {
                      label: "Resize",
                      onSelect: () => {
                        setDesired(p.desired_count)
                        setError(null)
                        setResizing(p)
                      },
                    },
                    { label: "Delete", variant: "destructive", onSelect: () => void remove(p) },
                  ]}
                />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <FormDialog
        open={creating}
        onOpenChange={setCreating}
        title="Add pool"
        description="Placement is all-or-nothing: if the fleet can't fit every VM, none are placed."
        onSubmit={() => void create()}
        submitting={createMut.isPending}
        error={error}
        submitLabel="Create & place"
      >
        <div className="flex flex-col gap-2">
          <Label htmlFor="pool-name">Name</Label>
          <Input
            id="pool-name"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            required
          />
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="pool-dc">Datacenter</Label>
          <select
            id="pool-dc"
            className={cn(selectCls)}
            value={form.datacenter_id}
            onChange={(e) => setForm({ ...form, datacenter_id: e.target.value })}
          >
            {dcs.map((d) => (
              <option key={d.id} value={d.id}>{d.name}</option>
            ))}
          </select>
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="pool-slot">Slot</Label>
          <select
            id="pool-slot"
            className={cn(selectCls)}
            value={form.slot_id}
            onChange={(e) => setForm({ ...form, slot_id: e.target.value })}
          >
            {slotList.map((s) => (
              <option key={s.id} value={s.id}>{s.name}</option>
            ))}
          </select>
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="pool-count">Desired count</Label>
          <Input
            id="pool-count"
            type="number"
            min={1}
            value={form.desired_count}
            onChange={(e) =>
              setForm({ ...form, desired_count: Math.max(1, parseInt(e.target.value || "1", 10) || 1) })
            }
          />
        </div>
      </FormDialog>

      <FormDialog
        open={resizing !== null}
        onOpenChange={(o) => !o && setResizing(null)}
        title={`Resize pool — ${resizing?.name ?? ""}`}
        description="Growing places new VMs; shrinking releases the most recent placements."
        onSubmit={() => void resize()}
        submitting={resizeMut.isPending}
        error={error}
        submitLabel="Resize"
      >
        <div className="flex flex-col gap-2">
          <Label htmlFor="resize-count">Desired count</Label>
          <Input
            id="resize-count"
            type="number"
            min={0}
            value={desired}
            onChange={(e) => setDesired(Math.max(0, parseInt(e.target.value || "0", 10) || 0))}
          />
        </div>
      </FormDialog>

      <Dialog open={viewing !== null} onOpenChange={(o) => !o && setViewing(null)}>
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>Placements — {viewing?.name}</DialogTitle>
            <DialogDescription>
              {viewing
                ? `${placements.length} of ${viewing.desired_count} desired VM(s) placed.`
                : ""}
            </DialogDescription>
          </DialogHeader>
          {placementsQ.isPending ? (
            <p className="text-muted-foreground py-6 text-center text-sm">Loading…</p>
          ) : placements.length === 0 ? (
            <p className="text-muted-foreground py-6 text-center text-sm">
              No placements for this pool.
            </p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Hypervisor</TableHead>
                  <TableHead>VMID</TableHead>
                  <TableHead>Placed</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {placements.map((pl) => (
                  <TableRow key={pl.id}>
                    <TableCell className="font-medium">{pl.hypervisor_name}</TableCell>
                    <TableCell className="tabular-nums">{pl.vmid ?? "—"}</TableCell>
                    <TableCell className="text-muted-foreground">
                      {new Date(pl.created_at).toLocaleString()}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </DialogContent>
      </Dialog>
    </PageShell>
  )
}
