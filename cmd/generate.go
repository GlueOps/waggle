package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/glueops/waggle/internal/api"
	"github.com/glueops/waggle/internal/app"
	"github.com/glueops/waggle/internal/buildinfo"
	"github.com/glueops/waggle/internal/config"
	"github.com/glueops/waggle/internal/service"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Code generation tools (migrations, SDKs)",
}

// poolPlacementsDataSourceSrc is a hand-authored Terraform data source injected
// into the generated provider after each OAG run. The OpenAPI Generator cannot
// model the placements list (a list of nested objects), so this fills the gap.
//
//go:embed overlays/pool_placements_data_source.go.tmpl
var poolPlacementsDataSourceSrc string

// slotsDataSourceSrc is a hand-authored replacement for the generated slots
// data source. The OpenAPI Generator emits a by-id-only lookup; this overlay
// supports looking a slot up by id OR by its (tenant-unique) name via the
// /slots?name= filter. It overwrites the generated slots_data_source.go after
// each OAG run.
//
//go:embed overlays/slots_data_source.go.tmpl
var slotsDataSourceSrc string

// placementResourceSrc is a hand-authored replacement for the generated
// placements_resource.go. The generator produces a stub with no schema and a
// broken Create; this overlay provides a full CRUD implementation that adopts
// an existing waggle placement into Terraform state and backfills its vmid.
//
//go:embed overlays/placement_resource.go.tmpl
var placementResourceSrc string

// clientBodyTestSrc holds regression tests for the generated provider client:
// that create bodies omit read-only fields, and that the API key is sent as a
// Bearer token. Both assert the effect of a generator patch, so they live with
// the generator — internal/client in the provider repo is replaced wholesale on
// every regenerate and a test committed there would not survive.
//
//go:embed overlays/client_body_test.go.tmpl
var clientBodyTestSrc string

var generateMigrationsCmd = &cobra.Command{
	Use:   "migrations [name]",
	Short: "Generate Goose migrations from GORM models",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		log.Printf("Generating Control DB migration: %s", name)

		atlasControl := exec.Command("atlas", "migrate", "diff", name, "--env", "control")
		atlasControl.Stdout = os.Stdout
		atlasControl.Stderr = os.Stderr
		if err := atlasControl.Run(); err != nil {
			return fmt.Errorf("failed to generate control migrations: %w", err)
		}

		log.Printf("Generating Tenant DB migration: %s", name)

		atlasTenant := exec.Command("atlas", "migrate", "diff", name, "--env", "tenant")
		atlasTenant.Stdout = os.Stdout
		atlasTenant.Stderr = os.Stderr
		if err := atlasTenant.Run(); err != nil {
			return fmt.Errorf("failed to generate tenant migrations: %w", err)
		}

		log.Println("Migrations generated successfully.")
		return nil
	},
}

var generateSdkCmd = &cobra.Command{
	Use:   "sdk",
	Short: "Generate OpenAPI spec and TypeScript/Go SDKs",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("0. Removing existing generated SDKs...")
		for _, dir := range []string{
			"ui/src/sdk",
			"sdk/ts",
			"sdk/go",
		} {
			if err := os.RemoveAll(dir); err != nil {
				return fmt.Errorf("failed to remove %s: %w", dir, err)
			}
		}

		log.Println("1. Extracting OpenAPI Spec from Huma...")

		// The SDK generators need the version as an argument (npmVersion,
		// packageVersion). The spec itself gets it from buildOpenAPISpecBytes,
		// which stamps buildinfo.Version for every caller.
		specVersion := releaseVersion()
		log.Printf("spec/SDK version: %s", specVersion)

		b, err := buildOpenAPISpecBytes()
		if err != nil {
			return err
		}

		if err := os.MkdirAll("docs", 0o755); err != nil {
			return err
		}
		if err := os.WriteFile("docs/openapi.json", b, 0o644); err != nil {
			return fmt.Errorf("failed to write docs/openapi.json: %w", err)
		}

		log.Println("2. Generating embedded TypeScript SDK for UI with Hey API (ui/src/sdk)...")

		if err := generateTSUIHey(); err != nil {
			return err
		}

		log.Println("3. Generating standalone TypeScript SDK with OpenAPI Generator (sdk/ts)...")
		if err := generateTSPackagedOAG("docs/openapi.json", "sdk/ts", []string{
			"--additional-properties=supportsES6=true,npmName=@glueops/waggle-sdk,npmVersion=" + specVersion,
		}); err != nil {
			return err
		}
		if err := patchTSBuildConfig("sdk/ts"); err != nil {
			return err
		}

		log.Println("4. Generating standalone Go SDK (sdk/go)...")
		if err := generateGo("docs/openapi.json", "sdk/go", specVersion); err != nil {
			return err
		}

		log.Println("SDKs generated successfully.")
		return nil
	},
}

func generateGo(inputSpec, outDir, version string) error {
	_ = os.RemoveAll(outDir)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	gitUserID := "glueops"
	gitRepoID := filepath.ToSlash(filepath.Join("waggle", "sdk", "go"))

	args := []string{
		"--yes",
		"@openapitools/openapi-generator-cli",
		"generate",
		"-i", inputSpec,
		"-g", "go",
		"-o", outDir,
		// Go generator options like packageName/withGoMod/isGoSubmodule are supported.
		"--additional-properties=packageName=waggle,packageVersion=" + version + ",withGoMod=true,isGoSubmodule=true",
		"--git-user-id", gitUserID,
		"--git-repo-id", gitRepoID,
	}
	if err := run(exec.Command("npx", args...), "go -> "+outDir); err != nil {
		return err
	}

	// Post-process client.go to replace unsafe full request/response dumps with
	// metadata-only logging, so sensitive data (passwords, API keys) never appears
	// in debug logs regardless of which generator version produced the file.
	if err := patchGoClientLogging(outDir); err != nil {
		return err
	}

	// Force the generated module to require a newer Go version.
	edit := exec.Command("go", "mod", "edit", "-go=1.26.1")
	edit.Dir = outDir
	if err := run(edit, "set go version -> "+outDir); err != nil {
		return err
	}

	// Optional: normalize deps / let Go update related metadata.
	tidy := exec.Command("go", "mod", "tidy", "-go=1.26.1")
	tidy.Dir = outDir
	if err := run(tidy, "go mod tidy -> "+outDir); err != nil {
		return err
	}

	// Format the generated source — openapi-generator's output isn't
	// gofmt-clean out of the box, and CI's gofmt check (rightly) holds
	// the line at "every Go file in the repo, including generated, is
	// formatted." Run after tidy so any new files tidy added (none today,
	// but defensive) are formatted too.
	fmtCmd := exec.Command("gofmt", "-w", ".")
	fmtCmd.Dir = outDir
	if err := run(fmtCmd, "gofmt -w -> "+outDir); err != nil {
		return err
	}

	return nil
}

func generateTSUIHey() error {
	outDir := "ui/src/sdk"

	_ = os.RemoveAll(outDir)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	// The generator version lives in package.json / yarn.lock, NOT in an inline
	// `npx --yes pkg@x.y.z` pin here. Two reasons:
	//
	//  1. Renovate only reads manifest files, so a version pinned in Go source
	//     is invisible to it. This pin said 0.93.0 while package.json declared
	//     0.99.0 -- the drift is not hypothetical.
	//  2. `npx --yes` resolves the tool in its own sandbox, ignoring the repo's
	//     TypeScript pin. openapi-ts declares typescript as a peer with an open
	//     upper bound, so npx installs TS 7, which it does not yet support --
	//     it dies with "Cannot read properties of undefined". See
	//     https://github.com/hey-api/hey-api/issues/4235. Installing from the
	//     lockfile first is what holds TypeScript at 6.x.
	//
	// --immutable so the lockfile is authoritative and CI cannot silently float
	// to a different generator than a developer machine.
	if err := run(exec.Command("yarn", "install", "--immutable"), "yarn install (js toolchain)"); err != nil {
		return err
	}

	// Assumes openapi-ts.config.ts or openapi-ts.config.mjs exists at repo root.
	cmd := exec.Command(filepath.Join("node_modules", ".bin", "openapi-ts"))
	if err := run(cmd, "hey-api ui sdk -> "+outDir); err != nil {
		return err
	}

	return nil
}

func generateTSPackagedOAG(inputSpec, outDir string, extraArgs []string) error {
	_ = os.RemoveAll(outDir)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	gitUserID := "glueops"
	gitRepoID := "waggle"

	args := []string{
		"--yes",
		"@openapitools/openapi-generator-cli",
		"generate",
		"-i", inputSpec,
		"-g", "typescript-fetch",
		"-o", outDir,
		"--git-user-id", gitUserID,
		"--git-repo-id", gitRepoID,
	}
	args = append(args, extraArgs...)

	return run(exec.Command("npx", args...), "typescript-fetch -> "+outDir)
}

// patchTSBuildConfig fixes up the build configuration openapi-generator emits
// for the typescript-fetch client so it compiles under TypeScript 6.
//
// generateTSPackagedOAG wipes outDir on every run, so .openapi-generator-ignore
// cannot preserve hand edits here — these fixes have to be re-applied after
// each generation, the same way patchGoClientLogging handles the Go client.
//
// Two changes, both to tsconfig.json:
//   - ignoreDeprecations: TS 6 makes `moduleResolution: node` (node10) a hard
//     error without this opt-out (TS5107). Keeping node10 leaves the emitted
//     JavaScript byte-identical for consumers; a real migration to node16 /
//     bundler resolution will be required before TypeScript 7.
//   - rootDir: TS 6 requires the common source directory to be explicit
//     (TS5011). "src" is what TS previously inferred, so output layout is
//     unchanged.
//
// The generated devDependency range ("^4.0 || ^5.0") is also widened to admit
// TypeScript 6.
func patchTSBuildConfig(outDir string) error {
	tsconfigPath := filepath.Join(outDir, "tsconfig.json")
	b, err := os.ReadFile(tsconfigPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", tsconfigPath, err)
	}

	var tsconfig map[string]any
	if err := json.Unmarshal(b, &tsconfig); err != nil {
		return fmt.Errorf("parse %s: %w", tsconfigPath, err)
	}
	opts, ok := tsconfig["compilerOptions"].(map[string]any)
	if !ok {
		return fmt.Errorf("%s: compilerOptions is missing or not an object", tsconfigPath)
	}
	opts["ignoreDeprecations"] = "6.0"
	opts["rootDir"] = "src"

	out, err := json.MarshalIndent(tsconfig, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", tsconfigPath, err)
	}
	if err := os.WriteFile(tsconfigPath, append(out, '\n'), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", tsconfigPath, err)
	}

	pkgPath := filepath.Join(outDir, "package.json")
	pb, err := os.ReadFile(pkgPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", pkgPath, err)
	}
	var pkg map[string]any
	if err := json.Unmarshal(pb, &pkg); err != nil {
		return fmt.Errorf("parse %s: %w", pkgPath, err)
	}
	if dev, ok := pkg["devDependencies"].(map[string]any); ok {
		dev["typescript"] = "^6.0.0"
	}
	pout, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", pkgPath, err)
	}
	if err := os.WriteFile(pkgPath, append(pout, '\n'), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", pkgPath, err)
	}

	return nil
}

func patchGoClientLogging(outDir string) error {
	clientFile := filepath.Join(outDir, "client.go")
	b, err := os.ReadFile(clientFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", clientFile, err)
	}

	content := string(b)

	// Unsafe patterns produced by the openapi-generator Go template.
	// The format strings contain literal \n (backslash + n) as they appear in source.
	//
	// Unescaped request dump block:
	//   if c.cfg.Debug {
	//       dump, err := httputil.DumpRequestOut(request, true)
	//       if err != nil { return nil, err }
	//       log.Printf("\n%s\n", string(dump))
	//   }
	unsafeRequestDump := "\tif c.cfg.Debug {\n\t\tdump, err := httputil.DumpRequestOut(request, true)\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\tlog.Printf(\"\\n%s\\n\", string(dump))\n\t}"

	// Unescaped response dump block:
	//   if c.cfg.Debug {
	//       dump, err := httputil.DumpResponse(resp, true)
	//       if err != nil { resp.Body.Close(); return nil, err }
	//       log.Printf("\n%s\n", string(dump))
	//   }
	unsafeResponseDump := "\tif c.cfg.Debug {\n\t\tdump, err := httputil.DumpResponse(resp, true)\n\t\tif err != nil {\n\t\t\tresp.Body.Close()\n\t\t\treturn nil, err\n\t\t}\n\t\tlog.Printf(\"\\n%s\\n\", string(dump))\n\t}"

	// Safe replacements that log only metadata (method, URL, status) without bodies
	// or authentication headers.
	safeRequestLog := "\tif c.cfg.Debug {\n\t\t// Log request metadata only (method, URL) to avoid logging sensitive data in headers or body.\n\t\tlog.Printf(\"HTTP request: %s %s\", request.Method, request.URL.String())\n\t}"
	safeResponseLog := "\tif c.cfg.Debug {\n\t\t// Log response metadata only (status) to avoid logging sensitive data in headers or body.\n\t\tlog.Printf(\"HTTP response: %s\", resp.Status)\n\t}"

	patched := strings.ReplaceAll(content, unsafeRequestDump, safeRequestLog)
	patched = strings.ReplaceAll(patched, unsafeResponseDump, safeResponseLog)

	// Remove the httputil import since it is no longer used after patching.
	//patched = strings.ReplaceAll(patched, "\t\"net/http/httputil\"\n", "")

	if patched == content {
		// Distinguish between the idempotent case (already safe) and a template change.
		if strings.Contains(content, "request.Method, request.URL.String()") {
			log.Printf("note: %s already has metadata-only logging; skipping patch", clientFile)
		} else {
			log.Printf("WARNING: %s does not contain the expected unsafe logging patterns; the generator template may have changed - review debug logging in callAPI for sensitive data exposure", clientFile)
		}
		return nil
	}

	if err := os.WriteFile(clientFile, []byte(patched), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", clientFile, err)
	}
	log.Printf("patched %s: replaced unsafe request/response dumps with metadata-only logging", clientFile)
	return nil
}

func run(c *exec.Cmd, label string) error {
	log.Printf("running: %s: %s %s", label, c.Path, strings.Join(c.Args[1:], " "))
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", label, err)
	}
	return nil
}

var generateTerraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Generate Terraform provider artifacts",
}

var generateTerraformHashiCmd = &cobra.Command{
	Use:   "hashicorp",
	Short: "Generate Terraform provider code via HashiCorp codegen",
	RunE: func(cmd *cobra.Command, args []string) error {
		specPath := "docs/openapi.json"
		cfgPath := "terraform/generator_config.yml"
		specOut := "terraform/provider_code_spec.json"
		outDir := "terraform-provider-eyrie/internal/provider"

		// Ensure OpenAPI spec exists/fresh.
		if err := generateOpenAPISpec(specPath); err != nil {
			return err
		}

		// Clean generated artifacts to avoid dead code.
		for _, dir := range []string{
			specOut,
			outDir,
		} {
			if err := os.RemoveAll(dir); err != nil {
				return fmt.Errorf("failed to remove %s: %w", dir, err)
			}
		}
		if err := os.MkdirAll(filepath.Dir(specOut), 0o755); err != nil {
			return err
		}
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return err
		}

		openapiCmd := exec.Command(
			"tfplugingen-openapi",
			"generate",
			"--config", cfgPath,
			"--output", specOut,
			specPath,
		)
		if err := run(openapiCmd, "terraform hashicorp openapi -> "+specOut); err != nil {
			return err
		}

		frameworkCmd := exec.Command(
			"tfplugingen-framework",
			"generate",
			"all",
			"--input", specOut,
			"--output", outDir,
		)
		if err := run(frameworkCmd, "terraform hashicorp framework -> "+outDir); err != nil {
			return err
		}

		log.Println("HashiCorp Terraform provider artifacts generated successfully.")
		return nil
	},
}

var generateTerraformOAGCmd = &cobra.Command{
	Use:   "openapi-generator",
	Short: "Generate Terraform provider via OpenAPI Generator",
	RunE: func(cmd *cobra.Command, args []string) error {
		specPath := "docs/openapi.json"

		// The provider is its own repository (github.com/GlueOps/terraform-provider-waggle),
		// checked out alongside this one. Everything below is written into a
		// staging directory and only then synced into that checkout — see
		// syncProviderRepo for why this is not generated in place.
		target, err := filepath.Abs(providerOutDir)
		if err != nil {
			return fmt.Errorf("resolve provider output dir: %w", err)
		}
		if err := assertProviderCheckout(target); err != nil {
			return err
		}

		if err := generateOpenAPISpec(specPath); err != nil {
			return err
		}

		outDir, err := os.MkdirTemp("", "waggle-tf-provider-")
		if err != nil {
			return fmt.Errorf("create staging dir: %w", err)
		}
		defer os.RemoveAll(outDir)

		// The generator resolves -i relative to the working directory, and the
		// staging dir is elsewhere; hand it an absolute path.
		specAbs, err := filepath.Abs(specPath)
		if err != nil {
			return err
		}
		specPath = specAbs

		args = []string{
			"--yes",
			"@openapitools/openapi-generator-cli",
			"generate",
			"-i", specPath,
			"-g", "terraform-provider",
			"-o", outDir,
			"--additional-properties=providerName=waggle,packageName=waggle,providerAddress=registry.terraform.io/glueops/waggle",
			"--git-user-id", "glueops",
			"--git-repo-id", "terraform-provider-waggle",
		}
		if err := run(exec.Command("npx", args...), "terraform-provider -> "+outDir); err != nil {
			return err
		}

		// Post-process: Huma adds a "$schema" metadata property to response bodies.
		// The OpenAPI Generator renders it as a model field with a $-prefixed Go
		// identifier ($Schema) tagged tfsdk:"__schema" — which compiles but the
		// plugin framework rejects at apply time ("invalid tfsdk tag, must start
		// with a letter").
		//
		// De-dollar FIRST so $Schema -> Schema and the field reaches the canonical
		// form the drop pass below matches; running drop first misses it (the
		// generator emits the still-$-prefixed identifier), leaving the rejected
		// tfsdk:"__schema" tag behind. This also normalizes any other $-prefixed
		// identifiers the generator may emit.
		if err := patchTerraformDollarIdentifiers(outDir); err != nil {
			return err
		}
		// Now drop the $schema metadata field entirely from the provider models,
		// schema attributes, and mappings — it is metadata, not a real attribute.
		if err := dropTerraformSchemaField(outDir); err != nil {
			return err
		}

		// Auth + request-body fixes for the strict Huma API: send the org API key
		// as a Bearer token (the API's only scheme), and omitempty the read-only,
		// server-assigned fields so create/update bodies don't carry properties
		// the API rejects with 422 ("unexpected property").
		if err := patchTerraformAPIKeyBearer(outDir); err != nil {
			return err
		}
		if err := patchClientReadOnlyOmitempty(outDir); err != nil {
			return err
		}

		// Inject hand-authored provider overlays before patching schema roles so
		// patchResourceSchemaRoles can validate the overlay's attribute roles
		// idempotently (overlay already has correct roles; patch is a no-op).
		if err := writeProviderOverlays(outDir); err != nil {
			return err
		}

		// Fix the generated pool Update to use state.Id instead of plan.Id for
		// the PATCH URL. plan.Id is "(known after apply)" during an update, which
		// produces an empty path (/pools/) and an HTML error response.
		if err := patchPoolsResourceUpdate(outDir); err != nil {
			return err
		}

		// Correct resource attribute roles: the generator marks every field
		// Required, forcing server-assigned fields (id, created_at, ...) into
		// config and causing perpetual diffs. Reclassify per the spec.
		if err := patchResourceSchemaRoles(outDir); err != nil {
			return err
		}

		// Add 404 -> RemoveResource handling to all generated resource Read
		// methods. Without this, a resource deleted outside Terraform causes
		// a hard error on the next plan/apply instead of prompting recreation.
		if err := patchResourceRead404(outDir); err != nil {
			return err
		}

		// Force replacement on attributes the API refuses to change after
		// create. Runs after patchResourceSchemaRoles because it anchors on the
		// attribute's opening brace and would otherwise break that pass's regex.
		if err := patchResourceImmutability(outDir, specPath); err != nil {
			return err
		}

		// Map pools.metadata, which the generator emits as a schema attribute
		// but wires up in neither direction (free-form JSON has no Go type in
		// the spec, so the model generator skips it).
		if err := patchPoolsMetadataMapping(outDir); err != nil {
			return err
		}

		// gofmt the provider package: the schema-role and registration patches
		// emit loosely-spaced Go; gofmt normalizes alignment/indentation.
		fmtProvider := exec.Command("gofmt", "-w", ".")
		fmtProvider.Dir = filepath.Join(outDir, "internal", "provider")
		if err := run(fmtProvider, "gofmt provider -> "+outDir); err != nil {
			return err
		}

		// Copy the staged tree into the provider checkout, replacing only what
		// the generator owns.
		if err := syncProviderRepo(outDir, target); err != nil {
			return err
		}

		// internal/ was replaced wholesale and go.mod may have just gained a
		// dependency (jsontypes), so the checked-in go.sum is stale until this
		// runs. Mirrors what generateGo does for the Go SDK.
		tidyProvider := exec.Command("go", "mod", "tidy")
		tidyProvider.Dir = target
		if err := run(tidyProvider, "go mod tidy -> "+target); err != nil {
			return err
		}

		log.Printf("OpenAPI Generator Terraform provider generated successfully into %s.", target)
		return nil
	},
}

// providerModulePath is the Go module path of the Terraform provider
// repository. assertProviderCheckout uses it to confirm the output directory is
// really that checkout before anything is written into it.
const providerModulePath = "github.com/glueops/terraform-provider-waggle"

// providerOutDir is the checkout the generated Terraform provider is written
// into. The provider is a separate repository — it is published to the
// Terraform and OpenTofu registries from its own tags — so the default assumes
// it is cloned next to this one. Override with --out or WAGGLE_TF_PROVIDER_DIR.
var providerOutDir string

// defaultProviderOutDir resolves the default provider checkout location,
// honouring WAGGLE_TF_PROVIDER_DIR so CI can point at a checkout elsewhere.
func defaultProviderOutDir() string {
	if v := os.Getenv("WAGGLE_TF_PROVIDER_DIR"); v != "" {
		return v
	}
	return filepath.Join("..", "terraform-provider-waggle")
}

// assertProviderCheckout refuses to write into a directory that is neither
// empty nor the provider repository. syncProviderRepo deletes whole subtrees,
// so a mistyped --out would otherwise destroy unrelated work.
func assertProviderCheckout(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // created on first sync
		}
		return fmt.Errorf("inspect provider output dir %s: %w", dir, err)
	}
	if len(entries) == 0 {
		return nil
	}

	b, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return fmt.Errorf(
			"provider output dir %s is not empty and has no go.mod; refusing to write into it "+
				"(expected a checkout of %s — pass --out or set WAGGLE_TF_PROVIDER_DIR)",
			dir, providerModulePath,
		)
	}
	for line := range strings.SplitSeq(string(b), "\n") {
		if !strings.HasPrefix(line, "module ") {
			continue
		}
		if got := strings.TrimSpace(strings.TrimPrefix(line, "module ")); got != providerModulePath {
			return fmt.Errorf(
				"provider output dir %s declares module %q, want %q; refusing to write into it",
				dir, got, providerModulePath,
			)
		}
		return nil
	}
	return fmt.Errorf("provider output dir %s has a go.mod with no module directive; refusing to write into it", dir)
}

// generatorOwnedTrees are replaced wholesale on every run: everything under
// them is generated, so a file that disappears from the spec must disappear
// from the checkout too.
var generatorOwnedTrees = []string{
	filepath.Join("internal", "client"),
	filepath.Join("internal", "provider"),
	".openapi-generator",
}

// generatorScaffolding is written only when the provider checkout does not
// already have it. These are the files the OpenAPI Generator emits once to make
// a compilable module; the provider repository then owns them and has diverged
// on purpose — go.mod carries a newer toolchain and terraform-plugin-docs,
// README and GNUmakefile are hand-written, and examples/ is curated.
var generatorScaffolding = []string{
	".gitignore",
	".openapi-generator-ignore",
	"GNUmakefile",
	"README.md",
	"go.mod",
	"main.go",
	filepath.Join("examples", "provider", "provider.tf"),
}

// syncProviderRepo copies the staged provider tree into the real checkout.
//
// It replaces only generatorOwnedTrees and fills in generatorScaffolding when
// absent, rather than clearing the target first. The target is a git repository
// holding release machinery with no source in this repo — .git, .github/,
// .changes/, docs/, .goreleaser.yml — and a wholesale wipe would take all of it.
func syncProviderRepo(stage, target string) error {
	// Tripwire: the sync list is hand-maintained, so a new generated package
	// would otherwise be silently dropped on the floor.
	staged, err := os.ReadDir(filepath.Join(stage, "internal"))
	if err != nil {
		return fmt.Errorf("read staged internal/: %w", err)
	}
	for _, e := range staged {
		name := filepath.Join("internal", e.Name())
		if !slices.Contains(generatorOwnedTrees, name) {
			return fmt.Errorf(
				"generator emitted %s, which syncProviderRepo does not copy; add it to generatorOwnedTrees",
				name,
			)
		}
	}

	if err := os.MkdirAll(target, 0o755); err != nil {
		return err
	}

	for _, tree := range generatorOwnedTrees {
		src := filepath.Join(stage, tree)
		if _, err := os.Stat(src); err != nil {
			return fmt.Errorf("generator produced no %s: %w", tree, err)
		}
		dst := filepath.Join(target, tree)
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("clear %s: %w", dst, err)
		}
		if err := copyTreeOverwrite(dst, src); err != nil {
			return fmt.Errorf("copy %s: %w", tree, err)
		}
		log.Printf("synced %s -> %s", tree, dst)
	}

	for _, name := range generatorScaffolding {
		dst := filepath.Join(target, name)
		if _, err := os.Stat(dst); err == nil {
			continue // the provider repo owns it
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("stat %s: %w", dst, err)
		}
		b, err := os.ReadFile(filepath.Join(stage, name))
		if err != nil {
			if os.IsNotExist(err) {
				continue // generator did not emit it this run
			}
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, b, 0o644); err != nil {
			return err
		}
		log.Printf("wrote missing scaffolding %s", dst)
	}
	return nil
}

// schemaRole classifies a generated attribute. The OpenAPI Generator emits
// every attribute as Required; these roles correct that.
type schemaRole int

const (
	roleRequired schemaRole = iota
	roleOptional
	roleOptionalComputed
	roleComputed
)

// flags renders the terraform-plugin-framework attribute flag lines for a role.
// Spacing is loose on purpose — a gofmt pass after patching aligns it.
func (r schemaRole) flags() string {
	switch r {
	case roleRequired:
		return "Required: true,"
	case roleOptional:
		return "Optional: true,"
	case roleOptionalComputed:
		return "Optional: true,\nComputed: true,"
	default:
		return "Computed: true,"
	}
}

// resourceSchemaRoles classifies each resource view field: real inputs
// (Required/Optional), settable-with-server-default (OptionalComputed), or
// server-assigned read-only (Computed). Derived from each entity's create/
// update request body in the OpenAPI spec. Fields not listed (e.g. __schema,
// already Computed) are left as generated.
var resourceSchemaRoles = map[string]map[string]schemaRole{
	"organizations": {
		"name":       roleRequired,
		"id":         roleComputed,
		"created_at": roleComputed,
		"domain":     roleComputed,
		"role":       roleComputed,
		"slug":       roleComputed,
		"status":     roleComputed,
	},
	"datacenters": {
		"name":                 roleRequired,
		"url":                  roleRequired,
		"insecure_skip_verify": roleOptionalComputed,
		// Optional+computed: omitting it keeps the server's 1.0 default rather
		// than showing a perpetual diff against an unset attribute.
		"cpu_overcommit_ratio": roleOptionalComputed,
		"id":                   roleComputed,
		"created_at":           roleComputed,
		"updated_at":           roleComputed,
		"has_token":            roleComputed,
	},
	"slots": {
		"name":       roleRequired,
		"vcpu":       roleRequired,
		"ram_gb":     roleRequired,
		"disk_gb":    roleRequired,
		"id":         roleComputed,
		"created_at": roleComputed,
		"updated_at": roleComputed,
	},
	"hypervisors": {
		"datacenter_id":    roleRequired,
		"name":             roleRequired,
		"cpu_total":        roleRequired,
		"cpu_reserved":     roleRequired,
		"ram_gb_total":     roleRequired,
		"ram_gb_reserved":  roleRequired,
		"disk_gb_total":    roleRequired,
		"disk_gb_reserved": roleRequired,
		"schedulable":      roleOptionalComputed,
		// Optional+computed: omitting it inherits the datacenter's ratio, which
		// the server fills in — a plain Optional would diff against that.
		"cpu_overcommit_ratio": roleOptionalComputed,
		"cpu_used":             roleComputed,
		"cpu_effective_total":  roleComputed,
		"cpu_bookable":         roleComputed,
		"ram_gb_used":          roleComputed,
		"ram_gb_bookable":      roleComputed,
		"disk_gb_used":         roleComputed,
		"disk_gb_bookable":     roleComputed,
		"last_synced_at":       roleComputed,
		"id":                   roleComputed,
		"created_at":           roleComputed,
		"updated_at":           roleComputed,
	},
	"pools": {
		"name":          roleRequired,
		"datacenter_id": roleRequired,
		"slot_id":       roleRequired,
		"desired_count": roleRequired,
		"metadata":      roleOptional,
		"id":            roleComputed,
		"created_at":    roleComputed,
		"updated_at":    roleComputed,
	},
	// placements_resource.go is a full hand-authored overlay (writeProviderOverlays
	// runs before patchResourceSchemaRoles). The schema roles here serve two
	// purposes: (1) patchClientReadOnlyOmitempty adds omitempty to the generated
	// PlacementView client model for server-assigned fields; (2) patchResourceSchemaRoles
	// validates idempotently that the overlay already has the correct roles.
	"placements": {
		"placement_id":    roleRequired,
		"id":              roleComputed,
		"vmid":            roleOptional,
		"pool_id":         roleComputed,
		"hypervisor_id":   roleComputed,
		"hypervisor_name": roleComputed,
		"created_at":      roleComputed,
	},
}

// patchResourceSchemaRoles rewrites the flag line of each classified attribute
// in the generated <resource>_resource.go Schema() method. It runs on freshly
// generated files (every attribute is Required), so a single per-field flag
// line is replaced. Errors loudly if an expected attribute is missing.
func patchResourceSchemaRoles(outDir string) error {
	providerDir := filepath.Join(outDir, "internal", "provider")
	for resource, fields := range resourceSchemaRoles {
		file := filepath.Join(providerDir, resource+"_resource.go")
		b, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read %s: %w", file, err)
		}
		content := string(b)
		for field, role := range fields {
			re := regexp.MustCompile(`("` + regexp.QuoteMeta(field) + `":\s*schema\.\w+Attribute\{\s*)(?:Required|Optional|Computed):\s+true,`)
			if !re.MatchString(content) {
				return fmt.Errorf("%s: attribute %q not found; generator output may have changed", file, field)
			}
			content = re.ReplaceAllString(content, "${1}"+role.flags())
		}
		if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", file, err)
		}
	}
	log.Printf("patched resource schema roles (server-assigned fields -> Computed)")
	return nil
}

// resourceImmutableAttrs declares, per resource, the attributes that get a
// RequiresReplace plan modifier because the API refuses to change them after
// create.
//
// This table is a tripwire, not the source of truth: specImmutableAttrs
// recomputes the same sets from the OpenAPI spec on every run and
// patchResourceImmutability fails if the two disagree. Its job is to make a
// change in the API's write surface show up as a build failure here rather than
// as a silently un-replaced attribute in the provider.
//
// A resource is listed with an empty set when create and update accept the same
// body — that is a positive assertion that nothing is immutable, and it is what
// makes a newly-narrowed update endpoint fail loudly.
//
// placements is absent on purpose: the API has no create endpoint for it
// (placements are created by pools), so immutability cannot be derived. The
// hand-authored overlay declares RequiresReplace on placement_id directly.
var resourceImmutableAttrs = map[string][]string{
	"datacenters":   {},
	"hypervisors":   {},
	"organizations": {},
	"slots":         {},
	// PATCH /pools/{id} binds only desired_count (resizePoolInput), and
	// FleetService.ResizePool never touches the rest of the row. Without
	// RequiresReplace, Terraform plans an in-place update the provider cannot
	// perform, the PATCH response echoes the unchanged values back over the
	// plan, and apply fails with "Provider produced inconsistent result".
	"pools": {"datacenter_id", "metadata", "name", "slot_id"},
}

// planModifierPkgs maps a generated schema attribute type to the plugin
// framework subpackage holding its plan modifiers. An attribute type missing
// here is a hard error rather than a skip, so a new column type cannot quietly
// lose its RequiresReplace.
var planModifierPkgs = map[string]string{
	"Bool":    "boolplanmodifier",
	"Float64": "float64planmodifier",
	"Int64":   "int64planmodifier",
	"List":    "listplanmodifier",
	"Map":     "mapplanmodifier",
	"Number":  "numberplanmodifier",
	"Object":  "objectplanmodifier",
	"Set":     "setplanmodifier",
	"String":  "stringplanmodifier",
}

// schemaAttrRe matches the opening of a generated schema attribute block:
//
//	"slot_id": schema.StringAttribute{
//
// Group 1 is the attribute name, group 2 the framework type.
var schemaAttrRe = regexp.MustCompile(`"(\w+)":\s*schema\.(\w+)Attribute\{`)

// oasSchema is the sliver of an OpenAPI schema object specImmutableAttrs needs:
// either a reference into components/schemas or an inline property set.
type oasSchema struct {
	Ref        string                     `json:"$ref"`
	Properties map[string]json.RawMessage `json:"properties"`
}

type oasOperation struct {
	Tags        []string `json:"tags"`
	RequestBody *struct {
		Content map[string]struct {
			Schema oasSchema `json:"schema"`
		} `json:"content"`
	} `json:"requestBody"`
}

type oasDocument struct {
	Paths      map[string]map[string]json.RawMessage `json:"paths"`
	Components struct {
		Schemas map[string]oasSchema `json:"schemas"`
	} `json:"components"`
}

// specImmutableAttrs derives, per resource, the request-body properties that a
// client can set at create time but can never change afterwards: those present
// in the create body and absent from the update body.
//
// The mapping from spec to resource follows what the OpenAPI Generator does —
// operations are grouped by their first tag, and the tag becomes the resource
// file name (api-keys -> api_keys_resource.go). Within a tag group the create
// operation is the POST on the bare collection path (/pools) and the update
// operation is the PUT or PATCH on that path plus one path parameter
// (/pools/{id}). Anything that does not fit that shape — /auth/login, which is
// a POST but not a collection; /organizations/{id}/members/{userId}, which is
// nested a level deeper — is not a create/update pair and is ignored.
//
// A tag with a create and no matching update has every create property
// immutable: there is no way to change any of them.
func specImmutableAttrs(specJSON []byte) (map[string][]string, error) {
	var doc oasDocument
	if err := json.Unmarshal(specJSON, &doc); err != nil {
		return nil, fmt.Errorf("parse openapi spec: %w", err)
	}

	resolve := func(s oasSchema) (map[string]json.RawMessage, error) {
		if s.Ref != "" {
			name := s.Ref[strings.LastIndex(s.Ref, "/")+1:]
			target, ok := doc.Components.Schemas[name]
			if !ok {
				return nil, fmt.Errorf("unresolved schema ref %q", s.Ref)
			}
			return target.Properties, nil
		}
		return s.Properties, nil
	}

	bodyProps := func(op oasOperation) (map[string]json.RawMessage, error) {
		if op.RequestBody == nil {
			return nil, nil
		}
		content, ok := op.RequestBody.Content["application/json"]
		if !ok {
			return nil, nil
		}
		return resolve(content.Schema)
	}

	type opAt struct {
		path string
		op   oasOperation
	}
	creates := map[string][]opAt{}
	updates := map[string][]opAt{}

	for path, methods := range doc.Paths {
		for method, raw := range methods {
			method = strings.ToLower(method)
			if method != "post" && method != "put" && method != "patch" {
				continue
			}
			var op oasOperation
			if err := json.Unmarshal(raw, &op); err != nil {
				return nil, fmt.Errorf("parse %s %s: %w", strings.ToUpper(method), path, err)
			}
			if len(op.Tags) == 0 || op.RequestBody == nil {
				continue
			}
			tag := op.Tags[0]
			segs := strings.Split(strings.Trim(path, "/"), "/")
			switch {
			case method == "post" && len(segs) == 1 && !strings.Contains(path, "{"):
				creates[tag] = append(creates[tag], opAt{path, op})
			case method != "post" && len(segs) == 2 && !strings.Contains(segs[0], "{") &&
				strings.HasPrefix(segs[1], "{") && strings.HasSuffix(segs[1], "}"):
				updates[tag] = append(updates[tag], opAt{path, op})
			}
		}
	}

	out := map[string][]string{}
	for tag, cs := range creates {
		if len(cs) != 1 {
			// Ambiguous: more than one collection POST under this tag, so
			// there is no single "the create endpoint" to compare against.
			continue
		}
		create := cs[0]

		var update *oasOperation
		for _, u := range updates[tag] {
			if strings.HasPrefix(u.path, create.path+"/{") {
				if update != nil {
					update = nil // ambiguous; treat as underivable
					break
				}
				o := u.op
				update = &o
			}
		}

		createProps, err := bodyProps(create.op)
		if err != nil {
			return nil, fmt.Errorf("%s create body: %w", tag, err)
		}
		var updateProps map[string]json.RawMessage
		if update != nil {
			if updateProps, err = bodyProps(*update); err != nil {
				return nil, fmt.Errorf("%s update body: %w", tag, err)
			}
		}

		immutable := []string{}
		for name := range createProps {
			// Huma injects $schema into every body; it is metadata, and
			// dropTerraformSchemaField removes it from the provider entirely.
			if name == "$schema" {
				continue
			}
			if _, mutable := updateProps[name]; !mutable {
				immutable = append(immutable, name)
			}
		}
		sort.Strings(immutable)
		out[strings.ReplaceAll(tag, "-", "_")] = immutable
	}
	return out, nil
}

// patchResourceImmutability gives every server-immutable attribute a
// RequiresReplace plan modifier, so Terraform destroys and recreates the
// resource instead of planning an update the API cannot carry out.
//
// It must run after patchResourceSchemaRoles: it injects at the attribute's
// opening brace, which sits between the anchor and the flag line that pass
// rewrites.
func patchResourceImmutability(outDir, specPath string) error {
	specJSON, err := os.ReadFile(specPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", specPath, err)
	}
	derived, err := specImmutableAttrs(specJSON)
	if err != nil {
		return err
	}

	// Walk the union so a stale table entry for a resource the spec no longer
	// describes is caught too.
	resources := map[string]bool{}
	for name := range derived {
		resources[name] = true
	}
	for name := range resourceImmutableAttrs {
		resources[name] = true
	}

	providerDir := filepath.Join(outDir, "internal", "provider")
	for _, resource := range sortedKeys(resources) {
		declared := resourceImmutableAttrs[resource]
		immutable, inSpec := derived[resource]
		if !inSpec {
			return fmt.Errorf(
				"resourceImmutableAttrs declares %q, but the spec has no derivable create/update pair for it; "+
					"remove the entry or fix the derivation", resource,
			)
		}
		if _, managed := resourceSchemaRoles[resource]; !managed {
			if len(immutable) > 0 {
				log.Printf(
					"note: %s has spec-immutable attributes %v but no schema-role classification; not patched",
					resource, immutable,
				)
			}
			continue
		}

		file := filepath.Join(providerDir, resource+"_resource.go")
		b, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read %s: %w", file, err)
		}
		content := string(b)

		var applied []string
		var kinds []string
		var walkErr error
		patched := schemaAttrRe.ReplaceAllStringFunc(content, func(m string) string {
			sub := schemaAttrRe.FindStringSubmatch(m)
			field, kind := sub[1], sub[2]
			if !slices.Contains(immutable, field) {
				return m
			}
			pkg, ok := planModifierPkgs[kind]
			if !ok {
				walkErr = fmt.Errorf(
					"%s: attribute %q is schema.%sAttribute, which has no entry in planModifierPkgs",
					file, field, kind,
				)
				return m
			}
			applied = append(applied, field)
			if !slices.Contains(kinds, kind) {
				kinds = append(kinds, kind)
			}
			return m + fmt.Sprintf("\nPlanModifiers: []planmodifier.%s{%s.RequiresReplace()},", kind, pkg)
		})
		if walkErr != nil {
			return walkErr
		}

		sort.Strings(applied)
		if applied == nil {
			applied = []string{}
		}
		if !slices.Equal(applied, declared) {
			return fmt.Errorf(
				"%s: the spec makes %v immutable and %v of those are schema attributes, "+
					"but resourceImmutableAttrs declares %v; reconcile the table with the API",
				file, immutable, applied, declared,
			)
		}
		if len(applied) == 0 {
			continue
		}

		imports := []string{"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"}
		for _, kind := range kinds {
			imports = append(imports, "github.com/hashicorp/terraform-plugin-framework/resource/schema/"+planModifierPkgs[kind])
		}
		patched, err = addProviderImports(patched, imports)
		if err != nil {
			return fmt.Errorf("%s: %w", file, err)
		}
		if err := os.WriteFile(file, []byte(patched), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", file, err)
		}
		log.Printf("patched %s: RequiresReplace on %v", file, applied)
	}
	return nil
}

// addProviderImports inserts import paths into a generated provider file's
// import block, after the terraform-plugin-framework schema import the
// generator always emits. Already-present paths are left alone.
func addProviderImports(content string, paths []string) (string, error) {
	const anchor = "\t\"github.com/hashicorp/terraform-plugin-framework/resource/schema\"\n"
	var add []string
	for _, p := range paths {
		if !strings.Contains(content, `"`+p+`"`) {
			add = append(add, "\t\""+p+"\"\n")
		}
	}
	if len(add) == 0 {
		return content, nil
	}
	if !strings.Contains(content, anchor) {
		return content, fmt.Errorf("schema import anchor not found; generator output may have changed")
	}
	return strings.Replace(content, anchor, anchor+strings.Join(add, ""), 1), nil
}

// sortedKeys returns a map's keys in sorted order, so generator passes emit a
// stable ordering of log lines and errors run to run.
func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// writeProviderOverlays drops hand-authored provider source files into the
// generated provider and registers them in provider.go. These cover resources
// and data sources the OpenAPI Generator cannot model (e.g. lists of nested
// objects). It runs on every regenerate, so the overlays are always present.
func writeProviderOverlays(outDir string) error {
	providerDir := filepath.Join(outDir, "internal", "provider")

	dst := filepath.Join(providerDir, "pool_placements_data_source.go")
	if err := os.WriteFile(dst, []byte(poolPlacementsDataSourceSrc), 0o644); err != nil {
		return fmt.Errorf("write overlay %s: %w", dst, err)
	}

	if err := registerDataSource(filepath.Join(providerDir, "provider.go"), "NewPoolPlacementsDataSource"); err != nil {
		return err
	}

	// Overwrite the generated by-id-only slots data source with the id-or-name
	// overlay. The generator already declares NewSlotsDataSource and registers
	// it in provider.go, so no separate registration is needed here.
	slotsDst := filepath.Join(providerDir, "slots_data_source.go")
	if err := os.WriteFile(slotsDst, []byte(slotsDataSourceSrc), 0o644); err != nil {
		return fmt.Errorf("write overlay %s: %w", slotsDst, err)
	}

	// Overwrite the generated placements_resource.go with the hand-authored
	// overlay that implements full CRUD (adopt/backfill-vmid/delete). The
	// generator produces a correct stub (no Create, schema from view model)
	// but the overlay provides a richer schema and proper Create semantics.
	// The generator already declares NewPlacementsResource and registers it in
	// provider.go, so no separate registration is needed here.
	placementsDst := filepath.Join(providerDir, "placements_resource.go")
	if err := os.WriteFile(placementsDst, []byte(placementResourceSrc), 0o644); err != nil {
		return fmt.Errorf("write overlay %s: %w", placementsDst, err)
	}

	clientTestDst := filepath.Join(outDir, "internal", "client", "auth_body_test.go")
	if err := os.WriteFile(clientTestDst, []byte(clientBodyTestSrc), 0o644); err != nil {
		return fmt.Errorf("write overlay %s: %w", clientTestDst, err)
	}

	log.Printf("injected provider overlays: pool_placements_data_source.go, slots_data_source.go, placements_resource.go, client/auth_body_test.go")
	return nil
}

// registerDataSource inserts `<constructor>,` into the slice returned by the
// provider's DataSources() method, idempotently. Errors if the expected
// generator output shape is missing (so a template change is caught loudly).
func registerDataSource(providerFile, constructor string) error {
	b, err := os.ReadFile(providerFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", providerFile, err)
	}
	content := string(b)
	if strings.Contains(content, constructor) {
		return nil // already registered
	}
	const anchor = "return []func() datasource.DataSource{\n"
	if !strings.Contains(content, anchor) {
		return fmt.Errorf("%s: DataSources() anchor not found; generator output may have changed", providerFile)
	}
	patched := strings.Replace(content, anchor, anchor+"\t\t"+constructor+",\n", 1)
	if err := os.WriteFile(providerFile, []byte(patched), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", providerFile, err)
	}
	return nil
}

// copyTreeOverwrite recursively copies src into dst, creating directories as
// needed and overwriting any existing files (unlike os.CopyFS, which errors on
// pre-existing files). Used by syncProviderRepo to lay generator output over
// the provider checkout.
func copyTreeOverwrite(dst, src string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, b, 0o644)
	})
}

// patchTerraformDollarIdentifiers rewrites illegal Go identifiers of the form
// "$Name" (a dollar sign followed by an uppercase letter) into "Name" across
// all generated .go files under outDir. The only source of such tokens is the
// terraform-provider template emitting Huma's "$schema" property as a struct
// field and its references. The match is intentionally restricted to "$" + an
// uppercase letter so it never touches the lowercase "$schema" that legitimately
// appears inside `json:"$schema,omitempty"` struct tags.
// dropTerraformSchemaField removes Huma's injected "$schema" metadata field from
// the generated provider tree. The generator renders it as a model field tagged
// tfsdk:"__schema" (illegal — tfsdk tags must start with a letter), a schema
// attribute block, and to/from client mapping lines. Removed rather than renamed
// because it is metadata, not a real resource attribute.
func dropTerraformSchemaField(outDir string) error {
	res := []*regexp.Regexp{
		// struct field:  Schema  types.String `tfsdk:"__schema"`
		regexp.MustCompile("(?m)^[ \t]*Schema[ \t]+types\\.String[ \t]+`tfsdk:\"__schema\"`[ \t]*\n"),
		// schema attribute block:  "__schema": schema.StringAttribute{ ... },
		regexp.MustCompile("(?s)[ \t]*\"__schema\": schema\\.StringAttribute\\{.*?\\},\n"),
		// ToClientModel guard block
		regexp.MustCompile("(?s)[ \t]*if !m\\.Schema\\.IsNull\\(\\) && !m\\.Schema\\.IsUnknown\\(\\) \\{\n[ \t]*out\\.Schema = m\\.Schema\\.ValueString\\(\\)\n[ \t]*\\}\n"),
		// FromClientModel assignment
		regexp.MustCompile("(?m)^[ \t]*m\\.Schema = types\\.StringValue\\(c\\.Schema\\)\n"),
	}

	patched := 0
	err := filepath.WalkDir(outDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		out := b
		for _, re := range res {
			out = re.ReplaceAll(out, nil)
		}
		if string(out) == string(b) {
			return nil
		}
		if err := os.WriteFile(path, out, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
		log.Printf("patched %s: removed $schema metadata field", path)
		patched++
		return nil
	})
	if err != nil {
		return err
	}
	if patched == 0 {
		log.Printf("note: no $schema field found under %s; skipping drop", outDir)
	}
	return nil
}

// patchTerraformAPIKeyBearer makes the generated client send the configured
// api_key as a Bearer token. The OpenAPI Generator emits a raw
// `Authorization: <api_key>` header, but Waggle's only auth scheme is HTTP
// bearer (a wgl_ org key is recognized as a bearer token), so without the
// prefix the API returns 401 "missing bearer token".
func patchTerraformAPIKeyBearer(outDir string) error {
	path := filepath.Join(outDir, "internal", "client", "client.go")
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	patched := strings.ReplaceAll(string(b),
		`req.Header.Set("Authorization", c.ApiKey)`,
		`req.Header.Set("Authorization", "Bearer "+c.ApiKey)`)
	if patched == string(b) {
		log.Printf("note: api_key Authorization line not found in %s; skipping", path)
		return nil
	}
	if err := os.WriteFile(path, []byte(patched), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	log.Printf("patched %s: send api_key as Bearer token", path)
	return nil
}

// patchClientReadOnlyOmitempty adds ",omitempty" to the json tags of read-only
// (server-assigned) fields in the generated client models. Create/update bodies
// are built from the response view models, so without this the empty read-only
// fields (id, created_at, ...) are marshaled into the request and the strict
// Huma API rejects them with 422 "unexpected property". The set is derived from
// resourceSchemaRoles so it stays in sync with the schema-role classification.
func patchClientReadOnlyOmitempty(outDir string) error {
	readOnly := map[string]bool{}
	for _, fields := range resourceSchemaRoles {
		for name, role := range fields {
			if role == roleComputed {
				readOnly[name] = true
			}
		}
	}

	clientDir := filepath.Join(outDir, "internal", "client")
	patched := 0
	err := filepath.WalkDir(clientDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		s := string(b)
		orig := s
		for name := range readOnly {
			s = strings.ReplaceAll(s, `json:"`+name+`"`, `json:"`+name+`,omitempty"`)
		}
		if s == orig {
			return nil
		}
		if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
		log.Printf("patched %s: omitempty on read-only fields", path)
		patched++
		return nil
	})
	if err != nil {
		return err
	}
	if patched == 0 {
		log.Printf("note: no read-only json tags found under %s; skipping", clientDir)
	}
	return nil
}

func patchTerraformDollarIdentifiers(outDir string) error {
	dollarIdent := regexp.MustCompile(`\$([A-Z])`)

	patched := 0
	err := filepath.WalkDir(outDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if !strings.ContainsRune(string(b), '$') {
			return nil
		}

		fixed := dollarIdent.ReplaceAll(b, []byte("$1"))
		if string(fixed) == string(b) {
			return nil
		}
		if err := os.WriteFile(path, fixed, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
		log.Printf("patched %s: rewrote illegal $-prefixed Go identifiers", path)
		patched++
		return nil
	})
	if err != nil {
		return err
	}
	if patched == 0 {
		log.Printf("note: no $-prefixed Go identifiers found under %s; skipping patch", outDir)
	}
	return nil
}

func generateOpenAPISpec(outPath string) error {
	log.Println("Extracting OpenAPI Spec from Huma...")

	b, err := buildOpenAPISpecBytes()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(outPath, b, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outPath, err)
	}

	return nil
}

// buildOpenAPISpecBytes builds the Huma app with non-nil (but dependency-less)
// dummy services and marshals its OpenAPI document to JSON. Constructing every
// service is what makes each route group register and land in the emitted spec;
// handlers never run during extraction, so nil repos/DB are fine — only
// registration matters. Shared by the SDK and Terraform generators so both emit
// the same complete spec.
func buildOpenAPISpecBytes() ([]byte, error) {
	// Stamp the release version HERE rather than in the callers: the spec takes
	// its version from buildinfo.Version, which is "dev" under `go run`, and
	// both `generate sdk` and `generate terraform` write docs/openapi.json.
	// Stamping in one caller only meant the other silently overwrote the file
	// with "dev" -- which is exactly what `just release-prep` (sdk then
	// terraform) did on its first run.
	buildinfo.Version = releaseVersion()

	cfg, err := config.Load()
	if err != nil {
		log.Printf("config.Load() failed (%v); using defaults for spec generation", err)
		cfg = &config.Config{
			BindHost:          "127.0.0.1",
			BindPort:          "8080",
			BasePath:          "/api/v1",
			BaseURL:           "http://localhost:8080",
			FrontendMode:      "none",
			ViteDevURL:        "http://localhost:5173",
			DatabaseURL:       "postgres://user:pass@localhost:5432/db?sslmode=disable",
			Env:               "development",
			JWTSecret:         "spec-only-secret",
			JWTIssuer:         "nexus",
			JWTAccessTTLMin:   60,
			JWTRefreshTTLHour: 720,
		}
	}
	cfg.FrontendMode = "none"

	if strings.TrimSpace(cfg.JWTSecret) == "" {
		cfg.JWTSecret = "spec-only-secret"
	}
	if strings.TrimSpace(cfg.JWTIssuer) == "" {
		cfg.JWTIssuer = "nexus"
	}
	if cfg.JWTAccessTTLMin <= 0 {
		cfg.JWTAccessTTLMin = 60
	}
	if cfg.JWTRefreshTTLHour <= 0 {
		cfg.JWTRefreshTTLHour = 720
	}

	dummyClient, err := river.NewClient(riverpgxv5.New(nil), &river.Config{
		Workers: river.NewWorkers(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create dummy river client: %w", err)
	}

	dummyTokenSvc, _ := service.NewTokenService(cfg.JWTSecret, cfg.JWTIssuer, 0, 0, cfg.JWTAudience)
	dummySessionSvc := service.NewTokenSessionService(dummyTokenSvc, nil, nil, nil)
	dummyAuthSvc := service.NewAuthService(nil, dummyTokenSvc, nil, nil, nil, nil, nil, nil, service.LogSender{})
	dummyFleetSvc := service.NewFleetService(nil, service.ReservationDefaults{})
	dummyAPIKeySvc := service.NewAPIKeyService(nil)
	dummyOrgSvc := service.NewOrgService(nil, dummyTokenSvc, nil, service.LogSender{})

	apiApp, err := api.Build(*cfg, &app.Deps{
		Config:        cfg,
		Jobs:          nil,
		River:         dummyClient,
		TokenSessions: dummySessionSvc,
		Tokens:        dummyTokenSvc,
		Auth:          dummyAuthSvc,
		Fleet:         dummyFleetSvc,
		APIKeys:       dummyAPIKeySvc,
		Orgs:          dummyOrgSvc,
	})
	if err != nil {
		return nil, fmt.Errorf("api build: %w", err)
	}

	b, err := apiApp.API.OpenAPI().MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal openapi json: %w", err)
	}

	return b, nil
}
func init() {
	generateTerraformOAGCmd.Flags().StringVar(
		&providerOutDir,
		"out",
		defaultProviderOutDir(),
		"checkout of github.com/GlueOps/terraform-provider-waggle to generate into",
	)

	generateTerraformCmd.AddCommand(generateTerraformHashiCmd)
	generateTerraformCmd.AddCommand(generateTerraformOAGCmd)

	generateCmd.AddCommand(generateMigrationsCmd)
	generateCmd.AddCommand(generateSdkCmd)
	generateCmd.AddCommand(generateTerraformCmd)
	rootCmd.AddCommand(generateCmd)
}

// patchPoolsResourceUpdate fixes the generated pools Update to use state.Id
// for the PATCH URL instead of plan.Id, which is "(known after apply)" during
// an update and produces an empty path (/pools/) that returns an HTML error.
func patchPoolsResourceUpdate(outDir string) error {
	file := filepath.Join(outDir, "internal", "provider", "pools_resource.go")
	b, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read %s: %w", file, err)
	}
	content := string(b)

	const old = `func (r *PoolsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PoolsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := plan.ToClientModel()

	respBody, err := r.client.DoRequest(ctx, "PATCH", fmt.Sprintf("/pools/%v", plan.Id.ValueString()), reqBody)`

	const patched = `func (r *PoolsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PoolsModel
	var state PoolsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resize endpoint only accepts desired_count — not the full pool body.
	resizeBody := map[string]int64{"desired_count": plan.DesiredCount.ValueInt64()}

	// Use state.Id (known current value) not plan.Id (unknown during update).
	respBody, err := r.client.DoRequest(ctx, "PATCH", fmt.Sprintf("/pools/%v", state.Id.ValueString()), resizeBody)`

	if !strings.Contains(content, old) {
		log.Printf("note: pools_resource.go Update already patched or generator output changed; skipping")
		return nil
	}
	patched2 := strings.Replace(content, old, patched, 1)
	if err := os.WriteFile(file, []byte(patched2), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", file, err)
	}
	log.Printf("patched pools_resource.go: Update uses state.Id for PATCH URL")
	return nil
}

// patchPoolsMetadataMapping wires up pools.metadata, which the generator
// advertises as a schema attribute but maps in neither direction.
//
// metadata is json.RawMessage in the API, so the spec types it as a free-form
// object, the OpenAPI Generator renders it as interface{} in client.PoolView,
// and the model generator — which only knows how to emit types.String/Int64/
// Bool conversions — skips it. The attribute stays in the schema, so Terraform
// accepts it in config and (being Optional and untouched) keeps whatever was
// planned, while the value is silently dropped on create and never refreshed on
// read.
//
// The fix types it as jsontypes.Normalized rather than a plain string: two JSON
// documents differing only in key order or whitespace are then semantically
// equal, so a round-trip through the API does not read as a diff.
func patchPoolsMetadataMapping(outDir string) error {
	providerDir := filepath.Join(outDir, "internal", "provider")

	modelFile := filepath.Join(providerDir, "pools_model.go")
	b, err := os.ReadFile(modelFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", modelFile, err)
	}
	model := string(b)

	fieldRe := regexp.MustCompile("(?m)^\tMetadata\\s+types\\.String\\s+`tfsdk:\"metadata\"`$")
	if !fieldRe.MatchString(model) {
		return fmt.Errorf("%s: Metadata field not found in the expected shape; generator output may have changed", modelFile)
	}
	model = fieldRe.ReplaceAllString(model, "\tMetadata jsontypes.Normalized `tfsdk:\"metadata\"`")

	const modelImports = "import (\n\t\"github.com/hashicorp/terraform-plugin-framework/types\"\n"
	if !strings.Contains(model, modelImports) {
		return fmt.Errorf("%s: import block not found in the expected shape; generator output may have changed", modelFile)
	}
	model = strings.Replace(model, modelImports,
		"import (\n\t\"encoding/json\"\n\n"+
			"\t\"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes\"\n"+
			"\t\"github.com/hashicorp/terraform-plugin-framework/types\"\n", 1)

	toRe := regexp.MustCompile(`(?s)(func \(m \*PoolsModel\) ToClientModel\(\) \*client\.PoolView \{.*?)(\n\treturn out\n\})`)
	if !toRe.MatchString(model) {
		return fmt.Errorf("%s: ToClientModel not found in the expected shape; generator output may have changed", modelFile)
	}
	model = toRe.ReplaceAllString(model, "${1}\n"+
		"\tif !m.Metadata.IsNull() && !m.Metadata.IsUnknown() {\n"+
		"\t\t// json.RawMessage marshals through the interface{} field verbatim.\n"+
		"\t\tout.Metadata = json.RawMessage(m.Metadata.ValueString())\n"+
		"\t}${2}")

	fromRe := regexp.MustCompile(`(?s)(func \(m \*PoolsModel\) FromClientModel\(c \*client\.PoolView\) \{.*?)(\n\})`)
	if !fromRe.MatchString(model) {
		return fmt.Errorf("%s: FromClientModel not found in the expected shape; generator output may have changed", modelFile)
	}
	model = fromRe.ReplaceAllString(model, "${1}\n"+
		"\tm.Metadata = jsontypes.NewNormalizedNull()\n"+
		"\tif c.Metadata != nil {\n"+
		"\t\tif b, err := json.Marshal(c.Metadata); err == nil {\n"+
		"\t\t\tm.Metadata = jsontypes.NewNormalizedValue(string(b))\n"+
		"\t\t}\n"+
		"\t}${2}")

	if err := os.WriteFile(modelFile, []byte(model), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", modelFile, err)
	}

	resourceFile := filepath.Join(providerDir, "pools_resource.go")
	b, err = os.ReadFile(resourceFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", resourceFile, err)
	}
	resourceSrc := string(b)

	attrRe := regexp.MustCompile(`"metadata":\s*schema\.StringAttribute\{`)
	if !attrRe.MatchString(resourceSrc) {
		return fmt.Errorf("%s: metadata attribute not found; generator output may have changed", resourceFile)
	}
	if !strings.Contains(resourceSrc, "jsontypes.NormalizedType{}") {
		resourceSrc = attrRe.ReplaceAllString(resourceSrc, "${0}\nCustomType: jsontypes.NormalizedType{},")
	}
	resourceSrc, err = addProviderImports(resourceSrc, []string{"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"})
	if err != nil {
		return fmt.Errorf("%s: %w", resourceFile, err)
	}
	if err := os.WriteFile(resourceFile, []byte(resourceSrc), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", resourceFile, err)
	}

	log.Printf("patched pools metadata: jsontypes.Normalized round-trip")
	return nil
}

// patchResourceRead404 adds 404 → RemoveResource handling to every generated
// resource Read method. Without this, a resource deleted outside Terraform
// causes a hard error on the next plan/apply instead of prompting recreation.
// The fix uses isNotFound (defined in placements_resource.go, same package)
// which detects the "API error ... 404" string emitted by DoRequest.
//
// Indentation: inside a function body the generator uses one tab; inside an
// if block, two tabs. The matched block is:
//
//	\t\tresp.Diagnostics.AddError("Error reading X", err.Error())
//	\t\treturn
//	\t}
//
// The replacement prepends the isNotFound guard before the matched block so
// the closing \t} of "if err != nil" is preserved.
func patchResourceRead404(outDir string) error {
	providerDir := filepath.Join(outDir, "internal", "provider")
	// Matches the AddError+return+} block the generator emits inside Read.
	// Two tabs before AddError/return (inside if err != nil), one tab for }.
	re := regexp.MustCompile(
		`\t\tresp\.Diagnostics\.AddError\("Error reading [^"]+", err\.Error\(\)\)\n\t\treturn\n\t\}`,
	)
	patched := 0
	err := filepath.WalkDir(providerDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, "_resource.go") {
			return nil
		}
		// placements_resource.go is a hand-authored overlay that already handles
		// 404 correctly — skip it to avoid a double-injection.
		if strings.HasSuffix(path, "placements_resource.go") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		s := string(b)
		if !re.MatchString(s) {
			return nil
		}
		// Prepend the isNotFound guard before the matched block.
		// The \t} at the end of match closes the surrounding "if err != nil".
		fixed := re.ReplaceAllStringFunc(s, func(match string) string {
			return "\t\tif isNotFound(err) {\n\t\t\tresp.State.RemoveResource(ctx)\n\t\t\treturn\n\t\t}\n" + match
		})
		if fixed == s {
			return nil
		}
		if err := os.WriteFile(path, []byte(fixed), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
		log.Printf("patched %s: Read returns RemoveResource on 404", path)
		patched++
		return nil
	})
	if err != nil {
		return err
	}
	if patched == 0 {
		log.Printf("note: no unpatched resource Read 404 handlers found under %s; skipping", providerDir)
	}
	return nil
}

// semverish matches a bare MAJOR.MINOR.PATCH with optional pre-release/build
// metadata. Used to reject values like "dev" that npm and the Go SDK
// generator cannot use as a package version.
var semverish = regexp.MustCompile(`^\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?$`)

// releaseVersion is the version stamped into the generated spec and SDKs.
//
// It reads the VERSION file rather than `git describe`, because git describe
// returns the LAST tag reachable from HEAD -- so regenerating while preparing
// release N+1 would stamp N into artifacts that then ship inside the vN+1 tag,
// leaving every committed artifact one release behind. VERSION is bumped as
// the first step of a release, so `just sdk` before tagging produces artifacts
// that already carry the version the tag will use.
//
// This governs committed artifacts only. The RUNNING binary reports
// buildinfo.Version, injected from the git tag via -ldflags in the Dockerfile,
// so what the UI shows is always the real build regardless of this file.
func releaseVersion() string {
	const fallback = "0.0.0"
	b, err := os.ReadFile("VERSION")
	if err != nil {
		log.Printf("could not read VERSION (%v); using %s", err, fallback)
		return fallback
	}
	v := strings.TrimPrefix(strings.TrimSpace(string(b)), "v")
	if !semverish.MatchString(v) {
		log.Printf("VERSION %q is not semver; using %s", v, fallback)
		return fallback
	}
	return v
}
