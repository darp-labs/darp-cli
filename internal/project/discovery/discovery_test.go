package discovery

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/darpbr/darp-cli/internal/project/asset"
)

func TestDiscoverEmptyRepository(t *testing.T) {
	root := t.TempDir()
	report := Discover(root)
	if !report.Empty() {
		t.Fatalf("expected empty report, got %#v", report)
	}
}

func TestDiscoverIgnoresUnsupportedFilesInKnownDirectories(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, ".github/workflows/ci.yml")
	writeFixture(t, root, ".github/CODEOWNERS")
	writeFixture(t, root, ".github/copilot-instructions.md")

	report := Discover(root)

	want := asset.Asset{
		Path:   ".github/copilot-instructions.md",
		Family: asset.FamilyGitHub,
		Type:   asset.TypeInstruction,
	}
	if !reflect.DeepEqual(report.Found, []asset.Asset{want}) {
		t.Fatalf("unexpected found assets\nwant: %#v\ngot:  %#v", []asset.Asset{want}, report.Found)
	}
	if len(report.Ambiguous) != 0 || len(report.SemanticConflicts) != 0 {
		t.Fatalf("unsupported files must not create discovery conflicts: %#v", report)
	}
}

func TestDiscoverSingleRootFileAsset(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, ".github/copilot-instructions.md")

	report := Discover(root)
	if len(report.Found) != 1 {
		t.Fatalf("expected one found asset, got %#v", report.Found)
	}
	got := report.Found[0]
	want := asset.Asset{Path: ".github/copilot-instructions.md", Family: asset.FamilyGitHub, Type: asset.TypeInstruction}
	if got != want {
		t.Fatalf("unexpected asset\nwant: %#v\ngot:  %#v", want, got)
	}
	if len(report.Ambiguous) != 0 || len(report.Ignored) != 0 {
		t.Fatalf("unexpected report %#v", report)
	}
}

func TestDiscoverRecursiveAndMultiplePatterns(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, ".github/instructions/nested/caveman.instructions.md")
	writeFixture(t, root, ".github/prompts/greet.prompt.md")
	writeFixture(t, root, ".github/skills/review/SKILL.md")
	writeFixture(t, root, ".github/hooks/on-push.json")
	writeFixture(t, root, ".github/agents/coder.md")
	writeFixture(t, root, ".claude/commands/prompt.md")
	writeFixture(t, root, ".claude/skills/docs/SKILL.md")
	writeFixture(t, root, ".claude/agents/explorer.md")
	writeFixture(t, root, ".claude/settings.json")
	writeFixture(t, root, ".codex/skills/x/SKILL.md")
	writeFixture(t, root, ".cursor/rules/backend.mdc")
	writeFixture(t, root, ".cursor/skills/y/SKILL.md")
	writeFixture(t, root, ".cursor/hooks.json")
	writeFixture(t, root, ".gemini/commands/test.toml")
	writeFixture(t, root, ".gemini/skills/z/SKILL.md")
	writeFixture(t, root, ".gemini/settings.json")
	writeFixture(t, root, ".gemini/extensions/foo/gemini-extension.json")

	report := Discover(root)
	if len(report.Found) != 17 {
		t.Fatalf("expected 17 found assets, got %d: %#v", len(report.Found), report.Found)
	}

	byPath := make(map[string]asset.Asset, len(report.Found))
	for _, a := range report.Found {
		byPath[a.Path] = a
	}

	expect := map[string]asset.Asset{
		".github/instructions/nested/caveman.instructions.md": {Path: ".github/instructions/nested/caveman.instructions.md", Family: asset.FamilyGitHub, Type: asset.TypeInstruction},
		".github/prompts/greet.prompt.md":                     {Path: ".github/prompts/greet.prompt.md", Family: asset.FamilyGitHub, Type: asset.TypePrompt},
		".github/skills/review/SKILL.md":                      {Path: ".github/skills/review/SKILL.md", Family: asset.FamilyGitHub, Type: asset.TypeSkill},
		".github/hooks/on-push.json":                          {Path: ".github/hooks/on-push.json", Family: asset.FamilyGitHub, Type: asset.TypeHook},
		".github/agents/coder.md":                             {Path: ".github/agents/coder.md", Family: asset.FamilyGitHub, Type: asset.TypeAgent},
		".claude/commands/prompt.md":                          {Path: ".claude/commands/prompt.md", Family: asset.FamilyClaude, Type: asset.TypeCommand},
		".claude/skills/docs/SKILL.md":                        {Path: ".claude/skills/docs/SKILL.md", Family: asset.FamilyClaude, Type: asset.TypeSkill},
		".claude/agents/explorer.md":                          {Path: ".claude/agents/explorer.md", Family: asset.FamilyClaude, Type: asset.TypeAgent},
		".claude/settings.json":                               {Path: ".claude/settings.json", Family: asset.FamilyClaude, Type: asset.TypeHook},
		".codex/skills/x/SKILL.md":                            {Path: ".codex/skills/x/SKILL.md", Family: asset.FamilyCodex, Type: asset.TypeSkill},
		".cursor/rules/backend.mdc":                           {Path: ".cursor/rules/backend.mdc", Family: asset.FamilyCursor, Type: asset.TypeRule},
		".cursor/skills/y/SKILL.md":                           {Path: ".cursor/skills/y/SKILL.md", Family: asset.FamilyCursor, Type: asset.TypeSkill},
		".cursor/hooks.json":                                  {Path: ".cursor/hooks.json", Family: asset.FamilyCursor, Type: asset.TypeHook},
		".gemini/commands/test.toml":                          {Path: ".gemini/commands/test.toml", Family: asset.FamilyGemini, Type: asset.TypeCommand},
		".gemini/skills/z/SKILL.md":                           {Path: ".gemini/skills/z/SKILL.md", Family: asset.FamilyGemini, Type: asset.TypeSkill},
		".gemini/settings.json":                               {Path: ".gemini/settings.json", Family: asset.FamilyGemini, Type: asset.TypeConfig},
		".gemini/extensions/foo/gemini-extension.json":        {Path: ".gemini/extensions/foo/gemini-extension.json", Family: asset.FamilyGemini, Type: asset.TypeExtension},
	}
	if !reflect.DeepEqual(byPath, expect) {
		t.Fatalf("unexpected assets\nwant: %#v\ngot:  %#v", expect, byPath)
	}
}

func TestDiscoverAmbiguousRootFile(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "AGENTS.md")

	report := Discover(root)
	if len(report.Found) != 0 {
		t.Fatalf("expected no found assets, got %#v", report.Found)
	}
	if len(report.Ambiguous) != 1 || report.Ambiguous[0].Path != "AGENTS.md" {
		t.Fatalf("expected AGENTS.md ambiguity, got %#v", report.Ambiguous)
	}
	if len(report.Ambiguous[0].Assets) != 3 {
		t.Fatalf("expected three families for AGENTS.md, got %#v", report.Ambiguous[0].Assets)
	}
}

func TestDiscoverReportsSemanticConflictAcrossPaths(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, ".github/instructions/review.instructions.md")
	writeFixture(t, root, ".github/instructions/nested/review.instructions.md")

	report := Discover(root)
	if len(report.Found) != 0 {
		t.Fatalf("expected conflicting assets not to be found, got %#v", report.Found)
	}
	if len(report.SemanticConflicts) != 1 {
		t.Fatalf("expected one semantic conflict, got %#v", report.SemanticConflicts)
	}
	conflict := report.SemanticConflicts[0]
	if conflict.Category != string(asset.TypeInstruction) || conflict.Name != "review" || len(conflict.Assets) != 2 {
		t.Fatalf("unexpected semantic conflict %#v", conflict)
	}
}

func TestDiscoverIgnoresSymlinks(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.md")
	if err := os.WriteFile(target, []byte("x\n"), 0o644); err != nil {
		t.Fatalf("write target: %v", err)
	}
	commands := filepath.Join(root, ".claude", "commands")
	if err := os.MkdirAll(commands, 0o755); err != nil {
		t.Fatalf("mkdir commands: %v", err)
	}
	link := filepath.Join(commands, "link.md")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	report := Discover(root)
	if len(report.Found) != 0 {
		t.Fatalf("expected no found assets, got %#v", report.Found)
	}
	if len(report.Ignored) != 1 || report.Ignored[0].Path != ".claude/commands/link.md" {
		t.Fatalf("expected ignored symlink, got %#v", report.Ignored)
	}
}

func TestDiscoverIgnoresSymlinkedFamilyRoot(t *testing.T) {
	root := t.TempDir()
	realDir := filepath.Join(root, "real-claude")
	writeFixture(t, root, "real-claude/commands/prompt.md")
	if err := os.Symlink(realDir, filepath.Join(root, ".claude")); err != nil {
		t.Fatalf("symlink root: %v", err)
	}

	report := Discover(root)
	if len(report.Found) != 0 {
		t.Fatalf("expected no found assets, got %#v", report.Found)
	}
	if len(report.Ignored) != 1 || report.Ignored[0].Path != ".claude" {
		t.Fatalf("expected ignored .claude symlink, got %#v", report.Ignored)
	}
}

func TestDiscoverDoesNotTraverseMonorepo(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "services/api/.github/copilot-instructions.md")

	report := Discover(root)
	if !report.Empty() {
		t.Fatalf("expected empty report for monorepo subtree, got %#v", report)
	}
}

func TestDiscoverIsDeterministic(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, ".github/copilot-instructions.md")
	writeFixture(t, root, ".claude/commands/one.md")
	writeFixture(t, root, ".claude/commands/two.md")

	first := Discover(root)
	second := Discover(root)
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("discovery is not deterministic\nfirst: %#v\nsecond: %#v", first, second)
	}
}

func TestMatchGlob(t *testing.T) {
	tests := []struct {
		glob string
		rel  string
		want bool
	}{
		{".github/copilot-instructions.md", ".github/copilot-instructions.md", true},
		{".github/copilot-instructions.md", ".github/other.md", false},
		{".github/instructions/**/*.instructions.md", ".github/instructions/a/b/c.instructions.md", true},
		{".github/instructions/**/*.instructions.md", ".github/instructions/c.instructions.md", true},
		{".github/skills/**/SKILL.md", ".github/skills/review/SKILL.md", true},
		{".github/skills/**/SKILL.md", ".github/skills/SKILL.md", true},
		{".github/skills/**/SKILL.md", ".github/skills/other.md", false},
		{".github/hooks/*.json", ".github/hooks/on-push.json", true},
		{".github/hooks/*.json", ".github/hooks/sub/on-push.json", false},
		{".gemini/extensions/*/gemini-extension.json", ".gemini/extensions/foo/gemini-extension.json", true},
	}
	for _, test := range tests {
		if got := matchGlob(test.glob, test.rel); got != test.want {
			t.Fatalf("matchGlob(%q, %q) = %v, want %v", test.glob, test.rel, got, test.want)
		}
	}
}

func writeFixture(t *testing.T, root, rel string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir fixture: %v", err)
	}
	if err := os.WriteFile(full, []byte("fixture\n"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
}
