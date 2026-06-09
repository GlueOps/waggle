import { useQuery } from "@tanstack/react-query"

import { listPlacementsOptions } from "@/sdk/@tanstack/react-query.gen"
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

export function PlacementsPage() {
  const list = useQuery(listPlacementsOptions())
  const items = list.data?.items ?? []

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
      empty={items.length === 0}
      emptyText="No placements yet. Create a pool to place VMs."
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
            <TableHead>Pool</TableHead>
            <TableHead>Hypervisor</TableHead>
            <TableHead>Slot</TableHead>
            <TableHead>VMID</TableHead>
            <TableHead>Placed</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {items.map((p) => (
            <TableRow key={p.id}>
              <TableCell className="font-medium">{p.pool_name}</TableCell>
              <TableCell>{p.hypervisor_name}</TableCell>
              <TableCell className="text-muted-foreground">
                {p.slot_name} ({p.vcpu}c/{p.ram_gb}G/{p.disk_gb}G)
              </TableCell>
              <TableCell className="tabular-nums">{p.vmid ?? "—"}</TableCell>
              <TableCell className="text-muted-foreground">
                {new Date(p.created_at).toLocaleString()}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </PageShell>
  )
}
