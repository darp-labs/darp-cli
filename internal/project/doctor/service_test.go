package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiagnoseHealthyProject(t *testing.T) {
	root := healthyProject(t, "1.0")
	result := NewService().Diagnose(root)

	passed, warnings, errors := result.Counts()
	if passed != 8 || warnings != 0 || errors != 0 || result.ExitCode() != 0 {
		t.Fatalf("unexpected healthy result: %#v", result)
	}
}

func TestDiagnoseOlderContractWarns(t *testing.T) {
	result := NewService().Diagnose(healthyProject(t, "0.9"))
	passed, warnings, errors := result.Counts()
	if passed != 7 || warnings != 1 || errors != 0 || result.ExitCode() != 0 {
		t.Fatalf("unexpected warning result: %#v", result)
	}
}

func TestDiagnoseInvalidProjectReportsFailuresAndDoesNotWrite(t *testing.T) {
	root := healthyProject(t, "2.0")
	workflow := filepath.Join(root, ".darp", "workflows", "implement.yaml")
	if err := os.Remove(workflow); err != nil {
		t.Fatalf("remove workflow: %v", err)
	}
	before, err := projectSnapshot(root)
	if err != nil {
		t.Fatalf("snapshot before: %v", err)
	}

	result := NewService().Diagnose(root)

	after, err := projectSnapshot(root)
	if err != nil {
		t.Fatalf("snapshot after: %v", err)
	}
	if before != after {
		t.Fatalf("doctor modified the project\nbefore: %q\nafter: %q", before, after)
	}
	passed, warnings, errors := result.Counts()
	if passed != 6 || warnings != 0 || errors != 2 || result.ExitCode() != 1 {
		t.Fatalf("unexpected failing result: %#v", result)
	}
}

func TestDiagnoseRejectsConfigurationOutsideProject(t *testing.T) {
	root := healthyProject(t, "1.0")
	config := `version: "1.0"
project:
  name: demo
governance:
  lifecycle: /tmp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
`
	if err := os.WriteFile(filepath.Join(root, "darp.yml"), []byte(config), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	result := NewService().Diagnose(root)
	if result.Checks[0].State != Fail {
		t.Fatalf("expected configuration failure, got %#v", result.Checks[0])
	}
}

func TestDiagnoseRejectsSkillTraversal(t *testing.T) {
	root := healthyProject(t, "1.0")
	workflow := filepath.Join(root, ".darp", "workflows", "implement.yaml")
	if err := os.WriteFile(workflow, []byte("name: implement\nsteps:\n  - ../outside\n"), 0o644); err != nil {
		t.Fatalf("write workflow: %v", err)
	}
	result := NewService().Diagnose(root)
	if result.Checks[2].State != Fail {
		t.Fatalf("expected workflow failure, got %#v", result.Checks[2])
	}
}

func TestDiagnoseRejectsSkillMarkdownDirectory(t *testing.T) {
	root := healthyProject(t, "1.0")
	path := filepath.Join(root, ".agents", "skills", "documentation", "SKILL.md")
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove skill file: %v", err)
	}
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("create skill directory: %v", err)
	}
	result := NewService().Diagnose(root)
	if result.Checks[3].State != Fail {
		t.Fatalf("expected skills failure, got %#v", result.Checks[3])
	}
}

func TestDiagnoseRejectsMissingDefaultWorkflow(t *testing.T) {
	root := healthyProject(t, "1.0")
	configPath := filepath.Join(root, "darp.yml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	content = []byte(string(content) + "\n")
	content = []byte(strings.Replace(string(content), "default: implement", "default: missing", 1))
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	result := NewService().Diagnose(root)
	if result.Checks[2].State != Fail || !strings.Contains(result.Checks[2].Message, `workflow "missing" not found`) {
		t.Fatalf("expected missing default workflow failure, got %#v", result.Checks[2])
	}
}

func TestDiagnoseRejectsDuplicateWorkflowNames(t *testing.T) {
	root := healthyProject(t, "1.0")
	if err := os.WriteFile(filepath.Join(root, ".darp", "workflows", "other.yaml"), []byte("name: implement\nsteps:\n  - documentation\n"), 0o644); err != nil {
		t.Fatalf("write duplicate workflow: %v", err)
	}

	result := NewService().Diagnose(root)
	if result.Checks[2].State != Fail || !strings.Contains(result.Checks[2].Message, "duplicate workflow name") {
		t.Fatalf("expected duplicate workflow failure, got %#v", result.Checks[2])
	}
}

func TestDiagnoseRejectsInvalidWorkflowSteps(t *testing.T) {
	tests := map[string]string{
		"empty list":         "name: implement\nsteps: []\n",
		"empty item":         "name: implement\nsteps:\n  - \" \"\n",
		"invalid skill path": "name: implement\nsteps:\n  - ../documentation\n",
		"invalid yaml type":  "name: implement\nsteps: documentation\n",
	}
	for name, workflow := range tests {
		t.Run(name, func(t *testing.T) {
			root := healthyProject(t, "1.0")
			if err := os.WriteFile(filepath.Join(root, ".darp", "workflows", "implement.yaml"), []byte(workflow), 0o644); err != nil {
				t.Fatalf("write workflow: %v", err)
			}
			result := NewService().Diagnose(root)
			if result.Checks[2].State != Fail {
				t.Fatalf("expected workflow failure, got %#v", result.Checks[2])
			}
		})
	}
}

func TestDiagnoseRejectsOptionalSkillResourceFile(t *testing.T) {
	root := healthyProject(t, "1.0")
	path := filepath.Join(root, ".agents", "skills", "documentation", "references")
	if err := os.WriteFile(path, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("write optional resource: %v", err)
	}

	result := NewService().Diagnose(root)
	if result.Checks[3].State != Fail || !strings.Contains(result.Checks[3].Message, "references must be a directory") {
		t.Fatalf("expected optional resource failure, got %#v", result.Checks[3])
	}
}

func TestDiagnoseRejectsSkillRootFile(t *testing.T) {
	root := healthyProject(t, "1.0")
	if err := os.WriteFile(filepath.Join(root, ".agents", "skills", "not-a-skill"), []byte("invalid"), 0o644); err != nil {
		t.Fatalf("write skill root file: %v", err)
	}

	result := NewService().Diagnose(root)
	if result.Checks[3].State != Fail {
		t.Fatalf("expected skill root failure, got %#v", result.Checks[3])
	}
}

func TestDiagnoseValidAssetPasses(t *testing.T) {
	root := healthyProject(t, "1.0")
	rewriteConfig(t, root, `version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
assets:
  - path: .github/copilot-instructions.md
    family: github
    type: instruction
`)
	if err := os.MkdirAll(filepath.Join(root, ".github"), 0o755); err != nil {
		t.Fatalf("mkdir .github: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".github", "copilot-instructions.md"), []byte("# x\n"), 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}
	result := NewService().Diagnose(root)
	if check := result.Checks[7]; check.State != Pass {
		t.Fatalf("expected assets pass, got %#v", check)
	}
}

func TestDiagnoseMissingAssetWarns(t *testing.T) {
	root := healthyProject(t, "1.0")
	rewriteConfig(t, root, `version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
assets:
  - path: .github/copilot-instructions.md
    family: github
    type: instruction
`)
	result := NewService().Diagnose(root)
	if check := result.Checks[7]; check.State != Warning || !strings.Contains(check.Message, "missing") {
		t.Fatalf("expected missing asset warning, got %#v", check)
	}
	if result.ExitCode() != 0 {
		t.Fatalf("missing asset must not block the project, got %#v", result.Checks)
	}
}

func TestDiagnoseDirectoryAssetFails(t *testing.T) {
	root := healthyProject(t, "1.0")
	rewriteConfig(t, root, `version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
assets:
  - path: .github/copilot-instructions.md
    family: github
    type: instruction
`)
	if err := os.MkdirAll(filepath.Join(root, ".github", "copilot-instructions.md"), 0o755); err != nil {
		t.Fatalf("mkdir asset path: %v", err)
	}

	result := NewService().Diagnose(root)
	if check := result.Checks[7]; check.State != Fail || !strings.Contains(check.Message, "regular file") {
		t.Fatalf("expected directory asset failure, got %#v", check)
	}
	if result.ExitCode() != 1 {
		t.Fatalf("directory asset must block the project, got %#v", result.Checks)
	}
}

func TestDiagnoseUnknownAssetTypeFails(t *testing.T) {
	root := healthyProject(t, "1.0")
	rewriteConfig(t, root, `version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
assets:
  - path: .github/copilot-instructions.md
    family: github
    type: bogus
`)
	if err := os.MkdirAll(filepath.Join(root, ".github"), 0o755); err != nil {
		t.Fatalf("mkdir .github: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".github", "copilot-instructions.md"), []byte("# x\n"), 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}
	result := NewService().Diagnose(root)
	if check := result.Checks[7]; check.State != Fail || !strings.Contains(check.Message, "unknown type") {
		t.Fatalf("expected unknown type failure, got %#v", check)
	}
}

func TestDiagnoseIncompatibleAssetClassificationFails(t *testing.T) {
	root := healthyProject(t, "1.0")
	rewriteConfig(t, root, `version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
assets:
  - path: .github/copilot-instructions.md
    family: claude
    type: command
`)
	if err := os.MkdirAll(filepath.Join(root, ".github"), 0o755); err != nil {
		t.Fatalf("mkdir .github: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".github", "copilot-instructions.md"), []byte("# x\n"), 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}
	result := NewService().Diagnose(root)
	if check := result.Checks[7]; check.State != Fail || !strings.Contains(check.Message, "incompatible") {
		t.Fatalf("expected incompatible asset failure, got %#v", check)
	}
}

func TestDiagnoseRegisteredSymlinkFails(t *testing.T) {
	root := healthyProject(t, "1.0")
	rewriteConfig(t, root, `version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
assets:
  - path: .github/copilot-instructions.md
    family: github
    type: instruction
`)
	if err := os.MkdirAll(filepath.Join(root, ".github"), 0o755); err != nil {
		t.Fatalf("mkdir .github: %v", err)
	}
	target := filepath.Join(root, "instructions.md")
	if err := os.WriteFile(target, []byte("# x\n"), 0o644); err != nil {
		t.Fatalf("write target: %v", err)
	}
	if err := os.Symlink(target, filepath.Join(root, ".github", "copilot-instructions.md")); err != nil {
		t.Fatalf("symlink asset: %v", err)
	}
	result := NewService().Diagnose(root)
	if check := result.Checks[7]; check.State != Fail || !strings.Contains(check.Message, "symlink") {
		t.Fatalf("expected symlink failure, got %#v", check)
	}
}

func TestDiagnoseAssetTraversalFails(t *testing.T) {
	root := healthyProject(t, "1.0")
	rewriteConfig(t, root, `version: "1.0"
project:
  name: demo
governance:
  lifecycle: .darp/lifecycle.md
workflows:
  default: implement
skills:
  documentation: .agents/skills/documentation
assets:
  - path: ../outside.md
    family: github
    type: instruction
`)
	result := NewService().Diagnose(root)
	if check := result.Checks[7]; check.State != Fail || !strings.Contains(check.Message, "remain inside the project") {
		t.Fatalf("expected traversal failure, got %#v", check)
	}
}

func rewriteConfig(t *testing.T, root, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, "darp.yml"), []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func healthyProject(t *testing.T, version string) string {
	t.Helper()
	root := t.TempDir()
	for _, directory := range []string{
		".darp/governance", ".darp/workflows", ".agents/skills/documentation", ".darp/templates",
	} {
		if err := os.MkdirAll(filepath.Join(root, directory), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", directory, err)
		}
	}
	files := map[string]string{
		"darp.yml":                              "version: \"" + version + "\"\nproject:\n  name: demo\ngovernance:\n  lifecycle: .darp/lifecycle.md\nworkflows:\n  default: implement\nskills:\n  documentation: .agents/skills/documentation\n",
		".darp/lifecycle.md":                    "# Lifecycle\n",
		".darp/governance/quality-gates.md":     "# Gates\n",
		".darp/workflows/implement.yaml":        "name: implement\nsteps:\n  - documentation\n",
		".agents/skills/documentation/SKILL.md": "---\nname: documentation\ndescription: test\n---\n# Skill\n",
	}
	for path, content := range files {
		if err := os.WriteFile(filepath.Join(root, path), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	return root
}

func projectSnapshot(root string) (string, error) {
	var snapshot string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		snapshot += relative + ":"
		if !entry.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			snapshot += string(content)
		}
		snapshot += "\n"
		return nil
	})
	return snapshot, err
}
