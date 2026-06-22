import { useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Plus } from "lucide-react"

import {
  listPlacementsOptions,
  listPlacementsQueryKey,
  backfillPlacementVmidMutation,
  deletePlacementMutation,
} from "@/sdk/@tanstack/react-query.gen"
import type { FleetPlacementView } from "@/sdk/types.gen"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { PageShell } from "@/components/crud/page-shell"
import { RowActions } from "@/components/crud/row-actions"
import { FormDialog } from "@/components/crud/form-dialog"
import { CopyId } from "@/components/ui/copy-id"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { errMsg } from "@/lib/errors"

export function PlacementsPage() {
  const qc = useQueryClient()
  const list = useQuery(listPlacementsOptions())
  const items = list.data?.items ?? []

  const [backfilling, setBackfilling] = useState<FleetPlacementView | null>(null)
  const [assignOpen, setAssignOpen] = useState(false)
  const [assignTarget, setAssignTarget] = useState<string>("")
  const [vmidInput, setVmidInput] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)

  const invalidate = () => qc.invalidateQueries({ queryKey: listPlacementsQueryKey() })
  const backfillMut = useMutation(backfillPlacementVmidMutation())
  const deleteMut = useMutation(deletePlacementMutation())

  // Unassigned placements available for top-level assign dialog
  const unassigned = items.filter((p) => p.vmid == null)

  function openAssign() {
    setAssignTarget(unassigned[0]?.id ?? "")
    setVmidInput("")
    setError(null)
    setAssignOpen(true)
  }

  async function submitAssign() {
    const vmid = parseInt(vmidInput, 10)
    if (!assignTarget) { setError("Select a placement"); return }
    if (!vmid || vmid < 1) { setError("vmid must be a positive integer"); return }
    setError(null)
    try {
      await backfillMut.mutateAsync({ body: { vmid }, path: { id: assignTarget } })
      await invalidate()
      setAssignOpen(false)
    } catch (e) {
      setError(errMsg(e))
    }
  }

  async function submitBackfill() {
    if (!backfilling) return
    const vmid = parseInt(vmidInput, 10)
    if (!vmid || vmid < 1) {
      setError("vmid must be a positive integer")
      return
    }
    setError(null)
    try {
      await backfillMut.mutateAsync({ body: { vmid }, path: { id: backfilling.id } })
      await invalidate()
      setBackfilling(null)
      setVmidInput("")
    } catch (e) {
      setError(errMsg(e))
    }
  }

  async function remove(p: FleetPlacementView) {
    if (!confirm(`Delete placement ${p.id.slice(0, 8)}… on ${p.hypervisor_name}? The pool's desired count is NOT adjusted.`)) return
    try {
      await deleteMut.mutateAsync({ path: { id: p.id } })
      await invalidate()
    } catch (e) {
      setNotice(errMsg(e))
    }
  }

  function openBackfill(p: FleetPlacementView) {
    setVmidInput(p.vmid?.toString() ?? "")
    setError(null)
    setBackfilling(p)
  }

  // Per-hypervisor rollup of committed placement resources.
  const byHv = new Map<string, { count: number; vcpu: number; ram: number; disk: number }>()
  for (const p of items) {
    const cur = byHv.get(p.hypervisor_name) ?? { count: 0, vcpu: 0, ram: 0, disk: 0 }
    cur.count += 1
    cur.vcpu += p.vcpu
    cur.ram += p.ram_gb
    cur.disk += p.disk_gb
    byHv.set(p.hypervisor_name, cur)
  }
  const rollup = [...byHv.entries()].sort((a, b) => a[0].localeCompare(b[0]))

  return (
    <PageShell
      title="Placements"
      description="Every VM the oracle has placed across the fleet, and what each hypervisor is carrying."
      query={list}
      notice={notice}
      onDismissNotice={() => setNotice(null)}
      empty={items.length === 0}
      emptyText="No placements yet. Create a pool to place VMs."
      action={
        <Button onClick={openAssign} disabled={unassigned.length === 0}>
          <Plus className="size-4" /> Assign VM
        </Button>
      }
    >
      {rollup.length > 0 && (
        <Card className="mb-4">
          <CardHeader>
            <CardTitle className="text-base">Per-hypervisor</CardTitle>
            <CardDescription>Placements committed by Waggle on each node.</CardDescription>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Hypervisor</TableHead>
                  <TableHead>VMs</TableHead>
                  <TableHead>vCPU</TableHead>
                  <TableHead>RAM (GB)</TableHead>
                  <TableHead>Disk (GB)</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rollup.map(([name, r]) => (
                  <TableRow key={name}>
                    <TableCell className="font-medium">{name}</TableCell>
                    <TableCell className="tabular-nums">{r.count}</TableCell>
                    <TableCell className="tabular-nums">{r.vcpu}</TableCell>
                    <TableCell className="tabular-nums">{r.ram}</TableCell>
                    <TableCell className="tabular-nums">{r.disk}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Pool</TableHead>
            <TableHead>Hypervisor</TableHead>
            <TableHead>Slot</TableHead>
            <TableHead>VMID</TableHead>
            <TableHead>Placed</TableHead>
            <TableHead className="w-10" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {items.map((p) => (
            <TableRow key={p.id}>
              <TableCell>
                <CopyId id={p.id} />
              </TableCell>
              <TableCell className="font-medium">{p.pool_name}</TableCell>
              <TableCell>{p.hypervisor_name}</TableCell>
              <TableCell className="text-muted-foreground">
                {p.slot_name} ({p.vcpu}c/{p.ram_gb}G/{p.disk_gb}G)
              </TableCell>
              <TableCell className="tabular-nums">
                {p.vmid != null ? (
                  <Badge variant="secondary">{p.vmid}</Badge>
                ) : (
                  <span className="text-muted-foreground">—</span>
                )}
              </TableCell>
              <TableCell className="text-muted-foreground">
                {new Date(p.created_at).toLocaleString()}
              </TableCell>
              <TableCell>
                <RowActions
                  actions={[
                    { label: "Backfill VMID", onSelect: () => openBackfill(p) },
                    { label: "Delete", variant: "destructive", onSelect: () => void remove(p) },
                  ]}
                />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <FormDialog
        open={assignOpen}
        onOpenChange={(o) => { if (!o) setAssignOpen(false) }}
        title="Assign VM to placement"
        description="Select an unassigned placement and enter the Proxmox vmid of the VM that was created for it."
        onSubmit={() => void submitAssign()}
        submitting={backfillMut.isPending}
        error={error}
        submitLabel="Assign"
      >
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <Label htmlFor="assign-placement">Placement</Label>
            <Select value={assignTarget} onValueChange={setAssignTarget}>
              <SelectTrigger id="assign-placement">
                <SelectValue placeholder="Select a placement…" />
              </SelectTrigger>
              <SelectContent>
                {unassigned.map((p) => (
                  <SelectItem key={p.id} value={p.id}>
                    {p.pool_name} — {p.hypervisor_name} ({p.id.slice(0, 8)}…)
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="assign-vmid">Proxmox VMID</Label>
            <Input
              id="assign-vmid"
              type="number"
              min={1}
              placeholder="e.g. 100"
              value={vmidInput}
              onChange={(e) => setVmidInput(e.target.value)}
              required
            />
          </div>
        </div>
      </FormDialog>

      <FormDialog
        open={backfilling !== null}
        onOpenChange={(o) => !o && setBackfilling(null)}
        title="Backfill VMID"
        description={
          backfilling
            ? `Set the Proxmox vmid for placement ${backfilling.id.slice(0, 8)}… on ${backfilling.hypervisor_name}.`
            : ""
        }
        onSubmit={() => void submitBackfill()}
        submitting={backfillMut.isPending}
        error={error}
        submitLabel="Save"
      >
        <div className="flex flex-col gap-2">
          <Label htmlFor="vmid-input">Proxmox VMID</Label>
          <Input
            id="vmid-input"
            type="number"
            min={1}
            placeholder="e.g. 100"
            value={vmidInput}
            onChange={(e) => setVmidInput(e.target.value)}
            required
          />
        </div>
      </FormDialog>
    </PageShell>
  )
}
