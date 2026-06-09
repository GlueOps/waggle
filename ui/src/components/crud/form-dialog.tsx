import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"

export function FormDialog({
  open,
  onOpenChange,
  title,
  description,
  onSubmit,
  submitting,
  error,
  submitLabel = "Save",
  children,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  description?: string
  onSubmit: () => void
  submitting?: boolean
  error?: string | null
  submitLabel?: string
  children: React.ReactNode
}) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          {description && <DialogDescription>{description}</DialogDescription>}
        </DialogHeader>
        <form
          className="flex flex-col gap-4"
          onSubmit={(e) => {
            e.preventDefault()
            onSubmit()
          }}
        >
          {children}
          {error && <p className="text-destructive text-sm">{error}</p>}
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={submitting}>
              {submitting ? "Saving…" : submitLabel}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
