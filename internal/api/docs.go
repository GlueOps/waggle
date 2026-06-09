package api

import (
	"fmt"
	"net/http"
)

// rapidocVersion is pinned (not "latest") for supply-chain hygiene — bump
// deliberately. See https://www.npmjs.com/package/rapidoc for releases.
const rapidocVersion = "9.3.4"

// docsCSP relaxes the strict apiCSP for the docs page so RapiDoc can load its
// bundle + fonts from the CDN and call the same-origin API ("Authorize" +
// "Try"). 'unsafe-eval' is included for broad RapiDoc compatibility; if your
// pinned version renders without it, tighten by removing it. connect-src
// 'self' is what lets the bearer-authenticated try-it requests through.
const docsCSP = "default-src 'self'; " +
	"script-src 'self' https://unpkg.com 'unsafe-inline' 'unsafe-eval'; " +
	"style-src 'self' https://unpkg.com https://fonts.googleapis.com 'unsafe-inline'; " +
	"font-src 'self' https://fonts.gstatic.com data:; " +
	"img-src 'self' data: https:; " +
	"connect-src 'self'"

// rapiDocHandler serves a RapiDoc UI page bound to the API's OpenAPI spec.
// huma's built-in docs (Stoplight) are disabled via DocsPath="" so this
// handler owns the docs path instead.
func rapiDocHandler(basePath string) http.HandlerFunc {
	page := fmt.Sprintf(`<!doctype html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Waggle API</title>
    <script type="module" src="https://unpkg.com/rapidoc@%s/dist/rapidoc-min.js"></script>
  </head>
  <body>
    <rapi-doc
      spec-url="%s/openapi.json"
      theme="dark"
      render-style="read"
      show-header="false"
      allow-spec-url-load="false"
      allow-spec-file-load="false"
      allow-authentication="true"
      persist-auth="true"
      schema-style="table">
    </rapi-doc>
  </body>
</html>`, rapidocVersion, basePath)
	body := []byte(page)

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Security-Policy", docsCSP)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(body)
	}
}
