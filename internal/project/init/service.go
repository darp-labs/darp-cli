package init

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const configFileName = "darp.yml"
const skillsRoot = ".agents/skills"

var darpDirectories = []string{
	".darp",
	".darp/governance",
	".darp/workflows",
	".darp/templates",
	".agents/skills/documentation",
	".agents/skills/architecture",
	".agents/skills/testing",
	".agents/skills/release",
}

var managedSkills = []string{"documentation", "architecture", "testing", "release"}

// FileSystem abstracts filesystem interactions used by the initialization service.
type FileSystem interface {
	Exists(path string) (bool, error)
	MkdirAll(path string) error
	WriteFile(path string, data []byte) error
	ReadFile(path string) ([]byte, error)
	Rename(oldPath, newPath string) error
	Base(path string) string
}

// Result captures the visible outcome of the initialization flow.
type Result struct {
	AlreadyInitialized bool
	Messages           []string
}

// Service initializes a directory as a DARP project.
type Service struct {
	fs FileSystem
}

// NewService creates a new initialization service.
func NewService(fs FileSystem, _ ...string) Service {
	return Service{fs: fs}
}

// Initialize applies the project bootstrap structure to the provided directory.
func (s Service) Initialize(root string) (Result, error) {
	messages := []string{"✔ Initializing project"}

	projectName := strings.TrimSpace(s.fs.Base(root))
	if projectName == "" || projectName == "." || projectName == string(filepath.Separator) {
		return Result{}, errors.New("could not derive project name from current directory")
	}

	changed := false
	configPath := filepath.Join(root, configFileName)
	created, err := s.ensureConfig(configPath, projectName)
	if err != nil {
		return Result{}, fmt.Errorf("create darp.yml: %w", err)
	}
	if created {
		messages = append(messages, "✔ Creating darp.yml")
		changed = true
	}

	createdStructure := false
	for _, directory := range darpDirectories {
		created, err := s.createDirectoryIfMissing(filepath.Join(root, directory))
		if err != nil {
			return Result{}, fmt.Errorf("create %s: %w", directory, err)
		}
		createdStructure = createdStructure || created
		changed = changed || created
	}
	if createdStructure {
		messages = append(messages, "✔ Creating .darp structure")
	}

	contracts := embeddedContracts()
	createdContracts := false
	for _, contract := range contracts {
		created, err := s.writeContract(filepath.Join(root, contract.path), contract.path, []byte(contract.content))
		if err != nil {
			return Result{}, fmt.Errorf("write %s: %w", contract.path, err)
		}
		createdContracts = createdContracts || created
		changed = changed || created
	}
	if createdContracts {
		messages = append(messages, "✔ Creating DARP contracts")
	}

	if !changed {
		messages = append(messages, "✔ Project already initialized")
		return Result{AlreadyInitialized: true, Messages: messages}, nil
	}

	messages = append(messages, "✔ Project initialized")

	return Result{
		Messages: messages,
	}, nil
}

func (s Service) ensureConfig(path, projectName string) (bool, error) {
	exists, err := s.fs.Exists(path)
	if err != nil || !exists {
		if err != nil {
			return false, err
		}
		return true, s.fs.WriteFile(path, []byte(renderConfig(projectName)))
	}

	content, err := s.fs.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("read darp.yml: %w", err)
	}
	updated, err := updateExistingConfig(content)
	if err != nil {
		return false, err
	}
	if bytes.Equal(content, updated) {
		return false, nil
	}

	temporary := filepath.Join(filepath.Dir(path), ".darp-init-darp.yml.tmp")
	temporaryExists, err := s.fs.Exists(temporary)
	if err != nil {
		return false, fmt.Errorf("check temporary darp.yml: %w", err)
	}
	if temporaryExists {
		return false, errors.New("temporary darp.yml already exists")
	}
	if err := s.fs.WriteFile(temporary, updated); err != nil {
		return false, fmt.Errorf("write temporary darp.yml: %w", err)
	}
	if err := s.fs.Rename(temporary, path); err != nil {
		return false, fmt.Errorf("replace darp.yml: %w", err)
	}
	return true, nil
}

func updateExistingConfig(content []byte) ([]byte, error) {
	var document yaml.Node
	if err := yaml.Unmarshal(content, &document); err != nil {
		return nil, fmt.Errorf("invalid darp.yml: %w", err)
	}
	if len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return nil, errors.New("cannot safely update darp.yml: top-level value must be a mapping")
	}
	root := document.Content[0]
	if err := rejectDuplicateKeys(root); err != nil {
		return nil, fmt.Errorf("cannot safely update darp.yml: %w", err)
	}
	if err := validateExistingConfig(root); err != nil {
		return nil, fmt.Errorf("invalid darp.yml: %w", err)
	}

	skills := mappingValue(root, "skills")
	if skills == nil || skills.Kind != yaml.MappingNode {
		return nil, errors.New("cannot safely update darp.yml: skills must be a mapping")
	}
	changed := false
	for _, name := range managedSkills {
		if mappingValue(skills, name) == nil {
			changed = true
			skills.Content = append(skills.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: name, Tag: "!!str"},
				&yaml.Node{Kind: yaml.ScalarNode, Value: path.Join(skillsRoot, name), Tag: "!!str"},
			)
		}
	}
	if !changed {
		return content, nil
	}

	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(&document); err != nil {
		return nil, fmt.Errorf("encode darp.yml: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("close darp.yml encoder: %w", err)
	}
	return output.Bytes(), nil
}

func validateExistingConfig(root *yaml.Node) error {
	for _, field := range []string{"version", "project", "governance", "workflows", "skills"} {
		if mappingValue(root, field) == nil {
			return fmt.Errorf("missing required field %s", field)
		}
	}
	for _, field := range []string{"version"} {
		value := mappingValue(root, field)
		if strings.TrimSpace(value.Value) == "" {
			return fmt.Errorf("missing required field %s", field)
		}
	}
	for _, section := range []string{"project", "governance", "workflows"} {
		value := mappingValue(root, section)
		if value.Kind != yaml.MappingNode {
			return fmt.Errorf("%s must be a mapping", section)
		}
	}
	skills := mappingValue(root, "skills")
	if skills.Kind != yaml.MappingNode {
		return errors.New("skills must be a mapping")
	}
	if len(skills.Content) == 0 {
		return errors.New("missing required field skills")
	}
	for _, field := range []struct{ section, name string }{
		{"project", "name"}, {"governance", "lifecycle"}, {"workflows", "default"},
	} {
		section := mappingValue(root, field.section)
		value := mappingValue(section, field.name)
		if value == nil || strings.TrimSpace(value.Value) == "" {
			return fmt.Errorf("missing required field %s.%s", field.section, field.name)
		}
	}
	return nil
}

func mappingValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			return mapping.Content[index+1]
		}
	}
	return nil
}

func rejectDuplicateKeys(node *yaml.Node) error {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.MappingNode {
		seen := make(map[string]struct{}, len(node.Content)/2)
		for index := 0; index+1 < len(node.Content); index += 2 {
			key := node.Content[index].Value
			if _, exists := seen[key]; exists {
				return fmt.Errorf("duplicate key %q", key)
			}
			seen[key] = struct{}{}
			if err := rejectDuplicateKeys(node.Content[index+1]); err != nil {
				return err
			}
		}
	} else {
		for _, child := range node.Content {
			if err := rejectDuplicateKeys(child); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s Service) createDirectoryIfMissing(path string) (bool, error) {
	exists, err := s.fs.Exists(path)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}
	return true, s.fs.MkdirAll(path)
}

func (s Service) writeContract(path, relativePath string, data []byte) (bool, error) {
	exists, err := s.fs.Exists(path)
	if err != nil || !exists {
		if err != nil {
			return false, err
		}
		return true, s.fs.WriteFile(path, data)
	}

	current, err := s.fs.ReadFile(path)
	if err != nil {
		return false, err
	}
	if placeholder, ok := historicalPlaceholders[relativePath]; ok && bytes.Equal(current, []byte(placeholder)) {
		return true, s.fs.WriteFile(path, data)
	}
	return false, nil
}

func embeddedContracts() []struct{ path, content string } {
	contracts := make([]struct{ path, content string }, 0, len(assetPaths))
	for _, asset := range assetPaths {
		contracts = append(contracts, struct{ path, content string }{asset.target, assetFiles[asset.target]})
	}
	return contracts
}

func renderConfig(projectName string) string {
	return fmt.Sprintf(`version: "1.0"

project:
  name: %q

governance:
  lifecycle: .darp/lifecycle.md

workflows:
  default: implement

skills:
  documentation: .agents/skills/documentation
  architecture: .agents/skills/architecture
  testing: .agents/skills/testing
  release: .agents/skills/release
`, projectName)
}
