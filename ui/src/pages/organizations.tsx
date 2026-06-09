import { useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Plus, UserPlus } from "lucide-react"

import {
  listOrgsOptions,
  listOrgsQueryKey,
  createOrgMutation,
  updateOrgMutation,
  deleteOrgMutation,
  listMembersOptions,
  listMembersQueryKey,
  addMemberMutation,
  updateMemberMutation,
  removeMemberMutation,
} from "@/sdk/@tanstack/react-query.gen"
import type { OrgFullView } from "@/sdk/types.gen"
import { useAuth } from "@/lib/auth"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Dialog,
  DialogContent,
  DialogDescription,
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
import { RowActions, type RowAction } from "@/components/crud/row-actions"
import { PageShell } from "@/components/crud/page-shell"
import { errMsg } from "@/lib/errors"
import { cn } from "@/lib/utils"

const selectCls =
  "border-input h-8 rounded-md border bg-transparent px-2 text-sm outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]"

type Role = "owner" | "admin" | "member"

type EditOrg = { mode: "create" } | { mode: "edit"; org: OrgFullView } | null

export function OrganizationsPage() {
  const qc = useQueryClient()
  const { session } = useAuth()
  const list = useQuery(listOrgsOptions())
  const [editing, setEditing] = useState<EditOrg>(null)
  const [name, setName] = useState("")
  const [managing, setManaging] = useState<OrgFullView | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)

  const invalidate = () => qc.invalidateQueries({ queryKey: listOrgsQueryKey() })
  const createMut = useMutation(createOrgMutation())
  const updateMut = useMutation(updateOrgMutation())
  const deleteMut = useMutation(deleteOrgMutation())

  const orgs = list.data?.items ?? []

  async function saveOrg() {
    setError(null)
    try {
      if (editing?.mode === "edit") {
        await updateMut.mutateAsync({ body: { name }, path: { id: editing.org.id } })
      } else {
        await createMut.mutateAsync({ body: { name } })
      }
      await invalidate()
      setEditing(null)
    } catch (e) {
      setError(errMsg(e))
    }
  }

  async function remove(org: OrgFullView) {
    if (!confirm(`Delete "${org.name}"? This tears down its tenant database.`)) return
    try {
      await deleteMut.mutateAsync({ path: { id: org.id } })
      await invalidate()
      setNotice(`Deleting ${org.name} — tenant teardown enqueued.`)
    } catch (e) {
      setNotice(errMsg(e))
    }
  }

  return (
    <PageShell
      title="Organizations"
      description="Your organizations and their members. Each org is an isolated tenant."
      action={
        <Button
          onClick={() => {
            setName("")
            setError(null)
            setEditing({ mode: "create" })
          }}
        >
          <Plus className="size-4" /> New organization
        </Button>
      }
      notice={notice}
      onDismissNotice={() => setNotice(null)}
      query={list}
      empty={orgs.length === 0}
      emptyText="No organizations."
    >
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Slug</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Your role</TableHead>
            <TableHead className="w-10" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {orgs.map((o) => {
            const actions: RowAction[] = [{ label: "Members", onSelect: () => setManaging(o) }]
            if (o.role === "owner" || o.role === "admin")
              actions.push({
                label: "Rename",
                onSelect: () => {
                  setName(o.name)
                  setError(null)
                  setEditing({ mode: "edit", org: o })
                },
              })
            if (o.role === "owner")
              actions.push({ label: "Delete", variant: "destructive", onSelect: () => void remove(o) })
            return (
              <TableRow key={o.id}>
                <TableCell className="font-medium">
                  {o.name}
                  {o.id === session?.org?.id && (
                    <Badge variant="secondary" className="ml-2">current</Badge>
                  )}
                </TableCell>
                <TableCell className="text-muted-foreground">{o.slug}</TableCell>
                <TableCell>
                  <Badge variant={o.status === "active" ? "default" : "outline"}>{o.status}</Badge>
                </TableCell>
                <TableCell className="text-muted-foreground">{o.role}</TableCell>
                <TableCell>
                  <RowActions actions={actions} />
                </TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>

      <FormDialog
        open={editing !== null}
        onOpenChange={(o) => !o && setEditing(null)}
        title={editing?.mode === "edit" ? "Rename organization" : "New organization"}
        description={
          editing?.mode === "edit"
            ? undefined
            : "You'll be the owner. A tenant database is provisioned in the background."
        }
        onSubmit={() => void saveOrg()}
        submitting={createMut.isPending || updateMut.isPending}
        error={error}
        submitLabel={editing?.mode === "edit" ? "Save" : "Create"}
      >
        <div className="flex flex-col gap-2">
          <Label htmlFor="org-name">Name</Label>
          <Input id="org-name" value={name} onChange={(e) => setName(e.target.value)} required />
        </div>
      </FormDialog>

      <MembersDialog org={managing} onClose={() => setManaging(null)} selectCls={selectCls} />
    </PageShell>
  )
}

function MembersDialog({
  org,
  onClose,
  selectCls,
}: {
  org: OrgFullView | null
  onClose: () => void
  selectCls: string
}) {
  const qc = useQueryClient()
  const members = useQuery({
    ...listMembersOptions({ path: { id: org?.id ?? "" } }),
    enabled: !!org,
  })
  const addMut = useMutation(addMemberMutation())
  const roleMut = useMutation(updateMemberMutation())
  const removeMut = useMutation(removeMemberMutation())
  const [email, setEmail] = useState("")
  const [role, setRole] = useState("member")
  const [err, setErr] = useState<string | null>(null)
  const [info, setInfo] = useState<string | null>(null)

  const canManage = org?.role === "owner" || org?.role === "admin"
  const items = members.data?.items ?? []
  const invalidate = () =>
    org && qc.invalidateQueries({ queryKey: listMembersQueryKey({ path: { id: org.id } }) })

  async function add() {
    if (!org) return
    setErr(null)
    setInfo(null)
    try {
      const res = await addMut.mutateAsync({ body: { email, role: role as Role }, path: { id: org.id } })
      await invalidate()
      setEmail("")
      setInfo(res?.invited ? `Invite sent to ${email}.` : `${email} added.`)
    } catch (e) {
      setErr(errMsg(e))
    }
  }

  async function changeRole(userId: string, newRole: string) {
    if (!org) return
    setErr(null)
    try {
      await roleMut.mutateAsync({ body: { role: newRole as Role }, path: { id: org.id, userId } })
      await invalidate()
    } catch (e) {
      setErr(errMsg(e))
    }
  }

  async function kick(userId: string, label: string) {
    if (!org) return
    if (!confirm(`Remove ${label} from ${org.name}?`)) return
    setErr(null)
    try {
      await removeMut.mutateAsync({ path: { id: org.id, userId } })
      await invalidate()
    } catch (e) {
      setErr(errMsg(e))
    }
  }

  return (
    <Dialog open={org !== null} onOpenChange={(o) => !o && onClose()}>
      <DialogContent className="sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>Members — {org?.name}</DialogTitle>
          <DialogDescription>
            {canManage ? "Invite by email and manage roles." : "You can view members of this org."}
          </DialogDescription>
        </DialogHeader>

        {canManage && (
          <form
            className="flex items-end gap-2"
            onSubmit={(e) => {
              e.preventDefault()
              void add()
            }}
          >
            <div className="flex flex-1 flex-col gap-1">
              <Label htmlFor="m-email" className="text-xs">Email</Label>
              <Input id="m-email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
            </div>
            <div className="flex flex-col gap-1">
              <Label htmlFor="m-role" className="text-xs">Role</Label>
              <select id="m-role" className={cn(selectCls, "h-9")} value={role} onChange={(e) => setRole(e.target.value)}>
                <option value="member">member</option>
                <option value="admin">admin</option>
                {org?.role === "owner" && <option value="owner">owner</option>}
              </select>
            </div>
            <Button type="submit" disabled={addMut.isPending}>
              <UserPlus className="size-4" /> Add
            </Button>
          </form>
        )}
        {info && <p className="text-muted-foreground text-sm">{info}</p>}
        {err && <p className="text-destructive text-sm">{err}</p>}

        {members.isPending ? (
          <p className="text-muted-foreground py-4 text-center text-sm">Loading…</p>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Member</TableHead>
                <TableHead>Role</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="w-10" />
              </TableRow>
            </TableHeader>
            <TableBody>
              {items.map((m) => (
                <TableRow key={m.user_id}>
                  <TableCell>
                    <div className="font-medium">{m.display_name || m.email}</div>
                    <div className="text-muted-foreground text-xs">{m.email}</div>
                  </TableCell>
                  <TableCell>
                    {canManage ? (
                      <select
                        className={selectCls}
                        value={m.role}
                        onChange={(e) => void changeRole(m.user_id, e.target.value)}
                      >
                        <option value="member">member</option>
                        <option value="admin">admin</option>
                        <option value="owner">owner</option>
                      </select>
                    ) : (
                      m.role
                    )}
                  </TableCell>
                  <TableCell>
                    {m.pending ? (
                      <Badge variant="outline">invited</Badge>
                    ) : (
                      <Badge variant="secondary">active</Badge>
                    )}
                  </TableCell>
                  <TableCell>
                    {canManage && (
                      <RowActions
                        actions={[
                          {
                            label: "Remove",
                            variant: "destructive",
                            onSelect: () => void kick(m.user_id, m.display_name || m.email),
                          },
                        ]}
                      />
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </DialogContent>
    </Dialog>
  )
}
