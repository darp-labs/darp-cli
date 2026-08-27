// Package discovery scans a project for pre-existing AI assets from supported
// tool families. Discovery is local, deterministic, read-only and never
// executes, moves or overwrites files. It only evaluates the known roots and
// patterns approved in Spec 005.
package discovery

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/darpbr/darp-cli/internal/project/asset"
)

// rule maps a single supported pattern to a family and semantic type.
// glob is slash-separated and relative to the project root. A glob without a
// slash matches a root file directly.
type rule struct {
	family asset.Family
	typ    asset.Type
	glob   string
}

// Classify returns the supported asset identities for a normalized relative
// path. It does not access the filesystem.
func Classify(relativePath string) []asset.Asset {
	relativePath = filepath.ToSlash(filepath.Clean(filepath.FromSlash(relativePath)))
	if relativePath == "." || relativePath == ".." || strings.HasPrefix(relativePath, "../") || filepath.IsAbs(relativePath) {
		return nil
	}

	classified := make([]asset.Asset, 0)
	for _, r := range rules {
		if matchGlob(r.glob, relativePath) {
			classified = append(classified, asset.Asset{Path: relativePath, Family: r.family, Type: r.typ})
		}
	}
	sort.Slice(classified, func(i, j int) bool { return assetLess(classified[i], classified[j]) })
	return classified
}

// rules is the canonical table of patterns approved in Spec 005.
var rules = []rule{
	// GitHub
	{asset.FamilyGitHub, asset.TypeInstruction, ".github/copilot-instructions.md"},
	{asset.FamilyGitHub, asset.TypeInstruction, ".github/instructions/**/*.instructions.md"},
	{asset.FamilyGitHub, asset.TypePrompt, ".github/prompts/**/*.prompt.md"},
	{asset.FamilyGitHub, asset.TypeSkill, ".github/skills/**/SKILL.md"},
	{asset.FamilyGitHub, asset.TypeHook, ".github/hooks/*.json"},
	{asset.FamilyGitHub, asset.TypeAgent, ".github/agents/*.md"},
	{asset.FamilyGitHub, asset.TypeInstruction, "AGENTS.md"},
	{asset.FamilyGitHub, asset.TypeInstruction, "CLAUDE.md"},
	{asset.FamilyGitHub, asset.TypeInstruction, "GEMINI.md"},
	// Claude
	{asset.FamilyClaude, asset.TypeInstruction, "CLAUDE.md"},
	{asset.FamilyClaude, asset.TypeInstruction, ".claude/CLAUDE.md"},
	{asset.FamilyClaude, asset.TypeCommand, ".claude/commands/**/*.md"},
	{asset.FamilyClaude, asset.TypeSkill, ".claude/skills/**/SKILL.md"},
	{asset.FamilyClaude, asset.TypeAgent, ".claude/agents/**/*.md"},
	{asset.FamilyClaude, asset.TypeHook, ".claude/settings.json"},
	// Codex
	{asset.FamilyCodex, asset.TypeInstruction, "AGENTS.md"},
	{asset.FamilyCodex, asset.TypeSkill, ".codex/skills/**/SKILL.md"},
	// Cursor
	{asset.FamilyCursor, asset.TypeRule, ".cursor/rules/**/*.mdc"},
	{asset.FamilyCursor, asset.TypeSkill, ".cursor/skills/**/SKILL.md"},
	{asset.FamilyCursor, asset.TypeHook, ".cursor/hooks.json"},
	{asset.FamilyCursor, asset.TypeInstruction, "AGENTS.md"},
	// Gemini
	{asset.FamilyGemini, asset.TypeInstruction, "GEMINI.md"},
	{asset.FamilyGemini, asset.TypeCommand, ".gemini/commands/**/*.toml"},
	{asset.FamilyGemini, asset.TypeSkill, ".gemini/skills/**/SKILL.md"},
	{asset.FamilyGemini, asset.TypeConfig, ".gemini/settings.json"},
	{asset.FamilyGemini, asset.TypeExtension, ".gemini/extensions/*/gemini-extension.json"},
}

// knownDirs are the conventional family roots scanned by discovery. The v1
// never traverses arbitrary directories outside these roots.
var knownDirs = []string{".github", ".claude", ".codex", ".cursor", ".gemini"}

// Ambiguity groups assets that share a physical path but differ in family or
// type, which the CLI reports without persisting automatically.
type Ambiguity struct {
	Path   string
	Assets []asset.Asset
}

// SemanticConflict groups assets with the same category and normalized name.
type SemanticConflict struct {
	Category string
	Name     string
	Assets   []asset.Asset
}

// IgnoredItem reports a path that discovery skipped and the reason why.
type IgnoredItem struct {
	Path   string
	Reason string
}

// Report is the structured outcome of a discovery pass.
type Report struct {
	// Found are unambiguous assets ready for registration.
	Found []asset.Asset
	// Ambiguous lists physical paths recognized by more than one identity.
	Ambiguous []Ambiguity
	// SemanticConflicts lists candidates with the same semantic identity.
	SemanticConflicts []SemanticConflict
	// Ignored lists symlinks and other skipped paths.
	Ignored []IgnoredItem
}

// Empty reports whether the report contains no findings at all.
func (r Report) Empty() bool {
	return len(r.Found) == 0 && len(r.Ambiguous) == 0 && len(r.SemanticConflicts) == 0 && len(r.Ignored) == 0
}

// Discover scans root for supported assets. It is read-only and deterministic.
func Discover(root string) Report {
	report := Report{}
	candidates := map[string][]asset.Asset{}

	for _, r := range rules {
		if !isRootFileRule(r) {
			continue
		}
		abs := filepath.Join(root, filepath.FromSlash(r.glob))
		info, err := os.Lstat(abs)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			report.Ignored = append(report.Ignored, IgnoredItem{Path: r.glob, Reason: "symlink"})
			continue
		}
		if !info.Mode().IsRegular() {
			continue
		}
		candidates[r.glob] = append(candidates[r.glob], asset.Asset{Path: r.glob, Family: r.family, Type: r.typ})
	}

	for _, dir := range knownDirs {
		abs := filepath.Join(root, filepath.FromSlash(dir))
		info, err := os.Lstat(abs)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			report.Ignored = append(report.Ignored, IgnoredItem{Path: dir, Reason: "symlink"})
			continue
		}
		if !info.IsDir() {
			continue
		}
		_ = filepath.WalkDir(abs, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.Type()&os.ModeSymlink != 0 {
				rel, relErr := filepath.Rel(root, p)
				if relErr == nil {
					report.Ignored = append(report.Ignored, IgnoredItem{Path: filepath.ToSlash(rel), Reason: "symlink"})
				}
				return nil
			}
			if d.IsDir() {
				return nil
			}
			rel, relErr := filepath.Rel(root, p)
			if relErr != nil {
				return nil
			}
			slash := filepath.ToSlash(rel)
			candidates[slash] = append(candidates[slash], Classify(slash)...)
			return nil
		})
	}

	for _, assets := range candidates {
		if len(assets) > 1 {
			report.Ambiguous = append(report.Ambiguous, Ambiguity{Path: assets[0].Path, Assets: assets})
			continue
		}
		report.Found = append(report.Found, assets[0])
	}

	semanticGroups := make(map[string][]asset.Asset)
	for _, candidate := range report.Found {
		key := semanticKey(candidate)
		semanticGroups[key] = append(semanticGroups[key], candidate)
	}
	filtered := report.Found[:0]
	for _, candidate := range report.Found {
		key := semanticKey(candidate)
		if len(semanticGroups[key]) > 1 {
			continue
		}
		filtered = append(filtered, candidate)
	}
	report.Found = filtered
	for key, assets := range semanticGroups {
		if len(assets) < 2 {
			continue
		}
		category, name := splitSemanticKey(key)
		report.SemanticConflicts = append(report.SemanticConflicts, SemanticConflict{Category: category, Name: name, Assets: assets})
	}

	sortReport(&report)
	return report
}

func isRootFileRule(r rule) bool {
	return !strings.Contains(r.glob, "/")
}

func sortReport(r *Report) {
	sort.Slice(r.Found, func(i, j int) bool { return assetLess(r.Found[i], r.Found[j]) })
	sort.Slice(r.Ambiguous, func(i, j int) bool { return r.Ambiguous[i].Path < r.Ambiguous[j].Path })
	for i := range r.Ambiguous {
		sort.Slice(r.Ambiguous[i].Assets, func(a, b int) bool {
			return assetLess(r.Ambiguous[i].Assets[a], r.Ambiguous[i].Assets[b])
		})
	}
	sort.Slice(r.SemanticConflicts, func(i, j int) bool {
		if r.SemanticConflicts[i].Category != r.SemanticConflicts[j].Category {
			return r.SemanticConflicts[i].Category < r.SemanticConflicts[j].Category
		}
		return r.SemanticConflicts[i].Name < r.SemanticConflicts[j].Name
	})
	for i := range r.SemanticConflicts {
		sort.Slice(r.SemanticConflicts[i].Assets, func(a, b int) bool {
			return assetLess(r.SemanticConflicts[i].Assets[a], r.SemanticConflicts[i].Assets[b])
		})
	}
	sort.Slice(r.Ignored, func(i, j int) bool { return r.Ignored[i].Path < r.Ignored[j].Path })
}

func assetLess(a, b asset.Asset) bool {
	if a.Path != b.Path {
		return a.Path < b.Path
	}
	if a.Family != b.Family {
		return a.Family < b.Family
	}
	return a.Type < b.Type
}

func semanticKey(candidate asset.Asset) string {
	return string(candidate.Type) + "\x00" + semanticName(candidate)
}

func splitSemanticKey(key string) (string, string) {
	parts := strings.SplitN(key, "\x00", 2)
	return parts[0], parts[1]
}

func semanticName(candidate asset.Asset) string {
	base := path.Base(candidate.Path)
	switch candidate.Type {
	case asset.TypeSkill, asset.TypeExtension:
		base = path.Base(path.Dir(candidate.Path))
	}
	for _, suffix := range []string{".instructions.md", ".prompt.md", ".mdc", ".toml", ".json", ".md"} {
		if strings.HasSuffix(base, suffix) {
			base = strings.TrimSuffix(base, suffix)
			break
		}
	}
	return strings.ToLower(base)
}

// matchGlob reports whether the slash-separated rel path matches glob. It
// supports the "**" segment, which matches zero or more path segments.
func matchGlob(glob, rel string) bool {
	return matchSegments(strings.Split(glob, "/"), strings.Split(rel, "/"))
}

func matchSegments(glob, rel []string) bool {
	if len(glob) == 0 {
		return len(rel) == 0
	}
	switch glob[0] {
	case "**":
		if matchSegments(glob[1:], rel) {
			return true
		}
		if len(rel) > 0 {
			return matchSegments(glob, rel[1:])
		}
		return false
	default:
		if len(rel) == 0 {
			return false
		}
		ok, err := path.Match(glob[0], rel[0])
		if err != nil || !ok {
			return false
		}
		return matchSegments(glob[1:], rel[1:])
	}
}
