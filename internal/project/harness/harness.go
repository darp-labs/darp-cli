// Package harness reads and validates the declarative Harness manifest.
//
// The package is intentionally limited to local parsing and structural
// validation. It never executes providers, tools, workflows or other
// references declared by a manifest.
package harness

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// ManifestPath is the canonical path of a Harness manifest in a project.
	ManifestPath = ".darp/harness.yaml"
	// APIVersionV1Alpha1 is the only Harness contract supported by this package.
	APIVersionV1Alpha1 = "darp.dev/v1alpha1"
	// KindHarness identifies a Harness manifest.
	KindHarness = "Harness"
)

var knownBlocks = map[string]bool{
	"assets":        true,
	"providers":     true,
	"tools":         true,
	"policy":        true,
	"workflow":      true,
	"qualityGates":  true,
	"evaluation":    true,
	"compatibility": true,
}

// Metadata identifies a Harness within a project.
type Metadata struct {
	Name    string
	Version string
}

// Manifest is the structural representation of a Harness manifest.
// Blocks and Extensions retain their YAML nodes so unknown and not-yet-defined
// contract data is not discarded during reading.
type Manifest struct {
	APIVersion  string
	Kind        string
	Metadata    Metadata
	Blocks      map[string]*yaml.Node
	Extensions  map[string]*yaml.Node
	MetadataExt map[string]*yaml.Node
}

// Diagnostic describes one deterministic structural validation problem.
type Diagnostic struct {
	Path    string
	Message string
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("%s: %s", d.Path, d.Message)
}

// Parse parses and validates manifest data without accessing the filesystem.
func Parse(data []byte) (Manifest, []Diagnostic) {
	manifest := Manifest{
		Blocks:      make(map[string]*yaml.Node),
		Extensions:  make(map[string]*yaml.Node),
		MetadataExt: make(map[string]*yaml.Node),
	}

	var document yaml.Node
	if err := yaml.Unmarshal(data, &document); err != nil {
		return manifest, []Diagnostic{{Path: "$", Message: "invalid YAML: " + err.Error()}}
	}
	if document.Kind != yaml.DocumentNode || len(document.Content) != 1 {
		return manifest, []Diagnostic{{Path: "$", Message: "document must contain one root mapping"}}
	}
	root := document.Content[0]
	if root.Kind != yaml.MappingNode {
		return manifest, []Diagnostic{{Path: "$", Message: "root must be a mapping"}}
	}

	diagnostics := make([]Diagnostic, 0)
	seen := make(map[string]bool)
	for i := 0; i < len(root.Content); i += 2 {
		key := root.Content[i]
		value := root.Content[i+1]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
			diagnostics = append(diagnostics, Diagnostic{Path: "$", Message: "root keys must be strings"})
			continue
		}
		name := key.Value
		if seen[name] {
			diagnostics = append(diagnostics, Diagnostic{Path: name, Message: "duplicate field"})
			continue
		}
		seen[name] = true

		switch name {
		case "apiVersion":
			manifest.APIVersion = scalarString(value, "apiVersion", &diagnostics)
		case "kind":
			manifest.Kind = scalarString(value, "kind", &diagnostics)
		case "metadata":
			parseMetadata(value, &manifest, &diagnostics)
		default:
			if knownBlocks[name] {
				if value.Kind != yaml.MappingNode {
					diagnostics = append(diagnostics, Diagnostic{Path: name, Message: "must be a mapping"})
				} else {
					manifest.Blocks[name] = value
					if name == "assets" {
						validateAssets(value, &diagnostics)
					}
				}
				continue
			}
			manifest.Extensions[name] = value
		}
	}

	if manifest.APIVersion != APIVersionV1Alpha1 {
		diagnostics = append(diagnostics, Diagnostic{Path: "apiVersion", Message: "must be " + APIVersionV1Alpha1})
	}
	if manifest.Kind != KindHarness {
		diagnostics = append(diagnostics, Diagnostic{Path: "kind", Message: "must be Harness"})
	}
	if strings.TrimSpace(manifest.Metadata.Name) == "" {
		diagnostics = append(diagnostics, Diagnostic{Path: "metadata.name", Message: "is required and must be a non-empty string"})
	}
	if strings.TrimSpace(manifest.Metadata.Version) == "" {
		diagnostics = append(diagnostics, Diagnostic{Path: "metadata.version", Message: "is required and must be a non-empty string"})
	}

	sortDiagnostics(diagnostics)
	return manifest, diagnostics
}

// Load reads the canonical Harness manifest below root and validates it. It
// does not write files, access the network or inspect referenced resources.
func Load(root string) (Manifest, []Diagnostic, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(ManifestPath)))
	if err != nil {
		return Manifest{}, nil, err
	}
	manifest, diagnostics := Parse(data)
	return manifest, diagnostics, nil
}

func parseMetadata(node *yaml.Node, manifest *Manifest, diagnostics *[]Diagnostic) {
	if node.Kind != yaml.MappingNode {
		*diagnostics = append(*diagnostics, Diagnostic{Path: "metadata", Message: "must be a mapping"})
		return
	}
	seen := make(map[string]bool)
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
			*diagnostics = append(*diagnostics, Diagnostic{Path: "metadata", Message: "keys must be strings"})
			continue
		}
		name := key.Value
		if seen[name] {
			*diagnostics = append(*diagnostics, Diagnostic{Path: "metadata." + name, Message: "duplicate field"})
			continue
		}
		seen[name] = true
		switch name {
		case "name":
			manifest.Metadata.Name = scalarString(value, "metadata.name", diagnostics)
		case "version":
			manifest.Metadata.Version = scalarString(value, "metadata.version", diagnostics)
		default:
			manifest.MetadataExt[name] = value
		}
	}
}

func scalarString(node *yaml.Node, field string, diagnostics *[]Diagnostic) string {
	if node.Kind != yaml.ScalarNode || node.Tag != "!!str" {
		*diagnostics = append(*diagnostics, Diagnostic{Path: field, Message: "must be a string"})
		return ""
	}
	return node.Value
}

func validateAssets(node *yaml.Node, diagnostics *[]Diagnostic) {
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" || strings.TrimSpace(key.Value) == "" {
			*diagnostics = append(*diagnostics, Diagnostic{Path: "assets", Message: "asset identities must be non-empty strings"})
			continue
		}
		field := "assets." + key.Value
		switch value.Kind {
		case yaml.ScalarNode:
			if value.Tag != "!!str" || strings.TrimSpace(value.Value) == "" {
				*diagnostics = append(*diagnostics, Diagnostic{Path: field, Message: "must reference a non-empty string path"})
				continue
			}
			validateAssetPath(field, value.Value, diagnostics)
		case yaml.MappingNode:
			validateAssetReference(field, value, diagnostics)
		default:
			*diagnostics = append(*diagnostics, Diagnostic{Path: field, Message: "must be a path string or reference mapping"})
		}
	}
}

func validateAssetReference(field string, node *yaml.Node, diagnostics *[]Diagnostic) {
	foundIdentity := false
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
			*diagnostics = append(*diagnostics, Diagnostic{Path: field, Message: "reference keys must be strings"})
			continue
		}
		switch key.Value {
		case "id":
			if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || strings.TrimSpace(value.Value) == "" {
				*diagnostics = append(*diagnostics, Diagnostic{Path: field + ".id", Message: "must be a non-empty string"})
			} else {
				foundIdentity = true
			}
		case "path":
			if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || strings.TrimSpace(value.Value) == "" {
				*diagnostics = append(*diagnostics, Diagnostic{Path: field + ".path", Message: "must be a non-empty string"})
			} else {
				foundIdentity = true
				validateAssetPath(field+".path", value.Value, diagnostics)
			}
		}
	}
	if !foundIdentity {
		*diagnostics = append(*diagnostics, Diagnostic{Path: field, Message: "must contain id or path"})
	}
}

func validateAssetPath(field, value string, diagnostics *[]Diagnostic) {
	if strings.Contains(value, "\\") || path.IsAbs(value) || filepath.IsAbs(filepath.FromSlash(value)) {
		*diagnostics = append(*diagnostics, Diagnostic{Path: field, Message: "path must be relative to the project"})
		return
	}
	clean := path.Clean(value)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		*diagnostics = append(*diagnostics, Diagnostic{Path: field, Message: "path must remain inside the project"})
	}
}

func sortDiagnostics(diagnostics []Diagnostic) {
	sort.SliceStable(diagnostics, func(i, j int) bool {
		if diagnostics[i].Path != diagnostics[j].Path {
			return diagnostics[i].Path < diagnostics[j].Path
		}
		return diagnostics[i].Message < diagnostics[j].Message
	})
}
