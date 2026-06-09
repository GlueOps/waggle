import { X } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { errMsg } from "@/lib/errors"

type QueryLike = { isPending: boolean; isError: boolean; error?: unknown }

export function PageShell({
  title,
  description,
  action,
  notice,
  onDismissNotice,
  query,
  empty,
  emptyText = "Nothing here yet.",
  children,
}: {
  title: string
  description?: string
  action?: React.ReactNode
  notice?: string | null
  onDismissNotice?: () => void
  query?: QueryLike
  empty?: boolean
  emptyText?: string
  children: React.ReactNode
}) {
  // children ALWAYS render so dialogs/forms nested in them stay mounted
  // regardless of list state — otherwise "Add" can't open a dialog on an empty
  // list. The loading/error/empty status shows as a note above them; during
  // those states the list itself is empty, so no duplicate data is displayed.
  const status = query?.isPending
    ? { text: "Loading…", cls: "text-muted-foreground" }
    : query?.isError
      ? { text: errMsg(query.error), cls: "text-destructive" }
      : empty
        ? { text: emptyText, cls: "text-muted-foreground" }
        : null

  return (
    <div className="flex flex-col gap-6 p-6">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h2 className="text-lg font-semibold tracking-tight">{title}</h2>
          {description && (
            <p className="text-muted-foreground text-sm">{description}</p>
          )}
        </div>
        {action}
      </div>

      {notice && (
        <div className="bg-muted text-foreground flex items-center justify-between gap-3 rounded-md border px-3 py-2 text-sm">
          <span>{notice}</span>
          {onDismissNotice && (
            <Button variant="ghost" size="icon" className="size-6" onClick={onDismissNotice}>
              <X className="size-3.5" />
            </Button>
          )}
        </div>
      )}

      <Card>
        <CardContent>
          {status && (
            <p className={`py-8 text-center text-sm ${status.cls}`}>{status.text}</p>
          )}
          {children}
        </CardContent>
      </Card>
    </div>
  )
}
