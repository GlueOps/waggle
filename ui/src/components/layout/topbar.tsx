import { useQuery } from "@tanstack/react-query"
import { LogOut } from "lucide-react"

import { healthOptions } from "@/sdk/@tanstack/react-query.gen"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { ModeToggle } from "@/components/theme/mode-toggle"
import { OrgSwitcher } from "@/components/layout/org-switcher"
import { useAuth } from "@/lib/auth"

export function Topbar({ title }: { title: string }) {
  const health = useQuery({ ...healthOptions(), refetchInterval: 15_000 })
  const { logout } = useAuth()

  const status = health.isPending
    ? { label: "Checking…", variant: "secondary" as const }
    : health.isError
      ? { label: "API offline", variant: "destructive" as const }
      : { label: "API healthy", variant: "default" as const }

  return (
    <header className="bg-background/80 sticky top-0 z-10 flex h-14 items-center justify-between border-b px-6 backdrop-blur">
      <h1 className="text-sm font-semibold tracking-tight">{title}</h1>
      <div className="flex items-center gap-3">
        <OrgSwitcher />
        <Badge variant={status.variant}>
          <span className="mr-1 inline-block size-1.5 rounded-full bg-current" aria-hidden />
          {status.label}
        </Badge>
        <ModeToggle />
        <Button variant="ghost" size="icon" onClick={() => void logout()} aria-label="Sign out">
          <LogOut className="size-4" />
        </Button>
      </div>
    </header>
  )
}
