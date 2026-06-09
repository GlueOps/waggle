import { useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Plus } from "lucide-react"

import {
  listSlotsOptions,
  listSlotsQueryKey,
  createSlotMutation,
  updateSlotMutation,
  deleteSlotMutation,
} from "@/sdk/@tanstack/react-query.gen"
import type { SlotView } from "@/sdk/types.gen"
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
import { FormDialog } from "@/components/crud/form-dialog"
import { RowActions } from "@/components/crud/row-actions"
import { PageShell } from "@/components/crud/page-shell"
import { errMsg } from "@/lib/errors"

type Editing = { mode: "create" } | { mode: "edit"; slot: SlotView } | null

export function SlotsPage() {
  const qc = useQueryClient()
  const list = useQuery(listSlotsOptions())
  const [editing, setEditing] = useState<Editing>(null)
  const [form, setForm] = useState({ name: "", vcpu: 2, ram_gb: 4, disk_gb: 40 })
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)

  const invalidate = () => qc.invalidateQueries({ queryKey: listSlotsQueryKey() })
  const createMut = useMutation(createSlotMutation())
  const updateMut = useMutation(updateSlotMutation())
  const deleteMut = useMutation(deleteSlotMutation())

  function openCreate() {
    setForm({ name: "", vcpu: 2, ram_gb: 4, disk_gb: 40 })
    setError(null)
    setEditing({ mode: "create" })
  }
  function openEdit(slot: SlotView) {
    setForm({ name: slot.name, vcpu: slot.vcpu, ram_gb: slot.ram_gb, disk_gb: slot.disk_gb })
    setError(null)
    setEditing({ mode: "edit", slot })
  }

  async function save() {
    setError(null)
    try {
      if (editing?.mode === "edit") {
        await updateMut.mutateAsync({ body: form, path: { id: editing.slot.id } })
      } else {
        await createMut.mutateAsync({ body: form })
      }
      await invalidate()
      setEditing(null)
    } catch (e) {
      setError(errMsg(e))
    }
  }

  async function remove(slot: SlotView) {
    if (!confirm(`Delete slot "${slot.name}"?`)) return
    try {
      await deleteMut.mutateAsync({ path: { id: slot.id } })
      await invalidate()
    } catch (e) {
      setNotice(errMsg(e))
    }
  }

  const items = list.data?.items ?? []
  const num = (v: string) => Math.max(0, parseInt(v || "0", 10) || 0)

  return (
    <PageShell
      title="Slots"
      description="T-shirt-size VM templates used when sizing pools."
      action={
        <Button onClick={openCreate}>
          <Plus className="size-4" /> Add slot
        </Button>
      }
      notice={notice}
      onDismissNotice={() => setNotice(null)}
      query={list}
      empty={items.length === 0}
      emptyText="No slots yet."
    >
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>vCPU</TableHead>
            <TableHead>RAM (GB)</TableHead>
            <TableHead>Disk (GB)</TableHead>
            <TableHead className="w-10" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {items.map((s) => (
            <TableRow key={s.id}>
              <TableCell className="font-medium">{s.name}</TableCell>
              <TableCell>{s.vcpu}</TableCell>
              <TableCell>{s.ram_gb}</TableCell>
              <TableCell>{s.disk_gb}</TableCell>
              <TableCell>
                <RowActions
                  actions={[
                    { label: "Edit", onSelect: () => openEdit(s) },
                    { label: "Delete", variant: "destructive", onSelect: () => void remove(s) },
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
        title={editing?.mode === "edit" ? "Edit slot" : "Add slot"}
        onSubmit={() => void save()}
        submitting={createMut.isPending || updateMut.isPending}
        error={error}
      >
        <div className="flex flex-col gap-2">
          <Label htmlFor="slot-name">Name</Label>
          <Input
            id="slot-name"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            required
          />
        </div>
        <div className="grid grid-cols-3 gap-3">
          <div className="flex flex-col gap-2">
            <Label htmlFor="slot-vcpu">vCPU</Label>
            <Input
              id="slot-vcpu"
              type="number"
              min={1}
              value={form.vcpu}
              onChange={(e) => setForm({ ...form, vcpu: num(e.target.value) })}
            />
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="slot-ram">RAM (GB)</Label>
            <Input
              id="slot-ram"
              type="number"
              min={1}
              value={form.ram_gb}
              onChange={(e) => setForm({ ...form, ram_gb: num(e.target.value) })}
            />
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="slot-disk">Disk (GB)</Label>
            <Input
              id="slot-disk"
              type="number"
              min={1}
              value={form.disk_gb}
              onChange={(e) => setForm({ ...form, disk_gb: num(e.target.value) })}
            />
          </div>
        </div>
      </FormDialog>
    </PageShell>
  )
}
