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
	"strings"

	"github.com/glueops/waggle/internal/api"
	"github.com/glueops/waggle/internal/app"
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

		specVersion := "0.0.0"
		{
			var doc struct {
				Info struct {
					Version string `json:"version"`
				} `json:"info"`
			}
			if json.Unmarshal(b, &doc) == nil && strings.TrimSpace(doc.Info.Version) != "" {
				specVersion = strings.TrimPrefix(strings.TrimSpace(doc.Info.Version), "v")
			}
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

	// Assumes openapi-ts.config.ts or openapi-ts.config.mjs exists at repo root.
	// Pin exact version on purpose.
	cmd := exec.Command("npx", "--yes", "@hey-api/openapi-ts@0.93.0")
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
		outDir := "terraform-provider-waggle"

		if err := generateOpenAPISpec(specPath); err != nil {
			return err
		}

		// Preserve hand-authored examples across the clean+regenerate. The
		// generator owns only examples/provider/provider.tf; everything else
		// under examples/ is maintained by hand and would otherwise be lost to
		// the RemoveAll below. Snapshot it now and overlay it back after
		// generation (the restore wins on conflict, so curated examples persist).
		examplesDir := filepath.Join(outDir, "examples")
		snapshot, err := snapshotDir(examplesDir)
		if err != nil {
			return fmt.Errorf("snapshot examples: %w", err)
		}
		if snapshot != "" {
			defer os.RemoveAll(snapshot)
		}

		if err := os.RemoveAll(outDir); err != nil {
			return fmt.Errorf("failed to remove %s: %w", outDir, err)
		}
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return err
		}

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

		// Post-process: Huma adds a "$schema" metadata property to response
		// bodies. The OpenAPI Generator terraform-provider templates render it as
		// a model field tagged tfsdk:"__schema" — which compiles but the plugin
		// framework rejects at apply time ("invalid tfsdk tag, must start with a
		// letter"). It is metadata, not a real resource attribute, so drop it
		// entirely from the provider models, schema attributes, and mappings.
		if err := dropTerraformSchemaField(outDir); err != nil {
			return err
		}
		// Safety net for any other "$"-prefixed Go identifiers the generator may
		// emit (the $schema field above is removed before this runs).
		if err := patchTerraformDollarIdentifiers(outDir); err != nil {
			return err
		}

		// Correct resource attribute roles: the generator marks every field
		// Required, forcing server-assigned fields (id, created_at, ...) into
		// config and causing perpetual diffs. Reclassify per the spec.
		if err := patchResourceSchemaRoles(outDir); err != nil {
			return err
		}

		// Inject hand-authored provider overlays (data sources the generator
		// can't model) and register them in the generated provider.go.
		if err := writeProviderOverlays(outDir); err != nil {
			return err
		}

		// gofmt the provider package: the schema-role and registration patches
		// emit loosely-spaced Go; gofmt normalizes alignment/indentation.
		fmtProvider := exec.Command("gofmt", "-w", ".")
		fmtProvider.Dir = filepath.Join(outDir, "internal", "provider")
		if err := run(fmtProvider, "gofmt provider -> "+outDir); err != nil {
			return err
		}

		// Restore hand-authored examples, overlaying (and winning over) anything
		// the generator wrote under examples/.
		if snapshot != "" {
			if err := copyTreeOverwrite(examplesDir, snapshot); err != nil {
				return fmt.Errorf("restore examples: %w", err)
			}
			log.Printf("restored hand-authored examples under %s", examplesDir)
		}

		log.Println("OpenAPI Generator Terraform provider generated successfully.")
		return nil
	},
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
		"cpu_used":         roleComputed,
		"cpu_bookable":     roleComputed,
		"ram_gb_used":      roleComputed,
		"ram_gb_bookable":  roleComputed,
		"disk_gb_used":     roleComputed,
		"disk_gb_bookable": roleComputed,
		"last_synced_at":   roleComputed,
		"id":               roleComputed,
		"created_at":       roleComputed,
		"updated_at":       roleComputed,
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

	log.Printf("injected provider overlay pool_placements_data_source.go and registered it")
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

// snapshotDir copies srcDir into a fresh temp directory and returns its path,
// so the caller can restore it after a destructive regenerate. Returns "" (no
// error) when srcDir does not exist.
func snapshotDir(srcDir string) (string, error) {
	if _, err := os.Stat(srcDir); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	tmp, err := os.MkdirTemp("", "waggle-tf-examples-")
	if err != nil {
		return "", err
	}
	if err := copyTreeOverwrite(tmp, srcDir); err != nil {
		_ = os.RemoveAll(tmp)
		return "", err
	}
	return tmp, nil
}

// copyTreeOverwrite recursively copies src into dst, creating directories as
// needed and overwriting any existing files (unlike os.CopyFS, which errors on
// pre-existing files). Used to overlay preserved examples on top of generator
// output.
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
	generateTerraformCmd.AddCommand(generateTerraformHashiCmd)
	generateTerraformCmd.AddCommand(generateTerraformOAGCmd)

	generateCmd.AddCommand(generateMigrationsCmd)
	generateCmd.AddCommand(generateSdkCmd)
	generateCmd.AddCommand(generateTerraformCmd)
	rootCmd.AddCommand(generateCmd)
}
