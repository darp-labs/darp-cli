// Package doctor validates the DARP contracts in a project without modifying it.
package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/darpbr/darp-cli/internal/project/asset"
	"github.com/darpbr/darp-cli/internal/project/discovery"
)

const supportedMajor = 1

const skillsRoot = ".agents/skills"

// State is the outcome of an individual diagnosis check.
type State string

const (
	Pass    State = "PASS"
	Warning State = "WARNING"
	Fail    State = "FAIL"
)

// CheckResult is the outcome reported by one independent contract check.
type CheckResult struct {
	Name    string
	State   State
	Message string
}

// Result contains all contract check results.
type Result struct {
	Checks []CheckResult
}

// ExitCode returns the process code required by the diagnosis result.
func (r Result) ExitCode() int {
	for _, check := range r.Checks {
		if check.State == Fail {
			return 1
		}
	}
	return 0
}

// Counts returns the numbers used by the doctor summary.
func (r Result) Counts() (passed, warnings, errors int) {
	for _, check := range r.Checks {
		switch check.State {
		case Pass:
			passed++
		case Warning:
			warnings++
		case Fail:
			errors++
		}
	}
	return
}

// Service diagnoses DARP project contracts rooted at a directory.
type Service struct{}

// NewService creates a doctor service.
func NewService() Service { return Service{} }

// Diagnose executes all required checks. It only reads project files.
func (Service) Diagnose(root string) Result {
	return Result{Checks: []CheckResult{
		configurationCheck(root),
		structureCheck(root),
		workflowsCheck(root),
		skillsCheck(root),
		templatesCheck(root),
		governanceCheck(root),
		versionCompatibilityCheck(root),
		assetsCheck(root),
	}}
}

type projectConfig struct {
	Version    string            `yaml:"version"`
	Project    projectSection    `yaml:"project"`
	Governance governanceSection `yaml:"governance"`
	Workflows  workflowsSection  `yaml:"workflows"`
	Skills     map[string]string `yaml:"skills"`
	Assets     []assetEntry      `yaml:"assets"`
}

type assetEntry struct {
	Path   string `yaml:"path"`
	Family string `yaml:"family"`
	Type   string `yaml:"type"`
}

type projectSection struct {
	Name string `yaml:"name"`
}
type governanceSection struct {
	Lifecycle string `yaml:"lifecycle"`
}
type workflowsSection struct {
	Default string `yaml:"default"`
}
type workflow struct {
	Name  string   `yaml:"name"`
	Steps []string `yaml:"steps"`
}

func configurationCheck(root string) CheckResult {
	config, err := readConfig(root)
	if err != nil {
		return failed("Configuration", err)
	}
	if strings.TrimSpace(config.Version) == "" {
		return failMessage("Configuration", "missing required field version")
	}
	if strings.TrimSpace(config.Project.Name) == "" {
		return failMessage("Configuration", "missing required field project.name")
	}
	if strings.TrimSpace(config.Governance.Lifecycle) == "" {
		return failMessage("Configuration", "missing required field governance.lifecycle")
	}
	if strings.TrimSpace(config.Workflows.Default) == "" {
		return failMessage("Configuration", "missing required field workflows.default")
	}
	if len(config.Skills) == 0 {
		return failMessage("Configuration", "missing required field skills")
	}
	if err := validateProjectPath(root, config.Governance.Lifecycle); err != nil {
		return failed("Configuration", fmt.Errorf("governance.lifecycle: %w", err))
	}
	for _, name := range sortedSkillNames(config.Skills) {
		path := config.Skills[name]
		if strings.TrimSpace(name) == "" || strings.TrimSpace(path) == "" {
			return failMessage("Configuration", "skills entries cannot be empty")
		}
		if filepath.Base(name) != name || name == "." || name == ".." || strings.ContainsAny(name, `/\\`) {
			return failMessage("Configuration", fmt.Sprintf("invalid skill name %q", name))
		}
		if err := validateProjectPath(root, path); err != nil {
			return failed("Configuration", fmt.Errorf("skills.%s: %w", name, err))
		}
		if filepath.Clean(path) != filepath.Join(skillsRoot, name) {
			return failMessage("Configuration", fmt.Sprintf("skills.%s must point to %s", name, filepath.Join(skillsRoot, name)))
		}
		info, err := os.Stat(filepath.Join(root, path))
		if err != nil || !info.IsDir() {
			return failMessage("Configuration", fmt.Sprintf("skills.%s must point to a directory", name))
		}
	}
	return pass("Configuration")
}

func structureCheck(root string) CheckResult {
	for _, path := range []string{".darp/governance", ".darp/workflows", skillsRoot, ".darp/templates"} {
		info, err := os.Stat(filepath.Join(root, path))
		if err != nil || !info.IsDir() {
			return failMessage("Structure", "missing "+path)
		}
	}
	return pass("Structure")
}

func workflowsCheck(root string) CheckResult {
	config, err := readConfig(root)
	if err != nil {
		return failed("Workflows", err)
	}
	directory := filepath.Join(root, ".darp", "workflows")
	entries, err := os.ReadDir(directory)
	if err != nil {
		return failed("Workflows", fmt.Errorf("read workflows: %w", err))
	}
	seen := make(map[string]bool)
	foundDefault := false
	for _, entry := range entries {
		if entry.IsDir() {
			return failMessage("Workflows", fmt.Sprintf("workflow entry %q must be a YAML file", entry.Name()))
		}
		if filepath.Ext(entry.Name()) != ".yaml" {
			return failMessage("Workflows", fmt.Sprintf("workflow entry %q must use the .yaml extension", entry.Name()))
		}
		workflowPath := filepath.Join(directory, entry.Name())
		if err := validateContainedPath(root, workflowPath); err != nil {
			return failed("Workflows", fmt.Errorf("workflow %q: %w", entry.Name(), err))
		}
		content, err := os.ReadFile(workflowPath)
		if err != nil {
			return failed("Workflows", err)
		}
		var value workflow
		if err := yaml.Unmarshal(content, &value); err != nil {
			return failed("Workflows", fmt.Errorf("invalid %s: %w", entry.Name(), err))
		}
		if strings.TrimSpace(value.Name) == "" {
			return failMessage("Workflows", fmt.Sprintf("workflow %q is missing name", entry.Name()))
		}
		if seen[value.Name] {
			return failMessage("Workflows", fmt.Sprintf("duplicate workflow name %q", value.Name))
		}
		seen[value.Name] = true
		if value.Name == config.Workflows.Default {
			foundDefault = true
		}
		if len(value.Steps) == 0 {
			return failMessage("Workflows", fmt.Sprintf("workflow %q has empty steps", value.Name))
		}
		for _, step := range value.Steps {
			if strings.TrimSpace(step) == "" {
				return failMessage("Workflows", fmt.Sprintf("workflow %q has an empty step", value.Name))
			}
			if filepath.Base(step) != step || step == "." || step == ".." || strings.ContainsAny(step, `/\\`) {
				return failMessage("Workflows", fmt.Sprintf("workflow %q has invalid skill name %q", value.Name, step))
			}
			skillPath := filepath.Join(root, skillsRoot, step)
			if err := validateContainedPath(root, skillPath); err != nil {
				return failed("Workflows", fmt.Errorf("workflow %q skill %q: %w", value.Name, step, err))
			}
			info, err := os.Stat(skillPath)
			if err != nil || !info.IsDir() {
				return failMessage("Workflows", fmt.Sprintf("workflow %q references missing skill %q", value.Name, step))
			}
		}
	}
	if !foundDefault {
		return failMessage("Workflows", fmt.Sprintf("workflow %q not found", config.Workflows.Default))
	}
	return pass("Workflows")
}

func skillsCheck(root string) CheckResult {
	directory := filepath.Join(root, skillsRoot)
	entries, err := os.ReadDir(directory)
	if err != nil {
		return failed("Skills", fmt.Errorf("read skills: %w", err))
	}
	for _, entry := range entries {
		skillPath := filepath.Join(directory, entry.Name())
		info, err := os.Stat(skillPath)
		if err != nil {
			return failed("Skills", fmt.Errorf("skill %q: %w", entry.Name(), err))
		}
		if !info.IsDir() {
			return failMessage("Skills", fmt.Sprintf("skill entry %q must be a directory", entry.Name()))
		}
		if err := validateContainedPath(root, skillPath); err != nil {
			return failed("Skills", fmt.Errorf("skill %q: %w", entry.Name(), err))
		}
		skillFile := filepath.Join(skillPath, "SKILL.md")
		if err := validateContainedPath(root, skillFile); err != nil {
			return failed("Skills", fmt.Errorf("skill %q SKILL.md: %w", entry.Name(), err))
		}
		info, err = os.Stat(skillFile)
		if err != nil || !info.Mode().IsRegular() {
			return failMessage("Skills", fmt.Sprintf("skill %q is missing SKILL.md", entry.Name()))
		}
		for _, optional := range []string{"prompts", "examples", "references", "scripts", "templates"} {
			optionalPath := filepath.Join(directory, entry.Name(), optional)
			_, err = os.Lstat(optionalPath)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return failed("Skills", fmt.Errorf("skill %q optional resource %s: %w", entry.Name(), optional, err))
			}
			info, err = os.Stat(optionalPath)
			if err != nil {
				return failed("Skills", fmt.Errorf("skill %q optional resource %s: %w", entry.Name(), optional, err))
			}
			if !info.IsDir() {
				return failMessage("Skills", fmt.Sprintf("skill %q optional resource %s must be a directory", entry.Name(), optional))
			}
			if err := validateContainedPath(root, optionalPath); err != nil {
				return failed("Skills", fmt.Errorf("skill %q %s: %w", entry.Name(), optional, err))
			}
		}
	}
	return pass("Skills")
}

func templatesCheck(root string) CheckResult {
	info, err := os.Stat(filepath.Join(root, ".darp", "templates"))
	if err != nil || !info.IsDir() {
		return failMessage("Templates", "missing .darp/templates")
	}
	return pass("Templates")
}

func governanceCheck(root string) CheckResult {
	for _, path := range []string{".darp/lifecycle.md", ".darp/governance/quality-gates.md"} {
		info, err := os.Stat(filepath.Join(root, path))
		if err != nil || info.IsDir() {
			return failMessage("Governance", "missing "+path)
		}
	}
	return pass("Governance")
}

func versionCompatibilityCheck(root string) CheckResult {
	config, err := readConfig(root)
	if err != nil {
		return failed("Version Compatibility", err)
	}
	major, err := contractMajor(config.Version)
	if err != nil {
		return failed("Version Compatibility", err)
	}
	switch {
	case major > supportedMajor:
		return failMessage("Version Compatibility", "unsupported contract version "+config.Version)
	case major < supportedMajor:
		return CheckResult{Name: "Version Compatibility", State: Warning, Message: "project contract version is older than CLI recommendation"}
	default:
		return pass("Version Compatibility")
	}
}

func assetsCheck(root string) CheckResult {
	config, err := readConfig(root)
	if err != nil {
		return failed("Assets", err)
	}
	var missing []string
	for _, entry := range config.Assets {
		if strings.TrimSpace(entry.Path) == "" {
			return failMessage("Assets", "asset path cannot be empty")
		}
		if err := validateRelativePath(entry.Path); err != nil {
			return failMessage("Assets", fmt.Sprintf("asset %q: %v", entry.Path, err))
		}
		if !asset.ValidFamily(entry.Family) {
			return failMessage("Assets", fmt.Sprintf("asset %q has unknown family %q", entry.Path, entry.Family))
		}
		if !asset.ValidType(entry.Type) {
			return failMessage("Assets", fmt.Sprintf("asset %q has unknown type %q", entry.Path, entry.Type))
		}
		classified := discovery.Classify(entry.Path)
		matched := false
		for _, candidate := range classified {
			if candidate.Family == asset.Family(entry.Family) && candidate.Type == asset.Type(entry.Type) {
				matched = true
				break
			}
		}
		if !matched {
			return failMessage("Assets", fmt.Sprintf("asset %q has incompatible family/type %s/%s", entry.Path, entry.Family, entry.Type))
		}
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(entry.Path)))
		if err != nil {
			if !os.IsNotExist(err) {
				return failed("Assets", fmt.Errorf("asset %q: %w", entry.Path, err))
			}
			missing = append(missing, fmt.Sprintf("asset %q is missing", entry.Path))
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return failMessage("Assets", fmt.Sprintf("asset %q must not be a symlink", entry.Path))
		}
		if info.IsDir() {
			return failMessage("Assets", fmt.Sprintf("asset %q must be a regular file", entry.Path))
		}
	}
	if len(missing) > 0 {
		return CheckResult{Name: "Assets", State: Warning, Message: strings.Join(missing, "; ")}
	}
	return pass("Assets")
}

func readConfig(root string) (projectConfig, error) {
	content, err := os.ReadFile(filepath.Join(root, "darp.yml"))
	if err != nil {
		return projectConfig{}, fmt.Errorf("read darp.yml: %w", err)
	}
	var config projectConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return projectConfig{}, fmt.Errorf("invalid darp.yml: %w", err)
	}
	return config, nil
}

func validateRelativePath(value string) error {
	path := strings.TrimSpace(value)
	if path == "" || filepath.IsAbs(path) {
		return fmt.Errorf("must be a non-empty relative path")
	}
	clean := filepath.Clean(path)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return fmt.Errorf("must remain inside the project")
	}
	return nil
}

func validateProjectPath(root, value string) error {
	path := strings.TrimSpace(value)
	if path == "" || filepath.IsAbs(path) {
		return fmt.Errorf("must be a non-empty relative path")
	}
	clean := filepath.Clean(path)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return fmt.Errorf("must remain inside the project")
	}
	target := filepath.Join(root, clean)
	if _, err := os.Stat(target); err != nil {
		return fmt.Errorf("path does not exist: %s", path)
	}
	return validateContainedPath(root, target)
}

func validateContainedPath(root, target string) error {
	rootReal, err := filepath.EvalSymlinks(root)
	if err != nil {
		return fmt.Errorf("resolve project root: %w", err)
	}
	targetReal, err := filepath.EvalSymlinks(target)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	relative, err := filepath.Rel(rootReal, targetReal)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path must remain inside the project")
	}
	return nil
}

func contractMajor(version string) (int, error) {
	parts := strings.Split(strings.TrimSpace(version), ".")
	if len(parts) < 2 || parts[0] == "" {
		return 0, fmt.Errorf("invalid contract version %q", version)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil || major < 0 {
		return 0, fmt.Errorf("invalid contract version %q", version)
	}
	return major, nil
}

func pass(name string) CheckResult { return CheckResult{Name: name, State: Pass} }
func failMessage(name, message string) CheckResult {
	return CheckResult{Name: name, State: Fail, Message: message}
}
func failed(name string, err error) CheckResult { return failMessage(name, err.Error()) }

func sortedSkillNames(skills map[string]string) []string {
	names := make([]string, 0, len(skills))
	for name := range skills {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
