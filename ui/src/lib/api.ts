import { client } from "@/sdk/client.gen"

const TOKEN_KEY = "waggle-access-token"

let accessToken: string | null = localStorage.getItem(TOKEN_KEY)

// Wire the bearer token into every SDK request. The generated client invokes
// this for operations declared with the `bearer` security scheme and sets the
// Authorization header automatically.
client.setConfig({
  baseUrl: "/api/v1",
  auth: () => accessToken ?? undefined,
})

export function getAccessToken(): string | null {
  return accessToken
}

export function setAccessToken(token: string | null) {
  accessToken = token
  if (token) localStorage.setItem(TOKEN_KEY, token)
  else localStorage.removeItem(TOKEN_KEY)
}
