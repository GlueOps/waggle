import { useState } from "react"
import { Check, Copy } from "lucide-react"

import { cn } from "@/lib/utils"

/**
 * Shows a resource UUID in monospace with a click-to-copy affordance. These IDs
 * (datacenter, slot, hypervisor) are reused as inputs in Terraform, so make
 * them easy to grab.
 */
export function CopyId({ id, className }: { id: string; className?: string }) {
  const [copied, setCopied] = useState(false)

  async function copy() {
    try {
      await navigator.clipboard.writeText(id)
      setCopied(true)
      setTimeout(() => setCopied(false), 1200)
    } catch {
      // clipboard unavailable (e.g. non-secure context) — ignore
    }
  }

  return (
    <button
      type="button"
      onClick={() => void copy()}
      title={`Copy ID ${id}`}
      className={cn(
        "group inline-flex items-center gap-1.5 font-mono text-xs text-muted-foreground transition-colors hover:text-foreground",
        className,
      )}
    >
      <span>{id}</span>
      {copied ? (
        <Check className="size-3 shrink-0 text-emerald-500" />
      ) : (
        <Copy className="size-3 shrink-0 opacity-50 group-hover:opacity-100" />
      )}
    </button>
  )
}
