// errMsg extracts a human-readable message from an SDK/network error. huma
// returns RFC7807 problem bodies ({ detail, title }); fall back gracefully.
export function errMsg(e: unknown): string {
  if (!e) return "Something went wrong."
  if (typeof e === "string") return e
  if (e instanceof Error && e.message) return e.message
  if (typeof e === "object") {
    const o = e as { detail?: string; title?: string; message?: string }
    return o.detail || o.title || o.message || "Request failed."
  }
  return "Request failed."
}
