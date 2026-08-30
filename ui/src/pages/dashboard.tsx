import { useQuery } from "@tanstack/react-query"
import { Boxes, Cpu, HardDrive, MapPin, MemoryStick, Server } from "lucide-react"

import {
  listHypervisorsOptions,
  listPlacementsOptions,
  listPoolsOptions,
} from "@/sdk/@tanstack/react-query.gen"
import type { HypervisorView } from "@/sdk/types.gen"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

type Seg = { reserved: number; used: number; pending: number; total: number }

// StackedMeter shows three committed segments against capacity:
// reserved (OS headroom) + used (existing guests) + pending (Waggle placements).
function StackedMeter({ seg, unit }: { seg: Seg; unit: string }) {
  const committed = seg.reserved + seg.used + seg.pending
  const free = Math.max(0, seg.total - committed)
  const over = committed > seg.total
  const w = (n: number) => (seg.total > 0 ? `${Math.min(100, (n / seg.total) * 100)}%` : "0%")
  const p = seg.total > 0 ? Math.round((committed / seg.total) * 100) : 0

  return (
    <div className="min-w-40">
      <div className="text-muted-foreground mb-1 flex justify-between text-xs tabular-nums">
        <span>{committed.toLocaleString()} / {seg.total.toLocaleString()} {unit}</span>
        <span className={over ? "text-destructive font-medium" : ""}>{p}%</span>
      </div>
      <div className="bg-muted flex h-2 w-full overflow-hidden rounded-full">
        {over ? (
          <div className="h-full w-full bg-destructive" title="over-committed" />
        ) : (
          <>
            <div className="h-full bg-muted-foreground/40" style={{ width: w(seg.reserved) }} title={`reserved ${seg.reserved}`} />
            <div className="h-full bg-primary" style={{ width: w(seg.used) }} title={`used ${seg.used}`} />
            <div className="h-full bg-sky-500" style={{ width: w(seg.pending) }} title={`pending ${seg.pending}`} />
            <div className="h-full" style={{ width: w(free) }} />
          </>
        )}
      </div>
    </div>
  )
}

function LegendDot({ cls, label }: { cls: string; label: string }) {
  return (
    <span className="text-muted-foreground inline-flex items-center gap-1.5 text-xs">
      <span className={`size-2 rounded-full ${cls}`} /> {label}
    </span>
  )
}

export function DashboardPage() {
  const hvQ = useQuery(listHypervisorsOptions())
  const plQ = useQuery(listPlacementsOptions())
  const poolQ = useQuery(listPoolsOptions())

  const hvs: HypervisorView[] = hvQ.data?.items ?? []
  const placements = plQ.data?.items ?? []
  const pools = poolQ.data?.items ?? []

  // Capacity Waggle has committed per hypervisor, DERIVED from the server's own
  // bookable figure rather than recomputed here:
  //
  //   pending = total - reserved - used - bookable
  //
  // This used to sum placements grouped by hypervisor_name, which duplicated
  // consumedByHypervisor on the client. That duplication has now drifted twice.
  // Most recently the server began charging a placement to the hypervisor its
  // guest was actually OBSERVED on rather than the one it was assigned to; the
  // dashboard kept grouping by the assignment, so a node carrying a guest it
  // was never assigned looked emptier here than the API said it was, and the
  // node that lost one looked fuller.
  //
  // Deriving from bookable cannot drift: the arithmetic lives in exactly one
  // place, on the server, and the meter is a view of it. A placement's
  // hypervisor_name is the ASSIGNMENT and is deliberately not enough to compute
  // capacity from.
  const pend = (h: HypervisorView) => ({
    cpu: h.cpu_effective_total - h.cpu_reserved - h.cpu_used - h.cpu_bookable,
    ram: h.ram_gb_total - h.ram_gb_reserved - h.ram_gb_used - h.ram_gb_bookable,
    disk: h.disk_gb_total - h.disk_gb_reserved - h.disk_gb_used - h.disk_gb_bookable,
  })

  const seg = (
    pick: (h: HypervisorView) => { reserved: number; used: number; total: number; pending: number }
  ): Seg =>
    hvs.reduce(
      (a, h) => {
        const v = pick(h)
        return {
          reserved: a.reserved + v.reserved,
          used: a.used + v.used,
          pending: a.pending + v.pending,
          total: a.total + v.total,
        }
      },
      { reserved: 0, used: 0, pending: 0, total: 0 }
    )

  // CPU meters measure against the EFFECTIVE total (cores x overcommit ratio),
  // which is what the scheduler books against. Using cpu_total would render an
  // oversold node as permanently over-full and understate fleet capacity.
  const cpu = seg((h) => ({ reserved: h.cpu_reserved, used: h.cpu_used, pending: pend(h).cpu, total: h.cpu_effective_total }))
  const ram = seg((h) => ({ reserved: h.ram_gb_reserved, used: h.ram_gb_used, pending: pend(h).ram, total: h.ram_gb_total }))
  const disk = seg((h) => ({ reserved: h.disk_gb_reserved, used: h.disk_gb_used, pending: pend(h).disk, total: h.disk_gb_total }))

  // Physical cores, kept alongside the effective vCPU pool so the card can show
  // what the hardware actually is as well as what is being sold from it.
  const cores = hvs.reduce((a, h) => a + h.cpu_total, 0)
  const oversold = cores > 0 && cpu.total !== cores

  const schedulable = hvs.filter((h) => h.schedulable).length
  const stats = [
    { label: "Hypervisors", value: hvs.length, hint: `${schedulable} schedulable`, icon: Server },
    { label: "Pools", value: pools.length, hint: "declared", icon: Boxes },
    { label: "Placements", value: placements.length, hint: "VMs placed", icon: MapPin },
  ]

  return (
    <div className="flex flex-col gap-6 p-6">
      {hvQ.isError && (
        <p className="text-muted-foreground text-sm">
          Fleet data unavailable — the tenant database may still be provisioning.
        </p>
      )}

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {stats.map(({ label, value, hint, icon: Icon }) => (
          <Card key={label}>
            <CardHeader>
              <CardDescription className="flex items-center gap-2">
                <Icon className="size-4" />
                {label}
              </CardDescription>
              <CardTitle className="text-3xl tabular-nums">{value}</CardTitle>
            </CardHeader>
            <CardContent className="text-muted-foreground text-sm">{hint}</CardContent>
          </Card>
        ))}
      </div>

      <div className="grid gap-4 sm:grid-cols-3">
        <Card>
          <CardHeader><CardDescription className="flex items-center gap-2"><Cpu className="size-4" /> vCPU</CardDescription></CardHeader>
          <CardContent>
            <StackedMeter seg={cpu} unit="vCPU" />
            {/* The meter is vCPU (cores x ratio); this line keeps the physical
                hardware visible so an overcommitted fleet is not mistaken for
                having more cores than it does. */}
            <p className="text-muted-foreground mt-2 text-xs">
              {cores} physical {cores === 1 ? "core" : "cores"}
              {oversold && ` · ${(cpu.total / cores).toFixed(2)}x overcommit`}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader><CardDescription className="flex items-center gap-2"><MemoryStick className="size-4" /> Memory</CardDescription></CardHeader>
          <CardContent><StackedMeter seg={ram} unit="GB" /></CardContent>
        </Card>
        <Card>
          <CardHeader><CardDescription className="flex items-center gap-2"><HardDrive className="size-4" /> Disk</CardDescription></CardHeader>
          <CardContent><StackedMeter seg={disk} unit="GB" /></CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Hypervisors</CardTitle>
          <CardDescription className="flex flex-wrap gap-3">
            <LegendDot cls="bg-muted-foreground/40" label="reserved (OS)" />
            <LegendDot cls="bg-primary" label="used (other guests)" />
            <LegendDot cls="bg-sky-500" label="placed (waggle)" />
          </CardDescription>
        </CardHeader>
        <CardContent>
          {hvs.length === 0 ? (
            <p className="text-muted-foreground py-6 text-center text-sm">
              No hypervisors. Run discovery from a datacenter.
            </p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Node</TableHead>
                  <TableHead>CPU</TableHead>
                  <TableHead>Memory</TableHead>
                  <TableHead>Disk</TableHead>
                  <TableHead>Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {hvs.map((h) => {
                  const pp = pend(h)
                  return (
                    <TableRow key={h.id}>
                      <TableCell className="font-medium">{h.name}</TableCell>
                      <TableCell>
                        <StackedMeter seg={{ reserved: h.cpu_reserved, used: h.cpu_used, pending: pp.cpu, total: h.cpu_effective_total }} unit="" />
                        <p className="text-muted-foreground mt-1 text-xs">
                          {h.cpu_total} {h.cpu_total === 1 ? "core" : "cores"}
                          {h.cpu_overcommit_ratio !== 1 &&
                            ` × ${h.cpu_overcommit_ratio} = ${h.cpu_effective_total} vCPU`}
                        </p>
                      </TableCell>
                      <TableCell><StackedMeter seg={{ reserved: h.ram_gb_reserved, used: h.ram_gb_used, pending: pp.ram, total: h.ram_gb_total }} unit="GB" /></TableCell>
                      <TableCell><StackedMeter seg={{ reserved: h.disk_gb_reserved, used: h.disk_gb_used, pending: pp.disk, total: h.disk_gb_total }} unit="GB" /></TableCell>
                      <TableCell>
                        {h.schedulable ? (
                          <Badge variant="secondary">schedulable</Badge>
                        ) : (
                          <Badge variant="outline">drained</Badge>
                        )}
                      </TableCell>
                    </TableRow>
                  )
                })}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
