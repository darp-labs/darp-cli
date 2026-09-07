package harness

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const minimalManifest = `apiVersion: darp.dev/v1alpha1
kind: Harness
metadata:
  name: example-harness
  version: 0.1.0
`

func TestParseMinimalManifest(t *testing.T) {
	manifest, diagnostics := Parse([]byte(minimalManifest))
	if len(diagnostics) != 0 {
		t.Fatalf("expected valid manifest, got %#v", diagnostics)
	}
	if manifest.APIVersion != APIVersionV1Alpha1 || manifest.Kind != KindHarness {
		t.Fatalf("unexpected identity: %#v", manifest)
	}
	if manifest.Metadata.Name != "example-harness" || manifest.Metadata.Version != "0.1.0" {
		t.Fatalf("unexpected metadata: %#v", manifest.Metadata)
	}
}

func TestParsePreservesBlocksAndUnknownFields(t *testing.T) {
	manifest, diagnostics := Parse([]byte(`apiVersion: darp.dev/v1alpha1
kind: Harness
metadata:
  name: demo
  version: 1.0.0
  futureMetadata: true
providers:
  primary: codex
  fallback: local
vendorExtension:
  enabled: true
`))
	if len(diagnostics) != 0 {
		t.Fatalf("expected valid manifest, got %#v", diagnostics)
	}
	if _, ok := manifest.Blocks["providers"]; !ok {
		t.Fatal("expected providers block to be preserved")
	}
	if _, ok := manifest.Extensions["vendorExtension"]; !ok {
		t.Fatal("expected unknown root field to be preserved")
	}
	if _, ok := manifest.MetadataExt["futureMetadata"]; !ok {
		t.Fatal("expected unknown metadata field to be preserved")
	}
}

func TestParseValidatesKnownFieldsDeterministically(t *testing.T) {
	data := []byte(`apiVersion: wrong
kind: Wrong
metadata:
  name: 42
providers: []
`)
	_, first := Parse(data)
	_, second := Parse(data)
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("diagnostics are not deterministic\nfirst: %#v\nsecond: %#v", first, second)
	}
	joined := diagnosticText(first)
	for _, want := range []string{
		"apiVersion: must be darp.dev/v1alpha1",
		"kind: must be Harness",
		"metadata.name: must be a string",
		"metadata.version: is required",
		"providers: must be a mapping",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing diagnostic %q in %q", want, joined)
		}
	}
}

func TestParseValidatesAssetReferences(t *testing.T) {
	valid := `apiVersion: darp.dev/v1alpha1
kind: Harness
metadata: {name: demo, version: 1.0.0}
assets:
  review: .github/instructions/review.instructions.md
  notInstalledYet: .github/skills/future/SKILL.md
  skill:
    id: review-skill
  pathRef:
    path: .codex/skills/review/SKILL.md
`
	if _, diagnostics := Parse([]byte(valid)); len(diagnostics) != 0 {
		t.Fatalf("expected valid asset references, got %#v", diagnostics)
	}

	invalid := `apiVersion: darp.dev/v1alpha1
kind: Harness
metadata: {name: demo, version: 1.0.0}
assets:
  absolute: /tmp/asset.md
  traversal: ../outside.md
  incomplete: {}
  nested: [not-a-reference]
`
	_, diagnostics := Parse([]byte(invalid))
	joined := diagnosticText(diagnostics)
	for _, want := range []string{
		"assets.absolute: path must be relative to the project",
		"assets.traversal: path must remain inside the project",
		"assets.incomplete: must contain id or path",
		"assets.nested: must be a path string or reference mapping",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing diagnostic %q in %q", want, joined)
		}
	}
}

func TestParseRejectsInvalidYAMLAndRoot(t *testing.T) {
	for name, data := range map[string]string{
		"invalid YAML":  "apiVersion: [",
		"scalar root":   "just a string",
		"sequence root": "- item\n",
	} {
		t.Run(name, func(t *testing.T) {
			_, diagnostics := Parse([]byte(data))
			if len(diagnostics) == 0 {
				t.Fatal("expected diagnostic")
			}
		})
	}
}

func TestLoadUsesCanonicalPathAndDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	manifestPath := filepath.Join(root, filepath.FromSlash(ManifestPath))
	if err := os.MkdirAll(filepath.Dir(manifestPath), 0o755); err != nil {
		t.Fatalf("create manifest directory: %v", err)
	}
	if err := os.WriteFile(manifestPath, []byte(minimalManifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	before, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read before: %v", err)
	}
	manifest, diagnostics, err := Load(root)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	if len(diagnostics) != 0 || manifest.Metadata.Name != "example-harness" {
		t.Fatalf("unexpected load result: %#v, %#v", manifest, diagnostics)
	}
	after, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read after: %v", err)
	}
	if !reflect.DeepEqual(before, after) {
		t.Fatal("load modified the manifest")
	}
}

func diagnosticText(diagnostics []Diagnostic) string {
	parts := make([]string, 0, len(diagnostics))
	for _, diagnostic := range diagnostics {
		parts = append(parts, diagnostic.String())
	}
	return strings.Join(parts, "\n")
}
