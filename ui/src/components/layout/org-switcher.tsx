import { useState } from "react"
import { useQuery } from "@tanstack/react-query"
import { Check, ChevronsUpDown } from "lucide-react"

import { listOrgsOptions } from "@/sdk/@tanstack/react-query.gen"
import { useAuth } from "@/lib/auth"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { cn } from "@/lib/utils"

export function OrgSwitcher() {
  const { session, switchOrg } = useAuth()
  const orgs = useQuery(listOrgsOptions())
  const [busy, setBusy] = useState<string | null>(null)

  const items = orgs.data?.items ?? []
  const currentId = session?.org?.id
  const current = items.find((o) => o.id === currentId)

  async function pick(id: string) {
    if (id === currentId) return
    setBusy(id)
    await switchOrg(id)
    setBusy(null)
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="gap-2">
          <span className="max-w-40 truncate">
            {current?.name ?? session?.org?.slug ?? "Select org"}
          </span>
          <ChevronsUpDown className="size-3.5 opacity-60" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="min-w-56">
        <DropdownMenuLabel>Organizations</DropdownMenuLabel>
        <DropdownMenuSeparator />
        {items.length === 0 ? (
          <DropdownMenuItem disabled>No organizations</DropdownMenuItem>
        ) : (
          items.map((o) => (
            <DropdownMenuItem key={o.id} onSelect={() => void pick(o.id)} disabled={busy !== null}>
              <Check className={cn("size-4", o.id === currentId ? "opacity-100" : "opacity-0")} />
              <span className="flex-1 truncate">{o.name}</span>
              <span className="text-muted-foreground text-xs">{o.role}</span>
            </DropdownMenuItem>
          ))
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
