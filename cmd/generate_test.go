package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// pristineTSConfig is what openapi-generator's typescript-fetch template emits
// before patchTSBuildConfig runs. TypeScript 6 rejects it: node10 resolution is
// a deprecation error (TS5107) and the inferred rootDir is an error (TS5011).
const pristineTSConfig = `{
  "compilerOptions": {
    "declaration": true,
    "target": "es6",
    "module": "commonjs",
    "moduleResolution": "node",
    "outDir": "dist",
    "typeRoots": [
      "node_modules/@types"
    ]
  },
  "exclude": [
    "dist",
    "node_modules"
  ]
}
`

const pristinePackageJSON = `{
  "name": "@glueops/waggle-sdk",
  "version": "0.1.0",
  "main": "./dist/index.js",
  "scripts": {
    "build": "tsc && tsc -p tsconfig.esm.json"
  },
  "devDependencies": {
    "typescript": "^4.0 || ^5.0"
  }
}
`

func writeGeneratedSDK(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "tsconfig.json"), []byte(pristineTSConfig), 0o644); err != nil {
		t.Fatalf("write tsconfig.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pristinePackageJSON), 0o644); err != nil {
		t.Fatalf("write package.json: %v", err)
	}
	return dir
}

func readJSON(t *testing.T, path string) map[string]any {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return m
}

func TestPatchTSBuildConfig(t *testing.T) {
	dir := writeGeneratedSDK(t)

	if err := patchTSBuildConfig(dir); err != nil {
		t.Fatalf("patchTSBuildConfig: %v", err)
	}

	tsconfig := readJSON(t, filepath.Join(dir, "tsconfig.json"))
	opts, ok := tsconfig["compilerOptions"].(map[string]any)
	if !ok {
		t.Fatal("compilerOptions missing after patch")
	}

	if got := opts["ignoreDeprecations"]; got != "6.0" {
		t.Fatalf("ignoreDeprecations = %v, want \"6.0\" (TS6 rejects node10 resolution without it)", got)
	}
	if got := opts["rootDir"]; got != "src" {
		t.Fatalf("rootDir = %v, want \"src\" (TS6 requires it to be explicit)", got)
	}

	// The patch must not disturb the settings that define the published output.
	for key, want := range map[string]any{
		"declaration":      true,
		"target":           "es6",
		"module":           "commonjs",
		"moduleResolution": "node",
		"outDir":           "dist",
	} {
		if got := opts[key]; got != want {
			t.Fatalf("compilerOptions[%q] = %v, want %v — patching must not change emit", key, got, want)
		}
	}
	if _, ok := tsconfig["exclude"].([]any); !ok {
		t.Fatal("top-level \"exclude\" was dropped by the patch")
	}

	pkg := readJSON(t, filepath.Join(dir, "package.json"))
	dev, ok := pkg["devDependencies"].(map[string]any)
	if !ok {
		t.Fatal("devDependencies missing after patch")
	}
	if got := dev["typescript"]; got != "^6.0.0" {
		t.Fatalf("typescript devDependency = %v, want \"^6.0.0\"", got)
	}
	// Unrelated package.json fields must survive.
	if got := pkg["name"]; got != "@glueops/waggle-sdk" {
		t.Fatalf("package name = %v, want it preserved", got)
	}
	if scripts, ok := pkg["scripts"].(map[string]any); !ok || scripts["build"] == nil {
		t.Fatal("scripts.build was dropped by the patch")
	}
}

// Running the patch twice must be a no-op, since `generate sdk` may be re-run
// over an already-patched tree.
func TestPatchTSBuildConfigIsIdempotent(t *testing.T) {
	dir := writeGeneratedSDK(t)

	if err := patchTSBuildConfig(dir); err != nil {
		t.Fatalf("first patch: %v", err)
	}
	first, err := os.ReadFile(filepath.Join(dir, "tsconfig.json"))
	if err != nil {
		t.Fatalf("read tsconfig.json: %v", err)
	}
	firstPkg, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatalf("read package.json: %v", err)
	}

	if err := patchTSBuildConfig(dir); err != nil {
		t.Fatalf("second patch: %v", err)
	}
	second, err := os.ReadFile(filepath.Join(dir, "tsconfig.json"))
	if err != nil {
		t.Fatalf("read tsconfig.json: %v", err)
	}
	secondPkg, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatalf("read package.json: %v", err)
	}

	if string(first) != string(second) {
		t.Fatalf("tsconfig.json changed on the second patch:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	if string(firstPkg) != string(secondPkg) {
		t.Fatalf("package.json changed on the second patch:\nfirst:\n%s\nsecond:\n%s", firstPkg, secondPkg)
	}
}

func TestPatchTSBuildConfigErrors(t *testing.T) {
	t.Run("missing tsconfig", func(t *testing.T) {
		if err := patchTSBuildConfig(t.TempDir()); err == nil {
			t.Fatal("patchTSBuildConfig succeeded with no tsconfig.json present")
		}
	})

	t.Run("malformed tsconfig", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "tsconfig.json"), []byte("{not json"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		if err := patchTSBuildConfig(dir); err == nil {
			t.Fatal("patchTSBuildConfig succeeded on malformed JSON")
		}
	})

	t.Run("tsconfig without compilerOptions", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "tsconfig.json"), []byte(`{"exclude":[]}`), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		if err := patchTSBuildConfig(dir); err == nil {
			t.Fatal("patchTSBuildConfig succeeded without compilerOptions")
		}
	})
}

// The committed sdk/ts/tsconfig.json must already carry what the patch applies,
// so a fresh checkout builds under TS6 without anyone running `generate sdk`.
func TestCommittedSDKTSConfigIsPatched(t *testing.T) {
	tsconfig := readJSON(t, filepath.Join("..", "sdk", "ts", "tsconfig.json"))
	opts, ok := tsconfig["compilerOptions"].(map[string]any)
	if !ok {
		t.Fatal("sdk/ts/tsconfig.json has no compilerOptions")
	}
	if got := opts["ignoreDeprecations"]; got != "6.0" {
		t.Fatalf("sdk/ts/tsconfig.json ignoreDeprecations = %v, want \"6.0\"", got)
	}
	if got := opts["rootDir"]; got != "src" {
		t.Fatalf("sdk/ts/tsconfig.json rootDir = %v, want \"src\"", got)
	}

	pkg := readJSON(t, filepath.Join("..", "sdk", "ts", "package.json"))
	dev, ok := pkg["devDependencies"].(map[string]any)
	if !ok {
		t.Fatal("sdk/ts/package.json has no devDependencies")
	}
	if got := dev["typescript"]; got != "^6.0.0" {
		t.Fatalf("sdk/ts/package.json typescript = %v, want \"^6.0.0\"", got)
	}
}

// fixtureSpec exercises the shapes specImmutableAttrs has to tell apart:
// a create/update pair that narrows its body (widgets), a create-only
// collection (gadgets), a pair whose bodies match (crates), sub-resource
// operations that are not a create/update pair (widgets/{id}/parts and the
// non-collection POST under thing), and a $ref that has to be resolved.
const fixtureSpec = `{
  "paths": {
    "/widgets": {
      "post": {"tags": ["widgets"], "requestBody": {"content": {"application/json": {
        "schema": {"$ref": "#/components/schemas/WidgetCreate"}}}}},
      "get": {"tags": ["widgets"]}
    },
    "/widgets/{id}": {
      "patch": {"tags": ["widgets"], "requestBody": {"content": {"application/json": {
        "schema": {"properties": {"$schema": {}, "size": {}}}}}}},
      "delete": {"tags": ["widgets"]}
    },
    "/widgets/{id}/parts": {
      "post": {"tags": ["widgets"], "requestBody": {"content": {"application/json": {
        "schema": {"properties": {"part": {}}}}}}}
    },
    "/gadgets": {
      "post": {"tags": ["gadgets"], "requestBody": {"content": {"application/json": {
        "schema": {"properties": {"$schema": {}, "name": {}, "serial": {}}}}}}}
    },
    "/crates": {
      "post": {"tags": ["crates"], "requestBody": {"content": {"application/json": {
        "schema": {"properties": {"name": {}, "height": {}}}}}}}
    },
    "/crates/{id}": {
      "put": {"tags": ["crates"], "requestBody": {"content": {"application/json": {
        "schema": {"properties": {"name": {}, "height": {}}}}}}}
    },
    "/thing/login": {
      "post": {"tags": ["thing"], "requestBody": {"content": {"application/json": {
        "schema": {"properties": {"password": {}}}}}}}
    },
    "/multi-word": {
      "post": {"tags": ["multi-word"], "requestBody": {"content": {"application/json": {
        "schema": {"properties": {"fixed": {}}}}}}}
    }
  },
  "components": {"schemas": {"WidgetCreate": {"properties": {
    "$schema": {}, "name": {}, "colour": {}, "size": {}}}}}
}`

func TestSpecImmutableAttrsFixture(t *testing.T) {
	got, err := specImmutableAttrs([]byte(fixtureSpec))
	if err != nil {
		t.Fatalf("specImmutableAttrs: %v", err)
	}

	want := map[string][]string{
		// create has name/colour/size, update only size.
		"widgets": {"colour", "name"},
		// no update endpoint at all, so nothing can ever be changed.
		"gadgets": {"name", "serial"},
		// create and update accept the same body.
		"crates": {},
		// tag becomes the resource file stem.
		"multi_word": {"fixed"},
	}
	for resource, wantAttrs := range want {
		gotAttrs, ok := got[resource]
		if !ok {
			t.Errorf("%s: missing from derivation", resource)
			continue
		}
		if !slices.Equal(gotAttrs, wantAttrs) {
			t.Errorf("%s: got %v, want %v", resource, gotAttrs, wantAttrs)
		}
	}
	// /thing/login is a POST but not on a collection root, so "thing" is not a
	// resource with a create endpoint.
	if attrs, ok := got["thing"]; ok {
		t.Errorf("thing: derived %v from a non-collection POST, want no entry", attrs)
	}
	if len(got) != len(want) {
		t.Errorf("derived %d resources (%v), want %d", len(got), sortedKeys(got), len(want))
	}
}

func TestSpecImmutableAttrsRejectsDanglingRef(t *testing.T) {
	const spec = `{"paths": {"/widgets": {"post": {"tags": ["widgets"], "requestBody":
	  {"content": {"application/json": {"schema": {"$ref": "#/components/schemas/Nope"}}}}}}}}`
	if _, err := specImmutableAttrs([]byte(spec)); err == nil {
		t.Fatal("expected an error for an unresolvable $ref, got nil")
	}
}

// TestResourceImmutableAttrsMatchesSpec is the same tripwire
// patchResourceImmutability enforces at generate time, run against the
// committed spec so CI catches a narrowed update endpoint without anyone having
// to regenerate the provider.
func TestResourceImmutableAttrsMatchesSpec(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("..", "docs", "openapi.json"))
	if err != nil {
		t.Fatalf("read spec: %v", err)
	}
	derived, err := specImmutableAttrs(b)
	if err != nil {
		t.Fatalf("specImmutableAttrs: %v", err)
	}

	for _, resource := range sortedKeys(resourceImmutableAttrs) {
		if _, ok := derived[resource]; !ok {
			t.Errorf("%s: declared in resourceImmutableAttrs but not derivable from the spec", resource)
		}
	}

	for _, resource := range sortedKeys(derived) {
		if _, managed := resourceSchemaRoles[resource]; !managed {
			continue
		}
		// Only attributes the provider actually models can carry a plan
		// modifier; a write-only create field has nothing to attach one to.
		var modelled []string
		for _, attr := range derived[resource] {
			if _, ok := resourceSchemaRoles[resource][attr]; ok {
				modelled = append(modelled, attr)
			}
		}
		if modelled == nil {
			modelled = []string{}
		}
		if !slices.Equal(modelled, resourceImmutableAttrs[resource]) {
			t.Errorf("%s: spec says %v is immutable, resourceImmutableAttrs declares %v",
				resource, modelled, resourceImmutableAttrs[resource])
		}
	}
}

func TestAddProviderImports(t *testing.T) {
	const src = `package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)
`
	out, err := addProviderImports(src, []string{
		"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier",
		"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier",
	})
	if err != nil {
		t.Fatalf("addProviderImports: %v", err)
	}
	for _, want := range []string{"schema/planmodifier", "schema/stringplanmodifier"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing import %q in:\n%s", want, out)
		}
	}

	// Re-running must not duplicate what is already there.
	again, err := addProviderImports(out, []string{
		"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier",
	})
	if err != nil {
		t.Fatalf("addProviderImports (repeat): %v", err)
	}
	if again != out {
		t.Errorf("second call changed the file:\n%s", again)
	}

	if _, err := addProviderImports("package provider\n", []string{"x"}); err == nil {
		t.Fatal("expected an error when the schema import anchor is absent, got nil")
	}
}

func TestAssertProviderCheckout(t *testing.T) {
	t.Run("missing directory is allowed", func(t *testing.T) {
		if err := assertProviderCheckout(filepath.Join(t.TempDir(), "nope")); err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("empty directory is allowed", func(t *testing.T) {
		if err := assertProviderCheckout(t.TempDir()); err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("the provider checkout is allowed", func(t *testing.T) {
		dir := t.TempDir()
		mustWrite(t, filepath.Join(dir, "go.mod"), "module "+providerModulePath+"\n\ngo 1.25.8\n")
		if err := assertProviderCheckout(dir); err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("some other module is refused", func(t *testing.T) {
		dir := t.TempDir()
		mustWrite(t, filepath.Join(dir, "go.mod"), "module github.com/glueops/waggle\n")
		if err := assertProviderCheckout(dir); err == nil {
			t.Error("expected an error for a foreign module, got nil")
		}
	})

	t.Run("a non-empty directory with no go.mod is refused", func(t *testing.T) {
		dir := t.TempDir()
		mustWrite(t, filepath.Join(dir, "notes.txt"), "important\n")
		if err := assertProviderCheckout(dir); err == nil {
			t.Error("expected an error for a non-module directory, got nil")
		}
	})
}

func TestSyncProviderRepoPreservesRepoOwnedFiles(t *testing.T) {
	stage := t.TempDir()
	mustWrite(t, filepath.Join(stage, "internal", "client", "client.go"), "package client\n")
	mustWrite(t, filepath.Join(stage, "internal", "provider", "provider.go"), "package provider\n")
	mustWrite(t, filepath.Join(stage, ".openapi-generator", "FILES"), "main.go\n")
	mustWrite(t, filepath.Join(stage, "go.mod"), "module "+providerModulePath+"\n\ngo 1.24.0\n")
	mustWrite(t, filepath.Join(stage, "README.md"), "generated readme\n")

	target := t.TempDir()
	// Things the provider repo owns and the generator must not touch.
	mustWrite(t, filepath.Join(target, "go.mod"), "module "+providerModulePath+"\n\ngo 1.25.8\n")
	mustWrite(t, filepath.Join(target, "README.md"), "hand-written readme\n")
	mustWrite(t, filepath.Join(target, ".goreleaser.yml"), "builds: []\n")
	mustWrite(t, filepath.Join(target, ".changes", "0.1.20.md"), "notes\n")
	mustWrite(t, filepath.Join(target, "docs", "index.md"), "docs\n")
	mustWrite(t, filepath.Join(target, "examples", "resources", "waggle_pools", "resource.tf"), "curated\n")
	// A generated file that no longer exists upstream must be swept away.
	mustWrite(t, filepath.Join(target, "internal", "provider", "stale_resource.go"), "package provider\n")

	if err := syncProviderRepo(stage, target); err != nil {
		t.Fatalf("syncProviderRepo: %v", err)
	}

	for path, want := range map[string]string{
		"go.mod":             "module " + providerModulePath + "\n\ngo 1.25.8\n",
		"README.md":          "hand-written readme\n",
		".goreleaser.yml":    "builds: []\n",
		".changes/0.1.20.md": "notes\n",
		"docs/index.md":      "docs\n",
		"examples/resources/waggle_pools/resource.tf": "curated\n",
		"internal/provider/provider.go":               "package provider\n",
		"internal/client/client.go":                   "package client\n",
	} {
		b, err := os.ReadFile(filepath.Join(target, filepath.FromSlash(path)))
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		if string(b) != want {
			t.Errorf("%s: got %q, want %q", path, b, want)
		}
	}

	if _, err := os.Stat(filepath.Join(target, "internal", "provider", "stale_resource.go")); !os.IsNotExist(err) {
		t.Error("stale generated file survived the sync")
	}
}

func TestSyncProviderRepoWritesMissingScaffolding(t *testing.T) {
	stage := t.TempDir()
	mustWrite(t, filepath.Join(stage, "internal", "client", "client.go"), "package client\n")
	mustWrite(t, filepath.Join(stage, "internal", "provider", "provider.go"), "package provider\n")
	mustWrite(t, filepath.Join(stage, ".openapi-generator", "FILES"), "main.go\n")
	mustWrite(t, filepath.Join(stage, "go.mod"), "module "+providerModulePath+"\n")
	mustWrite(t, filepath.Join(stage, "main.go"), "package main\n")

	target := filepath.Join(t.TempDir(), "fresh-clone")
	if err := syncProviderRepo(stage, target); err != nil {
		t.Fatalf("syncProviderRepo: %v", err)
	}
	for _, name := range []string{"go.mod", "main.go", "internal/provider/provider.go"} {
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(name))); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

func TestSyncProviderRepoRejectsUnknownGeneratedPackage(t *testing.T) {
	stage := t.TempDir()
	mustWrite(t, filepath.Join(stage, "internal", "client", "client.go"), "package client\n")
	mustWrite(t, filepath.Join(stage, "internal", "provider", "provider.go"), "package provider\n")
	mustWrite(t, filepath.Join(stage, "internal", "surprise", "x.go"), "package surprise\n")

	err := syncProviderRepo(stage, t.TempDir())
	if err == nil {
		t.Fatal("expected an error for an unsynced generated package, got nil")
	}
	if !strings.Contains(err.Error(), "surprise") {
		t.Errorf("error does not name the package: %v", err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
