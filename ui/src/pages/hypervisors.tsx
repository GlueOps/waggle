import { useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"

import {
  listHypervisorsOptions,
  listHypervisorsQueryKey,
  updateHypervisorMutation,
  deleteHypervisorMutation,
} from "@/sdk/@tanstack/react-query.gen"
import type { HypervisorView } from "@/sdk/types.gen"
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
import { CopyId } from "@/components/ui/copy-id"
import { FormDialog } from "@/components/crud/form-dialog"
import { RowActions } from "@/components/crud/row-actions"
import { PageShell } from "@/components/crud/page-shell"
import { errMsg } from "@/lib/errors"

function bodyFor(hv: HypervisorView, overrides: Partial<HypervisorView>) {
  const m = { ...hv, ...overrides }
  return {
    datacenter_id: m.datacenter_id,
    name: m.name,
    cpu_total: m.cpu_total,
    cpu_reserved: m.cpu_reserved,
    ram_gb_total: m.ram_gb_total,
    ram_gb_reserved: m.ram_gb_reserved,
    disk_gb_total: m.disk_gb_total,
    disk_gb_reserved: m.disk_gb_reserved,
    schedulable: m.schedulable,
  }
}

export function HypervisorsPage() {
  const qc = useQueryClient()
  const list = useQuery(listHypervisorsOptions())
  const [editing, setEditing] = useState<HypervisorView | null>(null)
  const [reserved, setReserved] = useState({ cpu: 0, ram: 0, disk: 0 })
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)

  const invalidate = () => qc.invalidateQueries({ queryKey: listHypervisorsQueryKey() })
  const updateMut = useMutation(updateHypervisorMutation())
  const deleteMut = useMutation(deleteHypervisorMutation())

  async function toggleSchedulable(hv: HypervisorView, schedulable: boolean) {
    setNotice(null)
    try {
      await updateMut.mutateAsync({ body: bodyFor(hv, { schedulable }), path: { id: hv.id } })
      await invalidate()
    } catch (e) {
      setNotice(errMsg(e))
    }
  }

  function openEdit(hv: HypervisorView) {
    setReserved({ cpu: hv.cpu_reserved, ram: hv.ram_gb_reserved, disk: hv.disk_gb_reserved })
    setError(null)
    setEditing(hv)
  }

  async function saveReserved() {
    if (!editing) return
    setError(null)
    try {
      await updateMut.mutateAsync({
        body: bodyFor(editing, {
          cpu_reserved: reserved.cpu,
          ram_gb_reserved: reserved.ram,
          disk_gb_reserved: reserved.disk,
        }),
        path: { id: editing.id },
      })
      await invalidate()
      setEditing(null)
    } catch (e) {
      setError(errMsg(e))
    }
  }

  async function remove(hv: HypervisorView) {
    if (!confirm(`Delete hypervisor "${hv.name}"?`)) return
    try {
      await deleteMut.mutateAsync({ path: { id: hv.id } })
      await invalidate()
    } catch (e) {
      setNotice(errMsg(e))
    }
  }

  const items = list.data?.items ?? []
  const num = (v: string) => Math.max(0, parseInt(v || "0", 10) || 0)

  return (
    <PageShell
      title="Hypervisors"
      description="Discovered Proxmox nodes. Toggle Schedulable to drain a node without deleting it."
      notice={notice}
      onDismissNotice={() => setNotice(null)}
      query={list}
      empty={items.length === 0}
      emptyText="No hypervisors yet. Run discovery from a datacenter."
    >
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>CPU (free/total)</TableHead>
            <TableHead>RAM GB (free/total)</TableHead>
            <TableHead>Disk GB (free/total)</TableHead>
            <TableHead>Schedulable</TableHead>
            <TableHead>ID</TableHead>
            <TableHead className="w-10" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {items.map((hv) => (
            <TableRow key={hv.id}>
              <TableCell className="font-medium">{hv.name}</TableCell>
              <TableCell>{hv.cpu_bookable}/{hv.cpu_total}</TableCell>
              <TableCell>{hv.ram_gb_bookable}/{hv.ram_gb_total}</TableCell>
              <TableCell>{hv.disk_gb_bookable}/{hv.disk_gb_total}</TableCell>
              <TableCell>
                <Switch
                  checked={hv.schedulable}
                  onCheckedChange={(c) => void toggleSchedulable(hv, c)}
                />
              </TableCell>
              <TableCell><CopyId id={hv.id} /></TableCell>
              <TableCell>
                <RowActions
                  actions={[
                    { label: "Edit reserved", onSelect: () => openEdit(hv) },
                    { label: "Delete", variant: "destructive", onSelect: () => void remove(hv) },
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
        title={`Reserved capacity — ${editing?.name ?? ""}`}
        description="Reserved capacity is held back from placement (e.g. for host overhead)."
        onSubmit={() => void saveReserved()}
        submitting={updateMut.isPending}
        error={error}
      >
        <div className="grid grid-cols-3 gap-3">
          <div className="flex flex-col gap-2">
            <Label htmlFor="hv-cpu">CPU</Label>
            <Input
              id="hv-cpu"
              type="number"
              min={0}
              value={reserved.cpu}
              onChange={(e) => setReserved({ ...reserved, cpu: num(e.target.value) })}
            />
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="hv-ram">RAM (GB)</Label>
            <Input
              id="hv-ram"
              type="number"
              min={0}
              value={reserved.ram}
              onChange={(e) => setReserved({ ...reserved, ram: num(e.target.value) })}
            />
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="hv-disk">Disk (GB)</Label>
            <Input
              id="hv-disk"
              type="number"
              min={0}
              value={reserved.disk}
              onChange={(e) => setReserved({ ...reserved, disk: num(e.target.value) })}
            />
          </div>
        </div>
      </FormDialog>
    </PageShell>
  )
}
