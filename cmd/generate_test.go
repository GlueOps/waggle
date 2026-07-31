package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
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
