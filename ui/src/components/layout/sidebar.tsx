import {
  LayoutDashboard,
  Server,
  Boxes,
  Cpu,
  Layers,
  KeyRound,
  MapPin,
  Building2,
} from "lucide-react"

import { cn } from "@/lib/utils"

export type NavKey =
  | "dashboard"
  | "datacenters"
  | "slots"
  | "hypervisors"
  | "pools"
  | "placements"
  | "organizations"
  | "api-keys"

const NAV: { key: NavKey; label: string; icon: typeof LayoutDashboard }[] = [
  { key: "dashboard", label: "Dashboard", icon: LayoutDashboard },
  { key: "datacenters", label: "Datacenters", icon: Server },
  { key: "slots", label: "Slots", icon: Layers },
  { key: "hypervisors", label: "Hypervisors", icon: Cpu },
  { key: "pools", label: "Pools", icon: Boxes },
  { key: "placements", label: "Placements", icon: MapPin },
  { key: "organizations", label: "Organizations", icon: Building2 },
  { key: "api-keys", label: "API keys", icon: KeyRound },
]

export function Sidebar({
  active,
  onNavigate,
}: {
  active: NavKey
  onNavigate: (key: NavKey) => void
}) {
  return (
    <aside className="bg-sidebar text-sidebar-foreground border-sidebar-border hidden w-60 shrink-0 flex-col border-r md:flex">
      <div className="flex h-14 items-center gap-2 px-5">
        <div className="bg-sidebar-primary text-sidebar-primary-foreground flex size-7 items-center justify-center rounded-md font-semibold">
          W
        </div>
        <span className="text-base font-semibold tracking-tight">Waggle</span>
      </div>
      <nav className="flex flex-1 flex-col gap-1 px-3 py-2">
        {NAV.map(({ key, label, icon: Icon }) => (
          <button
            key={key}
            type="button"
            onClick={() => onNavigate(key)}
            className={cn(
              "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors cursor-pointer",
              active === key
                ? "bg-sidebar-accent text-sidebar-accent-foreground"
                : "text-sidebar-foreground/70 hover:bg-sidebar-accent/60 hover:text-sidebar-accent-foreground"
            )}
          >
            <Icon className="size-4" />
            {label}
          </button>
        ))}
      </nav>
      <div className="text-sidebar-foreground/50 px-5 py-4 text-xs">
        Placement oracle &amp; ledger
      </div>
    </aside>
  )
}
