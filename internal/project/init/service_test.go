package init

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/darpbr/darp-cli/internal/project/doctor"
)

func TestInitializeCreatesExpectedProjectStructure(t *testing.T) {
	root := t.TempDir()
	projectDir := filepath.Join(root, "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}

	service := NewService(NewOSFileSystem(), "# ignored legacy content\n")

	result, err := service.Initialize(projectDir)
	if err != nil {
		t.Fatalf("initialize: %v", err)
	}

	if result.AlreadyInitialized {
		t.Fatalf("expected fresh initialization")
	}

	assertFileContent(t, filepath.Join(projectDir, "darp.yml"), `version: "1.0"

project:
  name: "demo-project"

governance:
  lifecycle: .darp/lifecycle.md

workflows:
  default: implement

skills:
  documentation: .agents/skills/documentation
  architecture: .agents/skills/architecture
  testing: .agents/skills/testing
  release: .agents/skills/release
`)
	assertFileContent(t, filepath.Join(projectDir, ".darp", "lifecycle.md"), assetFiles[".darp/lifecycle.md"])

	for _, directory := range darpDirectories {
		assertDirectoryExists(t, filepath.Join(projectDir, directory))
	}
	for _, skill := range managedSkills {
		assertFileContent(t, filepath.Join(projectDir, skillsRoot, skill, "SKILL.md"), assetFiles[filepath.Join(skillsRoot, skill, "SKILL.md")])
	}
}

func TestInitializeIsIdempotentAndDoesNotOverwriteFiles(t *testing.T) {
	root := t.TempDir()
	projectDir := filepath.Join(root, "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}

	service := NewService(NewOSFileSystem(), "# first version\n")

	if _, err := service.Initialize(projectDir); err != nil {
		t.Fatalf("first initialize: %v", err)
	}

	if err := os.WriteFile(filepath.Join(projectDir, "darp.yml"), []byte(`version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
  custom: .agents/skills/custom
metadata:
  owner: team
`), 0o644); err != nil {
		t.Fatalf("write custom darp.yml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, ".darp", "lifecycle.md"), []byte("# user changes\n"), 0o644); err != nil {
		t.Fatalf("write custom lifecycle.md: %v", err)
	}

	secondService := NewService(NewOSFileSystem(), "# second version\n")
	result, err := secondService.Initialize(projectDir)
	if err != nil {
		t.Fatalf("second initialize: %v", err)
	}

	if result.AlreadyInitialized {
		t.Fatalf("expected repair result after adding missing skills")
	}

	updated, err := os.ReadFile(filepath.Join(projectDir, "darp.yml"))
	if err != nil {
		t.Fatalf("read updated darp.yml: %v", err)
	}
	for _, skill := range managedSkills {
		if !strings.Contains(string(updated), skill+": .agents/skills/"+skill) {
			t.Fatalf("expected managed skill %q in updated config: %s", skill, updated)
		}
	}
	if !strings.Contains(string(updated), "custom: .agents/skills/custom") || !strings.Contains(string(updated), "owner: team") {
		t.Fatalf("expected custom configuration to be preserved: %s", updated)
	}
	assertFileContent(t, filepath.Join(projectDir, ".darp", "lifecycle.md"), "# user changes\n")
}

func TestInitializeRejectsInvalidExistingConfigWithoutCreatingAssets(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	original := []byte("name: custom\n")
	if err := os.WriteFile(filepath.Join(projectDir, "darp.yml"), original, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := NewService(NewOSFileSystem()).Initialize(projectDir)
	if err == nil || !strings.Contains(err.Error(), "invalid darp.yml") {
		t.Fatalf("expected invalid config error, got %v", err)
	}
	assertFileContent(t, filepath.Join(projectDir, "darp.yml"), string(original))
	if _, statErr := os.Stat(filepath.Join(projectDir, ".darp")); !os.IsNotExist(statErr) {
		t.Fatalf("expected no assets after invalid config, stat error: %v", statErr)
	}
}

func TestInitializeUpgradesOnlyExactHistoricalPlaceholders(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	if _, err := NewService(NewOSFileSystem()).Initialize(projectDir); err != nil {
		t.Fatalf("initial initialize: %v", err)
	}
	quality := filepath.Join(projectDir, ".darp/governance/quality-gates.md")
	documentation := filepath.Join(projectDir, ".agents/skills/documentation/SKILL.md")
	workflow := filepath.Join(projectDir, ".darp/workflows/implement.yaml")
	for path, content := range map[string]string{
		quality:       historicalPlaceholders[".darp/governance/quality-gates.md"],
		documentation: historicalPlaceholders[".agents/skills/documentation/SKILL.md"],
		workflow:      historicalPlaceholders[".darp/workflows/implement.yaml"],
	} {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write placeholder %s: %v", path, err)
		}
	}
	if _, err := NewService(NewOSFileSystem()).Initialize(projectDir); err != nil {
		t.Fatalf("placeholder upgrade: %v", err)
	}
	assertFileContent(t, quality, assetFiles[".darp/governance/quality-gates.md"])
	assertFileContent(t, documentation, assetFiles[".agents/skills/documentation/SKILL.md"])
	assertFileContent(t, workflow, assetFiles[".darp/workflows/implement.yaml"])
}

func TestInitializePreservesNearPlaceholdersAndIsIdempotent(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	if _, err := NewService(NewOSFileSystem()).Initialize(projectDir); err != nil {
		t.Fatalf("initial initialize: %v", err)
	}
	quality := filepath.Join(projectDir, ".darp/governance/quality-gates.md")
	original := "# Quality Gates\n\ncustom\n"
	if err := os.WriteFile(quality, []byte(original), 0o644); err != nil {
		t.Fatalf("write custom quality gates: %v", err)
	}
	if _, err := NewService(NewOSFileSystem()).Initialize(projectDir); err != nil {
		t.Fatalf("repair initialize: %v", err)
	}
	assertFileContent(t, quality, original)
	before := snapshotFiles(t, projectDir)
	if _, err := NewService(NewOSFileSystem()).Initialize(projectDir); err != nil {
		t.Fatalf("idempotent initialize: %v", err)
	}
	after := snapshotFiles(t, projectDir)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("second initialize changed files\nbefore: %#v\nafter: %#v", before, after)
	}
}

func TestInitializeRejectsDuplicateSkills(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	original := `version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
  documentation: .agents/skills/other
`
	if err := os.WriteFile(filepath.Join(projectDir, "darp.yml"), []byte(original), 0o644); err != nil {
		t.Fatalf("write duplicate config: %v", err)
	}
	_, err := NewService(NewOSFileSystem()).Initialize(projectDir)
	if err == nil || !strings.Contains(err.Error(), "duplicate key") {
		t.Fatalf("expected duplicate key error, got %v", err)
	}
	assertFileContent(t, filepath.Join(projectDir, "darp.yml"), original)
}

func TestCanonicalAssetsMatchRepositoryContracts(t *testing.T) {
	for _, asset := range assetPaths {
		t.Run(asset.target, func(t *testing.T) {
			rootPath := filepath.Join("..", "..", "..", asset.target)
			content, err := os.ReadFile(rootPath)
			if err != nil {
				t.Fatalf("read repository contract: %v", err)
			}
			if string(content) != assetFiles[asset.target] {
				t.Fatalf("embedded asset is not synchronized with %s", rootPath)
			}
		})
	}
}

func TestInitializePreservesConfigWhenAtomicReplacementFails(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	original := `version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
`
	configPath := filepath.Join(projectDir, "darp.yml")
	if err := os.WriteFile(configPath, []byte(original), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	fs := failingRenameFileSystem{FileSystem: NewOSFileSystem()}
	_, err := NewService(fs).Initialize(projectDir)
	if err == nil || !strings.Contains(err.Error(), "replace darp.yml") {
		t.Fatalf("expected replacement error, got %v", err)
	}
	assertFileContent(t, configPath, original)
	if _, statErr := os.Stat(filepath.Join(projectDir, ".darp-init-darp.yml.tmp")); !os.IsNotExist(statErr) {
		t.Fatalf("expected temporary file cleanup, stat error: %v", statErr)
	}
	if _, statErr := os.Stat(filepath.Join(projectDir, ".darp")); !os.IsNotExist(statErr) {
		t.Fatalf("expected no assets after replacement failure, stat error: %v", statErr)
	}
}

func TestInitializeCleansTemporaryAfterWriteFailure(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	original := []byte(`version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
`)
	configPath := filepath.Join(projectDir, "darp.yml")
	if err := os.WriteFile(configPath, original, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := NewService(failingWriteFileSystem{FileSystem: NewOSFileSystem()}).Initialize(projectDir)
	if err == nil || !strings.Contains(err.Error(), "write temporary darp.yml") {
		t.Fatalf("expected temporary write error, got %v", err)
	}
	assertFileContent(t, configPath, string(original))
	if _, statErr := os.Stat(filepath.Join(projectDir, ".darp-init-darp.yml.tmp")); !os.IsNotExist(statErr) {
		t.Fatalf("expected temporary file cleanup, stat error: %v", statErr)
	}
}

func TestInitializeRollsBackAfterLateContractFailure(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	before := snapshotFiles(t, projectDir)

	fs := failingContractWriteFileSystem{FileSystem: NewOSFileSystem()}
	if _, err := NewService(fs).Initialize(projectDir); err == nil {
		t.Fatal("expected contract write failure")
	}
	after := snapshotFiles(t, projectDir)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("expected rollback to preserve project\nbefore: %q\nafter: %q", before, after)
	}
}

func TestInitializeRejectsInvalidAssetsBeforeCreatingStructure(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	config := `version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
assets: invalid
`
	configPath := filepath.Join(projectDir, "darp.yml")
	if err := os.WriteFile(configPath, []byte(config), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := NewService(NewOSFileSystem()).Initialize(projectDir); err == nil || !strings.Contains(err.Error(), "assets must be a sequence") {
		t.Fatalf("expected invalid assets error, got %v", err)
	}
	assertFileContent(t, configPath, config)
	if _, statErr := os.Stat(filepath.Join(projectDir, ".darp")); !os.IsNotExist(statErr) {
		t.Fatalf("expected no structure after preflight failure, stat error: %v", statErr)
	}
}

func TestInitializeIsAgnosticToProjectStack(t *testing.T) {
	testCases := map[string][]string{
		"empty":          nil,
		"java-spring":    {"pom.xml", "src/main/java/Application.java"},
		"python-fastapi": {"pyproject.toml", "app/main.py"},
		"go":             {"go.mod", "cmd/server/main.go"},
		"monorepo":       {"frontend/package.json", "services/api/go.mod", "tools/script.py"},
	}
	for name, files := range testCases {
		t.Run(name, func(t *testing.T) {
			projectDir := filepath.Join(t.TempDir(), "demo-project")
			if err := os.Mkdir(projectDir, 0o755); err != nil {
				t.Fatalf("mkdir project dir: %v", err)
			}
			for _, file := range files {
				filePath := filepath.Join(projectDir, filepath.FromSlash(file))
				if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
					t.Fatalf("mkdir fixture path: %v", err)
				}
				if err := os.WriteFile(filePath, []byte("fixture\n"), 0o644); err != nil {
					t.Fatalf("write fixture file: %v", err)
				}
			}
			if _, err := NewService(NewOSFileSystem()).Initialize(projectDir); err != nil {
				t.Fatalf("initialize: %v", err)
			}
			if result := doctor.NewService().Diagnose(projectDir); result.ExitCode() != 0 {
				t.Fatalf("doctor rejected fixture: %#v", result.Checks)
			}
		})
	}
}

type failingRenameFileSystem struct {
	FileSystem
}

func (failingRenameFileSystem) Rename(_, _ string) error {
	return errors.New("simulated rename failure")
}

type failingWriteFileSystem struct {
	FileSystem
}

func (failingWriteFileSystem) WriteFile(path string, data []byte) error {
	if filepath.Base(path) == ".darp-init-darp.yml.tmp" {
		return errors.New("simulated temporary write failure")
	}
	return NewOSFileSystem().WriteFile(path, data)
}

type failingContractWriteFileSystem struct {
	FileSystem
}

func (failingContractWriteFileSystem) WriteFile(path string, data []byte) error {
	if filepath.Base(path) == "quality-gates.md" {
		return errors.New("simulated contract write failure")
	}
	return NewOSFileSystem().WriteFile(path, data)
}

func snapshotFiles(t *testing.T, root string) map[string]string {
	t.Helper()
	result := make(map[string]string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		result[relative] = string(content)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot files: %v", err)
	}
	return result
}

func TestInitializeRestoresMissingConfigWithoutChangingExistingDarpFiles(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}

	service := NewService(NewOSFileSystem(), "# lifecycle\n")
	if _, err := service.Initialize(projectDir); err != nil {
		t.Fatalf("first initialize: %v", err)
	}

	importantSpec := filepath.Join(projectDir, ".darp", "specs", "important.md")
	if err := os.MkdirAll(filepath.Dir(importantSpec), 0o755); err != nil {
		t.Fatalf("mkdir specs: %v", err)
	}
	if err := os.WriteFile(importantSpec, []byte("# Keep this\n"), 0o644); err != nil {
		t.Fatalf("write important spec: %v", err)
	}
	if err := os.Remove(filepath.Join(projectDir, "darp.yml")); err != nil {
		t.Fatalf("remove config: %v", err)
	}

	result, err := service.Initialize(projectDir)
	if err != nil {
		t.Fatalf("restore config: %v", err)
	}
	if result.AlreadyInitialized {
		t.Fatalf("expected repair result")
	}
	assertFileContent(t, importantSpec, "# Keep this\n")
	if _, err := os.Stat(filepath.Join(projectDir, "darp.yml")); err != nil {
		t.Fatalf("expected restored darp.yml: %v", err)
	}
}

func TestInitializeRegistersDiscoveredAssets(t *testing.T) {
	root := t.TempDir()
	projectDir := filepath.Join(root, "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(projectDir, ".github"), 0o755); err != nil {
		t.Fatalf("mkdir .github: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, ".github", "copilot-instructions.md"), []byte("# instructions\n"), 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}

	result, err := NewService(NewOSFileSystem()).Initialize(projectDir)
	if err != nil {
		t.Fatalf("initialize: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(projectDir, "darp.yml"))
	if err != nil {
		t.Fatalf("read darp.yml: %v", err)
	}
	text := string(content)
	for _, want := range []string{
		"assets:",
		"- path: .github/copilot-instructions.md",
		"family: github",
		"type: instruction",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected %q in darp.yml:\n%s", want, text)
		}
	}
	if !containsMessage(result, "Registered asset .github/copilot-instructions.md (github/instruction)") {
		t.Fatalf("expected registered asset message, got %#v", result.Messages)
	}
}

func TestInitializeDoesNotDuplicateAssetsOnRerun(t *testing.T) {
	root := t.TempDir()
	projectDir := filepath.Join(root, "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(projectDir, ".github"), 0o755); err != nil {
		t.Fatalf("mkdir .github: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, ".github", "copilot-instructions.md"), []byte("# instructions\n"), 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}

	if _, err := NewService(NewOSFileSystem()).Initialize(projectDir); err != nil {
		t.Fatalf("first initialize: %v", err)
	}
	result, err := NewService(NewOSFileSystem()).Initialize(projectDir)
	if err != nil {
		t.Fatalf("second initialize: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(projectDir, "darp.yml"))
	if err != nil {
		t.Fatalf("read darp.yml: %v", err)
	}
	if got := strings.Count(string(content), "path: .github/copilot-instructions.md"); got != 1 {
		t.Fatalf("expected a single asset entry, got %d:\n%s", got, content)
	}
	if !containsMessage(result, "Asset already registered .github/copilot-instructions.md (github/instruction)") {
		t.Fatalf("expected already registered message, got %#v", result.Messages)
	}
}

func TestInitializeReportsAmbiguousAssetsWithoutPersisting(t *testing.T) {
	root := t.TempDir()
	projectDir := filepath.Join(root, "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "AGENTS.md"), []byte("# agents\n"), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}

	result, err := NewService(NewOSFileSystem()).Initialize(projectDir)
	if err != nil {
		t.Fatalf("initialize: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(projectDir, "darp.yml"))
	if err != nil {
		t.Fatalf("read darp.yml: %v", err)
	}
	if strings.Contains(string(content), "assets:") {
		t.Fatalf("ambiguous asset must not be persisted:\n%s", content)
	}
	if !containsMessage(result, "Skipped ambiguous asset AGENTS.md") {
		t.Fatalf("expected ambiguous message, got %#v", result.Messages)
	}
}

func TestInitializeReportsSemanticConflictsWithoutPersisting(t *testing.T) {
	projectDir := filepath.Join(t.TempDir(), "demo-project")
	if err := os.Mkdir(projectDir, 0o755); err != nil {
		t.Fatalf("mkdir project dir: %v", err)
	}
	writeFixture := func(relative string) {
		path := filepath.Join(projectDir, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir fixture: %v", err)
		}
		if err := os.WriteFile(path, []byte("# review\n"), 0o644); err != nil {
			t.Fatalf("write fixture: %v", err)
		}
	}
	writeFixture(".github/instructions/review.instructions.md")
	writeFixture(".github/instructions/nested/review.instructions.md")

	result, err := NewService(NewOSFileSystem()).Initialize(projectDir)
	if err != nil {
		t.Fatalf("initialize: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(projectDir, "darp.yml"))
	if err != nil {
		t.Fatalf("read darp.yml: %v", err)
	}
	if strings.Contains(string(content), "assets:") || !containsMessage(result, "Skipped conflicting asset name review") {
		t.Fatalf("expected semantic conflict without persistence, config=%s messages=%#v", content, result.Messages)
	}
}

func containsMessage(result Result, want string) bool {
	for _, message := range result.Messages {
		if strings.Contains(message, want) {
			return true
		}
	}
	return false
}

func assertFileContent(t *testing.T, path string, want string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	if string(content) != want {
		t.Fatalf("unexpected content for %s\nwant:\n%s\ngot:\n%s", path, want, string(content))
	}
}

func assertDirectoryExists(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}

	if !info.IsDir() {
		t.Fatalf("expected directory at %s", path)
	}
}
