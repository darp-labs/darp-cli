package sync

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	projectinit "github.com/darpbr/darp-cli/internal/project/init"
)

func TestReconcileAddsLiteralSkillAndIsByteIdempotent(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "darp.yml"), "skills: {}\nassets: []\ncustom: keep\n")
	write(t, filepath.Join(root, ".agents/skills/my-skill/SKILL.md"), validSkill("my-skill"))
	s := Service{Root: root}
	r, err := s.Reconcile(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Added) != 1 || r.Added[0].Identity != "my-skill" {
		t.Fatalf("unexpected additions: %#v", r.Added)
	}
	first, err := os.ReadFile(filepath.Join(root, "darp.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(first), "my-skill: .agents/skills/my-skill") || !strings.Contains(string(first), "custom: keep") {
		t.Fatalf("manifest lost expected values:\n%s", first)
	}
	if _, err = s.Reconcile(false); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(filepath.Join(root, "darp.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatalf("second sync changed bytes:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestDryRunDoesNotWriteManifest(t *testing.T) {
	root := t.TempDir()
	manifest := filepath.Join(root, "darp.yml")
	original := []byte("skills: {}\nassets: []\n")
	if err := os.WriteFile(manifest, original, 0o644); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(root, ".agents/skills/preview/SKILL.md"), validSkill("preview"))
	r, err := (Service{Root: root}).Reconcile(true)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Added) != 1 {
		t.Fatalf("expected preview addition, got %#v", r.Added)
	}
	got, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Fatalf("dry-run modified manifest: %q", got)
	}
}

func TestSkillDiscoveryPreservesNamesAndSortsResults(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "darp.yml"), "skills: {}\nassets: []\n")
	write(t, filepath.Join(root, ".agents/skills/zed-name/SKILL.md"), validSkill("zed-name"))
	write(t, filepath.Join(root, ".agents/skills/a-b/SKILL.md"), validSkill("a-b"))
	r, err := (Service{Root: root}).Reconcile(true)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Added) != 2 || r.Added[0].Identity != "a-b" || r.Added[1].Identity != "zed-name" {
		t.Fatalf("skill identities/order changed: %#v", r.Added)
	}
}

func TestInvalidManifestIsUnchanged(t *testing.T) {
	root := t.TempDir()
	manifest := filepath.Join(root, "darp.yml")
	original := []byte("skills: {}\nskills: {}\nassets: []\n")
	if err := os.WriteFile(manifest, original, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := (Service{Root: root}).Reconcile(false); err == nil {
		t.Fatal("expected duplicate-key failure")
	}
	got, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Fatal("invalid manifest was modified")
	}
}

func TestMissingManifestDryRunSimulatesBootstrapWithoutWriting(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, ".agents/skills/extra/SKILL.md"), validSkill("extra"))
	before, err := snapshotTree(root)
	if err != nil {
		t.Fatal(err)
	}
	r, err := (Service{Root: root}).Reconcile(true)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Added) != 1 || r.Added[0].Identity != "extra" {
		t.Fatalf("unexpected simulated additions: %#v", r.Added)
	}
	if _, err := os.Stat(filepath.Join(root, "darp.yml")); !os.IsNotExist(err) {
		t.Fatalf("dry run created darp.yml: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".darp")); !os.IsNotExist(err) {
		t.Fatalf("dry run created bootstrap directory: %v", err)
	}
	after, err := snapshotTree(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("dry-run changed the filesystem\nbefore: %#v\nafter: %#v", before, after)
	}
}

func TestMissingManifestDryRunSerializesProjectNameAsYAML(t *testing.T) {
	root := filepath.Join(t.TempDir(), "quoted\"project\nname")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	r, err := (Service{Root: root}).Reconcile(true)
	if err != nil {
		t.Fatalf("dry-run rejected a project name requiring YAML escaping: %v", err)
	}
	if len(r.Kept) != 4 {
		t.Fatalf("expected four bootstrap skills, got %#v", r.Kept)
	}
}

func TestInvalidSkillEntriesAreRejectedWithoutWriting(t *testing.T) {
	tests := []struct {
		name, manifest string
	}{
		{
			name:     "mapping value",
			manifest: "skills:\n  invalid:\n    nested: value\nassets: []\n",
		},
		{
			name:     "sequence value",
			manifest: "skills:\n  invalid:\n    - .agents/skills/invalid\nassets: []\n",
		},
		{
			name:     "non-string value",
			manifest: "skills:\n  invalid: 42\nassets: []\n",
		},
		{
			name:     "non-string key",
			manifest: "skills:\n  42: .agents/skills/invalid\nassets: []\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			manifest := filepath.Join(root, "darp.yml")
			original := []byte(test.manifest)
			if err := os.WriteFile(manifest, original, 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := (Service{Root: root}).Reconcile(false); err == nil {
				t.Fatal("expected invalid skill entry to be rejected")
			}
			got, err := os.ReadFile(manifest)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != string(original) {
				t.Fatalf("invalid manifest changed: got %q, want %q", got, original)
			}
		})
	}
}

func TestNonCanonicalLegacySkillKeyIsPreserved(t *testing.T) {
	root := t.TempDir()
	manifest := filepath.Join(root, "darp.yml")
	original := []byte("skills:\n  ../manual: .agents/skills/../manual\nassets: []\n")
	if err := os.WriteFile(manifest, original, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := (Service{Root: root}).Reconcile(false); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "../manual: .agents/skills/../manual") {
		t.Fatalf("non-canonical manual entry was removed or rewritten:\n%s", got)
	}
}

func TestMissingManifestDryRunMatchesAppliedInventoryPlan(t *testing.T) {
	makeRoot := func() string {
		root := t.TempDir()
		write(t, filepath.Join(root, ".agents/skills/extra/SKILL.md"), validSkill("extra"))
		write(t, filepath.Join(root, ".github/copilot-instructions.md"), "instruction")
		return root
	}
	dryRoot := makeRoot()
	dry, err := (Service{Root: dryRoot}).Reconcile(true)
	if err != nil {
		t.Fatal(err)
	}
	applyRoot := makeRoot()
	if _, err := projectinit.NewService(projectinit.NewOSFileSystem()).Initialize(applyRoot); err != nil {
		t.Fatal(err)
	}
	applied, err := (Service{Root: applyRoot}).Reconcile(false)
	if err != nil {
		t.Fatal(err)
	}
	if !sameItems(dry.Added, applied.Added) || !sameItems(dry.Removed, applied.Removed) || !sameItems(dry.Kept, applied.Kept) {
		t.Fatalf("dry run differs from applied plan\ndry: %#v\napplied: %#v", dry, applied)
	}
}

func TestReconcileWithoutManifestRunsBootstrapAndSync(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, ".agents/skills/custom/SKILL.md"), validSkill("custom"))
	write(t, filepath.Join(root, ".github/copilot-instructions.md"), "instructions")
	result, err := (Service{Root: root}).Reconcile(false)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		"darp.yml",
		".darp/lifecycle.md",
		".agents/skills/documentation/SKILL.md",
	} {
		if _, err := os.Stat(filepath.Join(root, path)); err != nil {
			t.Fatalf("sync bootstrap did not create %s: %v", path, err)
		}
	}
	if len(result.Kept) < 1 || len(result.Added) != 1 || result.Added[0].Identity != "custom" {
		t.Fatalf("bootstrap and reconciliation result is incomplete: %#v", result)
	}
	manifest, err := os.ReadFile(filepath.Join(root, "darp.yml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"custom: .agents/skills/custom", ".github/copilot-instructions.md"} {
		if !strings.Contains(string(manifest), expected) {
			t.Fatalf("bootstrap sync manifest lacks %q:\n%s", expected, manifest)
		}
	}
}

func sameItems(a, b []Item) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestReconcilePreservesConflictingSkillAndRemovesManagedMissingEntry(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "darp.yml"), `skills:
  custom: elsewhere
  removed: .agents/skills/removed
assets: []
darp:
  sync:
    managed:
      skills:
        removed: .agents/skills/removed
`)
	write(t, filepath.Join(root, ".agents/skills/custom/SKILL.md"), validSkill("custom"))
	r, err := (Service{Root: root}).Reconcile(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Removed) != 1 || r.Removed[0].Identity != "removed" {
		t.Fatalf("unexpected removals: %#v", r.Removed)
	}
	if len(r.Ignored) == 0 {
		t.Fatal("expected conflicting manual skill to be reported")
	}
	got, err := os.ReadFile(filepath.Join(root, "darp.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "custom: elsewhere") || strings.Contains(string(got), "removed: .agents/skills/removed") {
		t.Fatalf("unexpected reconciled manifest:\n%s", got)
	}
}

func TestSymlinkSkillScopePreventsManagedRemoval(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "darp.yml"), `skills:
  linked: .agents/skills/linked
assets: []
darp:
  sync:
    managed:
      skills:
        linked: .agents/skills/linked
`)
	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".agents/skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(root, ".agents/skills/linked")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	r, err := (Service{Root: root}).Reconcile(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Removed) != 0 {
		t.Fatalf("symlink scope caused removal: %#v", r.Removed)
	}
	if len(r.Ignored) == 0 {
		t.Fatal("expected symlink to be reported")
	}
	got, err := os.ReadFile(filepath.Join(root, "darp.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "linked: .agents/skills/linked") {
		t.Fatalf("managed entry was removed:\n%s", got)
	}
}

func TestIncompleteSkillScopePreventsManagedRemoval(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "darp.yml"), `skills:
  partial: .agents/skills/partial
assets: []
darp:
  sync:
    managed:
      skills:
        partial: .agents/skills/partial
`)
	if err := os.MkdirAll(filepath.Join(root, ".agents/skills/partial"), 0o755); err != nil {
		t.Fatal(err)
	}
	r, err := (Service{Root: root}).Reconcile(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Removed) != 0 || len(r.Ignored) == 0 {
		t.Fatalf("incomplete skill was not treated as unknown: %#v", r)
	}
}

func TestNonRegularSkillCandidatePreventsManagedRemoval(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "darp.yml"), `skills:
  partial: .agents/skills/partial
assets: []
darp:
  sync:
    managed:
      skills:
        partial: .agents/skills/partial
`)
	if err := os.MkdirAll(filepath.Join(root, ".agents/skills/partial/SKILL.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	plan, err := (Service{Root: root}).BuildPlan()
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Removed) != 0 || len(plan.Ignored) == 0 || len(plan.Warnings) == 0 {
		t.Fatalf("non-regular skill candidate was not treated as unknown: %#v", plan.Result)
	}
}

func TestSemanticAssetConflictPreventsManagedRemoval(t *testing.T) {
	root := t.TempDir()
	managedPath := ".github/instructions/nested/review.instructions.md"
	write(t, filepath.Join(root, "darp.yml"), `skills: {}
assets:
  - path: .github/instructions/nested/review.instructions.md
    family: github
    type: instruction
darp:
  sync:
    managed:
      skills: {}
      assets:
        - path: .github/instructions/nested/review.instructions.md
          family: github
          type: instruction
`)
	write(t, filepath.Join(root, managedPath), "nested")
	write(t, filepath.Join(root, ".github/instructions/review.instructions.md"), "root")
	plan, err := (Service{Root: root}).BuildPlan()
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Removed) != 0 {
		t.Fatalf("semantic conflict caused a removal: %#v", plan.Removed)
	}
	foundConflict := false
	for _, item := range plan.Ignored {
		if item.Identity == "review" && item.Detail == "semantic asset conflict" {
			foundConflict = true
		}
	}
	if !foundConflict {
		t.Fatalf("semantic conflict was not reported: %#v", plan.Ignored)
	}
}

func TestInvalidAgentSkillsAreReportedAndPreserveUnknownOwnership(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "darp.yml"), `skills:
  BadName: .agents/skills/BadName
assets: []
darp:
  sync:
    managed:
      skills:
        BadName: .agents/skills/BadName
`)
	write(t, filepath.Join(root, ".agents/skills/BadName/SKILL.md"), "---\nname: BadName\ndescription: Invalid uppercase name.\n---\n")
	r, err := (Service{Root: root}).Reconcile(true)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Added) != 0 || len(r.Removed) != 0 {
		t.Fatalf("invalid skill changed ownership: %#v", r)
	}
	if len(r.Ignored) == 0 || r.Ignored[0].Detail != "invalid name or name does not match directory" {
		t.Fatalf("invalid skill was not reported: %#v", r.Ignored)
	}
}

func TestAgentSkillsFrontmatterRequirements(t *testing.T) {
	tests := []struct {
		name, directory, content, wantReason string
	}{
		{"missing frontmatter", "valid", "# instructions\n", "missing YAML frontmatter"},
		{"missing description", "valid", "---\nname: valid\n---\n", "description must contain 1 to 1024 characters"},
		{"name mismatch", "valid", "---\nname: other\ndescription: Valid description.\n---\n", "invalid name or name does not match directory"},
		{"malformed YAML", "valid", "---\nname: [\ndescription: broken\n---\n", "invalid YAML frontmatter"},
		{"metadata values must be strings", "valid", "---\nname: valid\ndescription: Valid description.\nmetadata:\n  version: 1\n---\n", "metadata must be a string map"},
		{"compatibility limit", "valid", "---\nname: valid\ndescription: Valid description.\ncompatibility: " + strings.Repeat("x", 501) + "\n---\n", "compatibility must contain 1 to 500 characters"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if reason := validateSkillDocument(test.directory, []byte(test.content)); reason != test.wantReason {
				t.Fatalf("reason=%q, want %q", reason, test.wantReason)
			}
		})
	}
}

func TestIncompleteExternalSkillPreventsManagedRemoval(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "darp.yml"), `skills: {}
assets:
  - path: .github/skills/incomplete/SKILL.md
    family: github
    type: skill
darp:
  sync:
    managed:
      skills: {}
      assets:
        - path: .github/skills/incomplete/SKILL.md
          family: github
          type: skill
`)
	if err := os.MkdirAll(filepath.Join(root, ".github/skills/incomplete"), 0o755); err != nil {
		t.Fatal(err)
	}
	r, err := (Service{Root: root}).Reconcile(true)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Removed) != 0 {
		t.Fatalf("incomplete external skill caused a removal: %#v", r.Removed)
	}
	if len(r.Ignored) == 0 || r.Ignored[0].Detail != "incomplete skill" {
		t.Fatalf("incomplete external skill was not reported: %#v", r.Ignored)
	}
}

func TestBuildPlanReportsReadFailuresAndProtectsAffectedScope(t *testing.T) {
	root := t.TempDir()
	manifest := `skills:
  custom: .agents/skills/custom
assets:
  - path: .github/skills/removed/SKILL.md
    family: github
    type: skill
darp:
  sync:
    managed:
      skills:
        custom: .agents/skills/custom
      assets:
        - path: .github/skills/removed/SKILL.md
          family: github
          type: skill
`
	write(t, filepath.Join(root, "darp.yml"), manifest)
	write(t, filepath.Join(root, ".github/copilot-instructions.md"), "instructions")
	write(t, filepath.Join(root, ".agents/skills/custom/SKILL.md"), validSkill("custom"))
	original, err := os.ReadFile(filepath.Join(root, "darp.yml"))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name             string
		fileSystem       failingScanFS
		wantAssetRemoved bool
	}{
		{
			name:             "agents root listing",
			wantAssetRemoved: true,
			fileSystem: failingScanFS{readDirErrors: map[string]error{
				filepath.Join(root, ".agents/skills"): os.ErrPermission,
			}},
		},
		{
			name:             "skill content read",
			wantAssetRemoved: true,
			fileSystem: failingScanFS{readFileErrors: map[string]error{
				filepath.Join(root, ".agents/skills/custom/SKILL.md"): os.ErrPermission,
			}},
		},
		{
			name: "external root walk",
			fileSystem: failingScanFS{
				walkRoot: filepath.Join(root, ".github"),
				walkPath: filepath.Join(root, ".github/skills"),
				walkErr:  os.ErrPermission,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			plan, err := (Service{Root: root, fileSystem: test.fileSystem}).BuildPlan()
			if err != nil {
				t.Fatal(err)
			}
			if len(plan.Ignored) == 0 || len(plan.Warnings) == 0 {
				t.Fatalf("read failure did not produce ignored item and warning: %#v", plan.Result)
			}
			if test.wantAssetRemoved && (len(plan.Removed) != 1 || plan.Removed[0].Identity != ".github/skills/removed/SKILL.md (github/skill)") {
				t.Fatalf("fully observed asset scope was not synchronized: %#v", plan.Removed)
			}
			if !test.wantAssetRemoved && len(plan.Removed) != 0 {
				t.Fatalf("unreadable external scope caused removal: %#v", plan.Removed)
			}
			for _, removed := range plan.Removed {
				if removed.Identity == "custom" {
					t.Fatalf("unreadable skill scope caused removal: %#v", plan.Removed)
				}
			}
			got, err := os.ReadFile(filepath.Join(root, "darp.yml"))
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != string(original) {
				t.Fatal("BuildPlan modified darp.yml")
			}
		})
	}
}

func TestBuildPlanReportsExternalCandidateStatFailureWithoutRemoval(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".github/copilot-instructions.md")
	write(t, path, "instructions")
	write(t, filepath.Join(root, "darp.yml"), `skills: {}
assets:
  - path: .github/copilot-instructions.md
    family: github
    type: instruction
darp:
  sync:
    managed:
      skills: {}
      assets:
        - path: .github/copilot-instructions.md
          family: github
          type: instruction
`)
	plan, err := (Service{Root: root, fileSystem: failingScanFS{lstatErrors: map[string]error{path: os.ErrPermission}}}).BuildPlan()
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Removed) != 0 || len(plan.Ignored) == 0 || len(plan.Warnings) == 0 {
		t.Fatalf("candidate stat failure was not treated as unknown: %#v", plan.Result)
	}
}

func TestNonRegularExternalAssetCandidatePreventsManagedRemoval(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "darp.yml"), `skills: {}
assets:
  - path: .github/copilot-instructions.md
    family: github
    type: instruction
darp:
  sync:
    managed:
      skills: {}
      assets:
        - path: .github/copilot-instructions.md
          family: github
          type: instruction
`)
	if err := os.MkdirAll(filepath.Join(root, ".github/copilot-instructions.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	plan, err := (Service{Root: root}).BuildPlan()
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Removed) != 0 || len(plan.Ignored) == 0 || len(plan.Warnings) == 0 {
		t.Fatalf("non-regular external candidate was not treated as unknown: %#v", plan.Result)
	}
}

func TestBuildPlanReturnsStructuredFailureAndApplyRejectsStalePlan(t *testing.T) {
	root := t.TempDir()
	manifest := filepath.Join(root, "darp.yml")
	write(t, manifest, "skills: {}\nassets: []\n")
	service := Service{Root: root}
	plan, err := service.BuildPlan()
	if err != nil {
		t.Fatal(err)
	}
	write(t, manifest, "skills: {}\nassets: []\nexternal: changed\n")
	got, err := service.Apply(plan)
	if err == nil {
		t.Fatal("expected stale plan to be rejected")
	}
	if len(got.Failures) != 1 || got.Failures[0].Identity != "darp.yml" {
		t.Fatalf("expected structured apply failure, got %#v", got.Failures)
	}
	content, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "external: changed") {
		t.Fatalf("Apply overwrote changes made after planning:\n%s", content)
	}
}

func TestBuildPlanIsDeterministicAndClassifiesInvalidManifest(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "darp.yml"), "skills: {}\nassets: []\n")
	write(t, filepath.Join(root, ".agents/skills/one/SKILL.md"), validSkill("one"))
	write(t, filepath.Join(root, ".agents/skills/two/SKILL.md"), validSkill("two"))
	service := Service{Root: root}
	first, err := service.BuildPlan()
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.BuildPlan()
	if err != nil {
		t.Fatal(err)
	}
	if !sameItems(first.Added, second.Added) || !sameItems(first.Removed, second.Removed) || !sameItems(first.Kept, second.Kept) || !sameItems(first.Ignored, second.Ignored) || !sameItems(first.Warnings, second.Warnings) {
		t.Fatalf("same workspace produced different plans\nfirst: %#v\nsecond: %#v", first.Result, second.Result)
	}
	write(t, filepath.Join(root, "darp.yml"), "skills: [\n")
	failed, err := service.BuildPlan()
	if err == nil {
		t.Fatal("expected invalid manifest failure")
	}
	if len(failed.Failures) != 1 || failed.Failures[0].Identity != "darp.yml" || failed.Failures[0].Kind != "failure" {
		t.Fatalf("invalid manifest was not classified in the plan: %#v", failed.Failures)
	}
}

func TestSyncMigratesRemovedGitHubSkillAndAddsTwoSkillsWithoutChangingHarness(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "darp.yml"), `skills: {}
assets:
  - path: .github/skills/legacy/SKILL.md
    family: github
    type: skill
`)
	write(t, filepath.Join(root, ".agents/skills/alpha/SKILL.md"), validSkill("alpha"))
	write(t, filepath.Join(root, ".agents/skills/beta/SKILL.md"), validSkill("beta"))
	harness := []byte("this Harness file is intentionally opaque and must not change\n")
	write(t, filepath.Join(root, ".darp/harness.yaml"), string(harness))
	skillsBefore, err := snapshotTree(filepath.Join(root, ".agents/skills"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := (Service{Root: root}).Reconcile(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Added) != 2 || result.Added[0].Identity != "alpha" || result.Added[1].Identity != "beta" {
		t.Fatalf("expected both DARP skills to be added: %#v", result.Added)
	}
	if len(result.Removed) != 1 || result.Removed[0].Identity != ".github/skills/legacy/SKILL.md (github/skill)" {
		t.Fatalf("expected absent legacy GitHub skill to be removed: %#v", result.Removed)
	}
	config, err := os.ReadFile(filepath.Join(root, "darp.yml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"alpha", "beta"} {
		if !strings.Contains(string(config), name+": .agents/skills/"+name) {
			t.Fatalf("skill %q not written to darp.yml:\n%s", name, config)
		}
	}
	if strings.Contains(string(config), ".github/skills/legacy/SKILL.md") {
		t.Fatalf("removed GitHub skill remains in darp.yml:\n%s", config)
	}
	gotHarness, err := os.ReadFile(filepath.Join(root, ".darp/harness.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(gotHarness) != string(harness) {
		t.Fatalf("sync modified Harness manifest: %q", gotHarness)
	}
	skillsAfter, err := snapshotTree(filepath.Join(root, ".agents/skills"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(skillsBefore, skillsAfter) {
		t.Fatalf("sync modified source skills\nbefore: %#v\nafter: %#v", skillsBefore, skillsAfter)
	}
}

func TestAssetLegacyAdoptionPreservesMetadataThenRemovesMissingEntry(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".github/copilot-instructions.md")
	write(t, path, "instructions")
	write(t, filepath.Join(root, "darp.yml"), `skills: {}
assets:
  - path: .github/copilot-instructions.md
    family: github
    type: instruction
    note: keep-me
`)
	s := Service{Root: root}
	first, err := s.Reconcile(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Kept) != 1 {
		t.Fatalf("expected legacy asset adoption, got %#v", first)
	}
	manifest, err := os.ReadFile(filepath.Join(root, "darp.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifest), "note: keep-me") || !strings.Contains(string(manifest), "family: github") {
		t.Fatalf("asset metadata not preserved:\n%s", manifest)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	second, err := s.Reconcile(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(second.Removed) != 1 {
		t.Fatalf("expected managed legacy asset removal, got %#v", second.Removed)
	}
}

func TestAtomicWriteFailurePreservesOriginalAndCleansTemporaryFile(t *testing.T) {
	root := t.TempDir()
	manifest := filepath.Join(root, "darp.yml")
	original := []byte("original\n")
	if err := os.WriteFile(manifest, original, 0o640); err != nil {
		t.Fatal(err)
	}
	wantErr := errors.New("injected rename failure")
	err := writeAtomic(manifest, []byte("replacement\n"), 0o640, func(string, string) error { return wantErr })
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected injected error, got %v", err)
	}
	got, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Fatalf("original changed: %q", got)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "darp.yml" {
		t.Fatalf("temporary file was not cleaned: %#v", entries)
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func validSkill(name string) string {
	return "---\nname: " + name + "\ndescription: Test skill used to validate inventory synchronization.\n---\n# Instructions\n"
}

func snapshotTree(root string) (map[string]string, error) {
	snapshot := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			snapshot[relative] = "directory"
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			snapshot[relative] = "symlink:" + target
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		snapshot[relative] = "file:" + string(content)
		return nil
	})
	return snapshot, err
}

type failingScanFS struct {
	osScanFS
	lstatErrors    map[string]error
	readDirErrors  map[string]error
	readFileErrors map[string]error
	walkRoot       string
	walkPath       string
	walkErr        error
}

func (f failingScanFS) Lstat(path string) (os.FileInfo, error) {
	if err := f.lstatErrors[filepath.Clean(path)]; err != nil {
		return nil, err
	}
	return f.osScanFS.Lstat(path)
}

func (f failingScanFS) ReadDir(path string) ([]os.DirEntry, error) {
	if err := f.readDirErrors[filepath.Clean(path)]; err != nil {
		return nil, err
	}
	return f.osScanFS.ReadDir(path)
}

func (f failingScanFS) ReadFile(path string) ([]byte, error) {
	if err := f.readFileErrors[filepath.Clean(path)]; err != nil {
		return nil, err
	}
	return f.osScanFS.ReadFile(path)
}

func (f failingScanFS) WalkDir(path string, fn fs.WalkDirFunc) error {
	if f.walkErr != nil && filepath.Clean(path) == filepath.Clean(f.walkRoot) {
		return fn(f.walkPath, nil, f.walkErr)
	}
	return f.osScanFS.WalkDir(path, fn)
}
