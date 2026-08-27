package init

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/darpbr/darp-cli/internal/project/asset"
	"github.com/darpbr/darp-cli/internal/project/discovery"
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
	Remove(path string) error
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

	configPath := filepath.Join(root, configFileName)
	report := discovery.Discover(root)
	configExists, err := s.fs.Exists(configPath)
	if err != nil {
		return Result{}, fmt.Errorf("check darp.yml: %w", err)
	}
	originalConfig := []byte(nil)
	configContent := []byte(renderConfig(projectName))
	if configExists {
		originalConfig, err = s.fs.ReadFile(configPath)
		if err != nil {
			return Result{}, fmt.Errorf("read darp.yml: %w", err)
		}
		configContent, err = updateExistingConfig(originalConfig)
		if err != nil {
			return Result{}, err
		}
	}
	updated, registrations, err := registerAssets(configContent, report.Found)
	if err != nil {
		return Result{}, err
	}

	transaction := newTransaction(s.fs)
	configChanged := !configExists || !bytes.Equal(originalConfig, updated)
	if configChanged {
		if configExists {
			transaction.rememberFile(configPath, originalConfig)
			if err := s.replaceConfig(configPath, updated); err != nil {
				return Result{}, transaction.fail(err)
			}
		} else {
			transaction.rememberMissing(configPath)
			if err := s.fs.WriteFile(configPath, updated); err != nil {
				return Result{}, transaction.fail(err)
			}
		}
	}

	createdStructure := false
	createdContracts := false
	changed := configChanged
	for _, directory := range darpDirectories {
		created, err := s.createDirectoryIfMissing(filepath.Join(root, directory))
		if err != nil {
			return Result{}, transaction.fail(fmt.Errorf("create %s: %w", directory, err))
		}
		if created {
			transaction.rememberMissing(filepath.Join(root, directory))
		}
		createdStructure = createdStructure || created
		changed = changed || created
	}
	for _, contract := range embeddedContracts() {
		contractPath := filepath.Join(root, contract.path)
		exists, err := s.fs.Exists(contractPath)
		if err != nil {
			return Result{}, transaction.fail(fmt.Errorf("check %s: %w", contract.path, err))
		}
		if exists {
			current, err := s.fs.ReadFile(contractPath)
			if err != nil {
				return Result{}, transaction.fail(fmt.Errorf("read %s: %w", contract.path, err))
			}
			if placeholder, ok := historicalPlaceholders[contract.path]; !ok || !bytes.Equal(current, []byte(placeholder)) {
				continue
			}
			transaction.rememberFile(contractPath, current)
		} else {
			transaction.rememberMissing(contractPath)
		}
		if err := s.fs.WriteFile(contractPath, []byte(contract.content)); err != nil {
			return Result{}, transaction.fail(fmt.Errorf("write %s: %w", contract.path, err))
		}
		createdContracts = true
		changed = true
	}
	if createdStructure {
		messages = append(messages, "✔ Creating .darp structure")
	}
	if createdContracts {
		messages = append(messages, "✔ Creating DARP contracts")
	}
	messages = append(messages, renderAssetMessages(report, registrations)...)

	if !changed {
		messages = append(messages, "✔ Project already initialized")
		return Result{AlreadyInitialized: true, Messages: messages}, nil
	}

	messages = append(messages, "✔ Project initialized")

	return Result{
		Messages: messages,
	}, nil
}

type transaction struct {
	fs    FileSystem
	undos []func() error
}

func newTransaction(fs FileSystem) *transaction { return &transaction{fs: fs} }

func (t *transaction) rememberFile(path string, content []byte) {
	t.undos = append(t.undos, func() error { return t.fs.WriteFile(path, content) })
}

func (t *transaction) rememberMissing(path string) {
	t.undos = append(t.undos, func() error {
		exists, err := t.fs.Exists(path)
		if err != nil || !exists {
			return err
		}
		return t.fs.Remove(path)
	})
}

func (t *transaction) fail(err error) error {
	rollbackErrs := make([]error, 0)
	for index := len(t.undos) - 1; index >= 0; index-- {
		if rollbackErr := t.undos[index](); rollbackErr != nil {
			rollbackErrs = append(rollbackErrs, rollbackErr)
		}
	}
	if len(rollbackErrs) > 0 {
		return errors.Join(err, fmt.Errorf("rollback failed: %w", errors.Join(rollbackErrs...)))
	}
	return err
}

func (s Service) replaceConfig(path string, content []byte) error {
	temporary := filepath.Join(filepath.Dir(path), ".darp-init-darp.yml.tmp")
	temporaryExists, err := s.fs.Exists(temporary)
	if err != nil {
		return fmt.Errorf("check temporary darp.yml: %w", err)
	}
	if temporaryExists {
		return errors.New("temporary darp.yml already exists")
	}
	if err := s.fs.WriteFile(temporary, content); err != nil {
		cleanupErr := s.fs.Remove(temporary)
		if cleanupErr != nil {
			return fmt.Errorf("write temporary darp.yml: %w; cleanup temporary darp.yml: %v", err, cleanupErr)
		}
		return fmt.Errorf("write temporary darp.yml: %w", err)
	}
	if err := s.fs.Rename(temporary, path); err != nil {
		cleanupErr := s.fs.Remove(temporary)
		if cleanupErr != nil {
			return fmt.Errorf("replace darp.yml: %w; cleanup temporary darp.yml: %v", err, cleanupErr)
		}
		return fmt.Errorf("replace darp.yml: %w", err)
	}
	return nil
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

type registrationStatus string

const (
	assetAdded             registrationStatus = "added"
	assetAlreadyRegistered registrationStatus = "already_registered"
)

type assetRegistration struct {
	asset  asset.Asset
	status registrationStatus
}

type assetIdentity struct {
	path   string
	family string
	typ    string
}

// registerAssets merges discovered assets into an existing darp.yml content.
// It is additive and idempotent: entries already present are reported as
// already registered and never duplicated. It never mutates the original
// content unless at least one new entry must be added.
func registerAssets(content []byte, found []asset.Asset) ([]byte, []assetRegistration, error) {
	var document yaml.Node
	if err := yaml.Unmarshal(content, &document); err != nil {
		return nil, nil, fmt.Errorf("invalid darp.yml: %w", err)
	}
	if len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return nil, nil, errors.New("cannot safely update darp.yml: top-level value must be a mapping")
	}
	root := document.Content[0]
	if err := rejectDuplicateKeys(root); err != nil {
		return nil, nil, fmt.Errorf("cannot safely update darp.yml: %w", err)
	}

	existing := make(map[assetIdentity]bool)
	assetsNode := mappingValue(root, "assets")
	if assetsNode != nil {
		if assetsNode.Kind == yaml.ScalarNode && assetsNode.Tag == "!!null" {
			assetsNode.Kind = yaml.SequenceNode
			assetsNode.Tag = "!!seq"
			assetsNode.Value = ""
			assetsNode.Content = nil
		}
		if assetsNode.Kind != yaml.SequenceNode {
			return nil, nil, errors.New("cannot safely update darp.yml: assets must be a sequence")
		}
		for _, entry := range assetsNode.Content {
			id, err := parseAssetIdentity(entry)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot safely update darp.yml: %w", err)
			}
			existing[id] = true
		}
	}

	registrations := make([]assetRegistration, 0, len(found))
	for _, a := range found {
		id := assetIdentity{path: a.Path, family: string(a.Family), typ: string(a.Type)}
		if existing[id] {
			registrations = append(registrations, assetRegistration{asset: a, status: assetAlreadyRegistered})
			continue
		}
		if assetsNode == nil {
			assetsNode = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			root.Content = append(root.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "assets", Tag: "!!str"},
				assetsNode,
			)
		}
		assetsNode.Content = append(assetsNode.Content, newAssetNode(a))
		existing[id] = true
		registrations = append(registrations, assetRegistration{asset: a, status: assetAdded})
	}

	if !hasAdded(registrations) {
		return content, registrations, nil
	}

	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(&document); err != nil {
		return nil, nil, fmt.Errorf("encode darp.yml: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, nil, fmt.Errorf("close darp.yml encoder: %w", err)
	}
	return output.Bytes(), registrations, nil
}

func parseAssetIdentity(entry *yaml.Node) (assetIdentity, error) {
	if entry.Kind != yaml.MappingNode {
		return assetIdentity{}, errors.New("each asset entry must be a mapping")
	}
	var id assetIdentity
	for i := 0; i+1 < len(entry.Content); i += 2 {
		key := entry.Content[i].Value
		value := entry.Content[i+1]
		switch key {
		case "path":
			id.path = scalarValue(value)
		case "family":
			id.family = scalarValue(value)
		case "type":
			id.typ = scalarValue(value)
		}
	}
	if strings.TrimSpace(id.path) == "" || strings.TrimSpace(id.family) == "" || strings.TrimSpace(id.typ) == "" {
		return assetIdentity{}, errors.New("asset entry must define path, family and type")
	}
	return id, nil
}

func scalarValue(node *yaml.Node) string {
	if node == nil || node.Kind != yaml.ScalarNode {
		return ""
	}
	return node.Value
}

func newAssetNode(a asset.Asset) *yaml.Node {
	return &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "path", Tag: "!!str"},
			{Kind: yaml.ScalarNode, Value: a.Path, Tag: "!!str"},
			{Kind: yaml.ScalarNode, Value: "family", Tag: "!!str"},
			{Kind: yaml.ScalarNode, Value: string(a.Family), Tag: "!!str"},
			{Kind: yaml.ScalarNode, Value: "type", Tag: "!!str"},
			{Kind: yaml.ScalarNode, Value: string(a.Type), Tag: "!!str"},
		},
	}
}

func hasAdded(registrations []assetRegistration) bool {
	for _, registration := range registrations {
		if registration.status == assetAdded {
			return true
		}
	}
	return false
}

func renderAssetMessages(report discovery.Report, registrations []assetRegistration) []string {
	messages := make([]string, 0)
	for _, registration := range registrations {
		a := registration.asset
		switch registration.status {
		case assetAdded:
			messages = append(messages, fmt.Sprintf("✔ Registered asset %s (%s/%s)", a.Path, a.Family, a.Type))
		case assetAlreadyRegistered:
			messages = append(messages, fmt.Sprintf("• Asset already registered %s (%s/%s)", a.Path, a.Family, a.Type))
		}
	}
	for _, ambiguity := range report.Ambiguous {
		messages = append(messages, fmt.Sprintf("⚠ Skipped ambiguous asset %s (%s)", ambiguity.Path, describeAmbiguity(ambiguity.Assets)))
	}
	for _, conflict := range report.SemanticConflicts {
		messages = append(messages, fmt.Sprintf("⚠ Skipped conflicting asset name %s (%s)", conflict.Name, describeAmbiguity(conflict.Assets)))
	}
	for _, ignored := range report.Ignored {
		messages = append(messages, fmt.Sprintf("• Ignored symlink %s", ignored.Path))
	}
	if report.Empty() {
		messages = append(messages, "• No supported assets found")
	}
	return messages
}

func describeAmbiguity(assets []asset.Asset) string {
	parts := make([]string, 0, len(assets))
	for _, a := range assets {
		parts = append(parts, fmt.Sprintf("%s/%s", a.Family, a.Type))
	}
	return strings.Join(parts, ", ")
}
