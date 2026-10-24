// Package sync reconciles the DARP asset inventory with the local workspace.
package sync

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/darpbr/darp-cli/internal/project/asset"
	"github.com/darpbr/darp-cli/internal/project/discovery"
	projectinit "github.com/darpbr/darp-cli/internal/project/init"
	"gopkg.in/yaml.v3"
)

type Item struct{ Kind, Identity, Detail string }
type Result struct {
	Added, Removed, Kept, Ignored, Warnings, Failures []Item
	Changed                                           bool
}
type Service struct {
	Root       string
	fileSystem scanFileSystem
}

// Plan contains a deterministic inventory diff and the serialized manifest
// that Apply may persist after checking that the source manifest is unchanged.
type Plan struct {
	Result
	manifest string
	source   []byte
	output   []byte
	missing  bool
}

// BuildPlan reads the workspace and manifest but does not write to the filesystem.
func (s Service) BuildPlan() (Plan, error) {
	root, err := filepath.Abs(s.Root)
	if err != nil {
		return planFailure(Plan{}, ".", err)
	}
	manifest := filepath.Join(root, "darp.yml")
	manifestInfo, statErr := os.Lstat(manifest)
	if statErr == nil && (manifestInfo.Mode()&os.ModeSymlink != 0 || !manifestInfo.Mode().IsRegular()) {
		return planFailure(Plan{}, "darp.yml", errors.New("cannot safely update darp.yml: manifest must be a regular file"))
	}
	if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
		return planFailure(Plan{}, "darp.yml", fmt.Errorf("inspect darp.yml: %w", statErr))
	}
	content, err := os.ReadFile(manifest)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return planFailure(Plan{}, "darp.yml", fmt.Errorf("read darp.yml: %w", err))
	}
	missing := errors.Is(err, os.ErrNotExist)
	if missing {
		content, err = bootstrapManifest(filepath.Base(root))
		if err != nil {
			return planFailure(Plan{}, "darp.yml", fmt.Errorf("build bootstrap manifest: %w", err))
		}
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(content, &doc); err != nil {
		return planFailure(Plan{}, "darp.yml", fmt.Errorf("invalid darp.yml: %w", err))
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		return planFailure(Plan{}, "darp.yml", errors.New("cannot safely update darp.yml: top-level value must be a mapping"))
	}
	if err := rejectDuplicates(doc.Content[0]); err != nil {
		return planFailure(Plan{}, "darp.yml", fmt.Errorf("cannot safely update darp.yml: %w", err))
	}
	rootNode := doc.Content[0]
	if err := validate(rootNode); err != nil {
		return planFailure(Plan{}, "darp.yml", err)
	}
	fileSystem := s.fileSystem
	if fileSystem == nil {
		fileSystem = osScanFS{}
	}
	found, ignored, uncertain := scanWithFS(root, fileSystem)
	if missing {
		// These are created by the standard bootstrap and therefore exist in the
		// in-memory state that an applied sync would observe.
		for _, name := range []string{"documentation", "architecture", "testing", "release"} {
			found.skills[name] = filepath.ToSlash(filepath.Join(".agents/skills", name))
		}
		// init's bootstrap registers the external assets it discovers. Seed that
		// same in-memory manifest so the dry-run report matches applied sync.
		assets := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, a := range found.assets {
			assets.Content = append(assets.Content, assetNode(a))
		}
		set(rootNode, "assets", assets)
	}
	before := yamlNodeString(rootNode)
	result := Result{Ignored: ignored}
	for _, i := range ignored {
		result.Warnings = append(result.Warnings, Item{Kind: "warning", Identity: i.Identity, Detail: i.Detail})
	}
	reconcileSkills(rootNode, found.skills, &result, uncertain)
	reconcileAssets(rootNode, found.assets, &result, uncertain)
	result.Changed = before != yamlNodeString(rootNode)
	sortResult(&result)
	plan := Plan{Result: result, manifest: manifest, source: append([]byte(nil), content...), missing: missing}
	if !result.Changed {
		return plan, nil
	}
	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(&doc); err != nil {
		return planFailure(plan, "darp.yml", fmt.Errorf("encode darp.yml: %w", err))
	}
	if err := encoder.Close(); err != nil {
		return planFailure(plan, "darp.yml", err)
	}
	plan.output = output.Bytes()
	if !missing && bytes.Equal(content, plan.output) {
		plan.Changed = false
	}
	return plan, nil
}

// Apply persists a plan only if its source manifest has not changed since planning.
func (s Service) Apply(plan Plan) (Plan, error) {
	root, err := filepath.Abs(s.Root)
	if err != nil {
		return planFailure(plan, ".", err)
	}
	manifest := filepath.Join(root, "darp.yml")
	if plan.manifest != manifest {
		return planFailure(plan, "darp.yml", errors.New("sync plan belongs to a different project"))
	}
	if plan.missing {
		if _, err := projectinit.NewService(projectinit.NewOSFileSystem()).Initialize(root); err != nil {
			return planFailure(plan, "darp.yml", fmt.Errorf("bootstrap project: %w", err))
		}
		plan, err = s.BuildPlan()
		if err != nil {
			return plan, err
		}
	}
	info, err := os.Lstat(manifest)
	if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		if err == nil {
			err = errors.New("manifest must be a regular file")
		}
		return planFailure(plan, "darp.yml", fmt.Errorf("recheck darp.yml before apply: %w", err))
	}
	current, err := os.ReadFile(manifest)
	if err != nil {
		return planFailure(plan, "darp.yml", fmt.Errorf("recheck darp.yml before apply: %w", err))
	}
	if !bytes.Equal(current, plan.source) {
		return planFailure(plan, "darp.yml", errors.New("darp.yml changed after the sync plan was built"))
	}
	if !plan.Changed {
		return plan, nil
	}
	if err := writeAtomic(manifest, plan.output, info.Mode().Perm(), os.Rename); err != nil {
		return planFailure(plan, "darp.yml", err)
	}
	return plan, nil
}

// Reconcile builds a plan and optionally applies it.
func (s Service) Reconcile(dryRun bool) (Result, error) {
	plan, err := s.BuildPlan()
	if err != nil {
		return plan.Result, err
	}
	if dryRun {
		return plan.Result, nil
	}
	applied, err := s.Apply(plan)
	return applied.Result, err
}

func planFailure(plan Plan, identity string, err error) (Plan, error) {
	plan.Failures = append(plan.Failures, Item{Kind: "failure", Identity: identity, Detail: err.Error()})
	sortItems(plan.Failures)
	return plan, err
}

func sortResult(result *Result) {
	sortItems(result.Added)
	sortItems(result.Removed)
	sortItems(result.Kept)
	sortItems(result.Ignored)
	sortItems(result.Warnings)
	sortItems(result.Failures)
}

func writeAtomic(manifest string, content []byte, mode os.FileMode, rename func(string, string) error) error {
	root := filepath.Dir(manifest)
	tmp, err := os.CreateTemp(root, ".darp-sync-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary darp.yml: %w", err)
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if _, err = tmp.Write(content); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temporary darp.yml: %w", err)
	}
	if err = tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync temporary darp.yml: %w", err)
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	if err = os.Chmod(tmpName, mode); err != nil {
		return fmt.Errorf("set temporary darp.yml permissions: %w", err)
	}
	if err = rename(tmpName, manifest); err != nil {
		return fmt.Errorf("replace darp.yml: %w", err)
	}
	return nil
}

type discoveries struct {
	skills map[string]string
	assets []asset.Asset
}

type scanFileSystem interface {
	Lstat(string) (os.FileInfo, error)
	ReadDir(string) ([]os.DirEntry, error)
	ReadFile(string) ([]byte, error)
	WalkDir(string, fs.WalkDirFunc) error
}

type osScanFS struct{}

func (osScanFS) Lstat(path string) (os.FileInfo, error)     { return os.Lstat(path) }
func (osScanFS) ReadDir(path string) ([]os.DirEntry, error) { return os.ReadDir(path) }
func (osScanFS) ReadFile(path string) ([]byte, error)       { return os.ReadFile(path) }
func (osScanFS) WalkDir(path string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(path, fn)
}

func scanWithFS(root string, fileSystem scanFileSystem) (discoveries, []Item, []string) {
	d := discoveries{skills: map[string]string{}}
	report := discovery.Discover(root)
	ignored := []Item{}
	uncertain := []string{}
	for _, candidate := range report.Found {
		info, err := fileSystem.Lstat(filepath.Join(root, filepath.FromSlash(candidate.Path)))
		if err != nil {
			addIgnored(&ignored, &uncertain, candidate.Path, "read error")
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
			addIgnored(&ignored, &uncertain, candidate.Path, "not a regular file")
			continue
		}
		d.assets = append(d.assets, candidate)
	}
	for _, x := range report.Ignored {
		ignored = append(ignored, Item{Kind: "ignored", Identity: x.Path, Detail: x.Reason})
		uncertain = append(uncertain, x.Path)
	}
	for _, x := range report.Ambiguous {
		ignored = append(ignored, Item{Kind: "ignored", Identity: x.Path, Detail: "ambiguous asset identity"})
		uncertain = append(uncertain, x.Path)
	}
	for _, x := range report.SemanticConflicts {
		ignored = append(ignored, Item{Kind: "ignored", Identity: x.Name, Detail: "semantic asset conflict"})
		for _, a := range x.Assets {
			uncertain = append(uncertain, a.Path)
		}
	}
	// The legacy discovery API intentionally omits I/O failures for init. Sync
	// separately records failures under each approved external asset root so an
	// unreadable tree cannot be mistaken for proof that managed assets vanished.
	for _, dir := range []string{".github", ".claude", ".codex", ".cursor", ".gemini"} {
		abs := filepath.Join(root, dir)
		info, err := fileSystem.Lstat(abs)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				addIgnored(&ignored, &uncertain, dir, "read error")
			}
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			addIgnored(&ignored, &uncertain, dir, "symlink")
			continue
		}
		if !info.IsDir() {
			addIgnored(&ignored, &uncertain, dir, "not a directory")
			continue
		}
		walkErr := fileSystem.WalkDir(abs, func(p string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				rel, e := filepath.Rel(root, p)
				if e == nil {
					addIgnored(&ignored, &uncertain, filepath.ToSlash(rel), "read error")
				}
				return nil
			}
			if entry.IsDir() {
				rel, e := filepath.Rel(root, p)
				if e == nil {
					rel = filepath.ToSlash(rel)
					if len(discovery.Classify(rel)) > 0 {
						addIgnored(&ignored, &uncertain, rel, "candidate is not a regular file")
					}
					candidate := filepath.ToSlash(filepath.Join(rel, "SKILL.md"))
					if isExternalSkillCandidate(candidate) && !hasSkillAncestor(root, rel) && !hasPrefixScope(rel, uncertain) {
						skillInfo, skillErr := fileSystem.Lstat(filepath.Join(p, "SKILL.md"))
						if errors.Is(skillErr, os.ErrNotExist) {
							addIgnored(&ignored, &uncertain, rel, "incomplete skill")
						} else if skillErr != nil {
							addIgnored(&ignored, &uncertain, rel, "read error")
						} else if skillInfo.Mode()&os.ModeSymlink != 0 || !skillInfo.Mode().IsRegular() {
							addIgnored(&ignored, &uncertain, rel, "SKILL.md is not a regular file")
						}
					}
				}
			}
			return nil
		})
		if walkErr != nil {
			addIgnored(&ignored, &uncertain, dir, "read error")
		}
	}
	for _, file := range []string{"AGENTS.md", "CLAUDE.md", "GEMINI.md"} {
		info, err := fileSystem.Lstat(filepath.Join(root, file))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				addIgnored(&ignored, &uncertain, file, "read error")
			}
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			addIgnored(&ignored, &uncertain, file, "symlink")
		} else if !info.Mode().IsRegular() {
			addIgnored(&ignored, &uncertain, file, "not a regular file")
		}
	}
	base := filepath.Join(root, ".agents", "skills")
	agentsRoot := filepath.Join(root, ".agents")
	agentsInfo, agentsErr := fileSystem.Lstat(agentsRoot)
	if agentsErr == nil && agentsInfo.Mode()&os.ModeSymlink != 0 {
		addIgnored(&ignored, &uncertain, ".agents", "symlink")
		sort.Slice(ignored, func(i, j int) bool { return ignored[i].Identity < ignored[j].Identity })
		return d, ignored, uncertain
	}
	if agentsErr == nil && !agentsInfo.IsDir() {
		addIgnored(&ignored, &uncertain, ".agents", "not a directory")
		sort.Slice(ignored, func(i, j int) bool { return ignored[i].Identity < ignored[j].Identity })
		return d, ignored, uncertain
	}
	if agentsErr != nil && !errors.Is(agentsErr, os.ErrNotExist) {
		addIgnored(&ignored, &uncertain, ".agents", "read error")
		sort.Slice(ignored, func(i, j int) bool { return ignored[i].Identity < ignored[j].Identity })
		return d, ignored, uncertain
	}
	info, err := fileSystem.Lstat(base)
	if err == nil && info.Mode()&os.ModeSymlink != 0 {
		ignored = append(ignored, Item{Kind: "ignored", Identity: ".agents/skills", Detail: "symlink"})
		uncertain = append(uncertain, ".agents/skills")
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		ignored = append(ignored, Item{Kind: "ignored", Identity: ".agents/skills", Detail: "read error"})
		uncertain = append(uncertain, ".agents/skills")
	} else if err == nil && !info.IsDir() {
		ignored = append(ignored, Item{Kind: "ignored", Identity: ".agents/skills", Detail: "not a directory"})
		uncertain = append(uncertain, ".agents/skills")
	} else if err == nil {
		entries, e := fileSystem.ReadDir(base)
		if e != nil {
			ignored = append(ignored, Item{Kind: "ignored", Identity: ".agents/skills", Detail: "read error"})
			uncertain = append(uncertain, ".agents/skills")
		} else {
			for _, entry := range entries {
				name := entry.Name()
				rel := filepath.ToSlash(filepath.Join(".agents/skills", name))
				p := filepath.Join(base, name)
				inf, e := fileSystem.Lstat(p)
				if e != nil {
					ignored = append(ignored, Item{Kind: "ignored", Identity: rel, Detail: "stat error"})
					uncertain = append(uncertain, rel)
					continue
				}
				if inf.Mode()&os.ModeSymlink != 0 {
					ignored = append(ignored, Item{Kind: "ignored", Identity: rel, Detail: "symlink"})
					uncertain = append(uncertain, rel)
					continue
				}
				if !inf.IsDir() {
					ignored = append(ignored, Item{Kind: "ignored", Identity: rel, Detail: "not a directory"})
					uncertain = append(uncertain, rel)
					continue
				}
				sk := filepath.Join(p, "SKILL.md")
				si, e := fileSystem.Lstat(sk)
				if e != nil {
					detail := "incomplete skill"
					if !errors.Is(e, os.ErrNotExist) {
						detail = "read error"
					}
					ignored = append(ignored, Item{Kind: "ignored", Identity: rel, Detail: detail})
					uncertain = append(uncertain, rel)
					continue
				}
				if si.Mode()&os.ModeSymlink != 0 || !si.Mode().IsRegular() {
					ignored = append(ignored, Item{Kind: "ignored", Identity: rel, Detail: "SKILL.md is not a regular file"})
					uncertain = append(uncertain, rel)
					continue
				}
				content, e := fileSystem.ReadFile(sk)
				if e != nil {
					ignored = append(ignored, Item{Kind: "ignored", Identity: rel, Detail: "read error"})
					uncertain = append(uncertain, rel)
					continue
				}
				if reason := validateSkillDocument(name, content); reason != "" {
					ignored = append(ignored, Item{Kind: "ignored", Identity: rel, Detail: reason})
					uncertain = append(uncertain, rel)
					continue
				}
				d.skills[name] = rel
			}
		}
	}
	sort.Slice(ignored, func(i, j int) bool {
		if ignored[i].Identity != ignored[j].Identity {
			return ignored[i].Identity < ignored[j].Identity
		}
		return ignored[i].Detail < ignored[j].Detail
	})
	return d, ignored, uncertain
}

func validateSkillDocument(directoryName string, content []byte) string {
	lines := bytes.Split(content, []byte("\n"))
	if len(lines) == 0 || strings.TrimSuffix(string(lines[0]), "\r") != "---" {
		return "missing YAML frontmatter"
	}
	closing := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSuffix(string(lines[i]), "\r") == "---" {
			closing = i
			break
		}
	}
	if closing < 0 {
		return "unterminated YAML frontmatter"
	}
	var document yaml.Node
	if err := yaml.Unmarshal(bytes.Join(lines[1:closing], []byte("\n")), &document); err != nil || len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return "invalid YAML frontmatter"
	}
	metadata := document.Content[0]
	if err := rejectDuplicates(metadata); err != nil {
		return "duplicate frontmatter field"
	}
	name := value(metadata, "name")
	if name == nil || name.Kind != yaml.ScalarNode || name.Tag != "!!str" || !validSkillName(name.Value) || name.Value != directoryName {
		return "invalid name or name does not match directory"
	}
	description := value(metadata, "description")
	if description == nil || description.Kind != yaml.ScalarNode || description.Tag != "!!str" || strings.TrimSpace(description.Value) == "" || utf8.RuneCountInString(description.Value) > 1024 {
		return "description must contain 1 to 1024 characters"
	}
	if compatibility := value(metadata, "compatibility"); compatibility != nil {
		if compatibility.Kind != yaml.ScalarNode || compatibility.Tag != "!!str" || strings.TrimSpace(compatibility.Value) == "" || utf8.RuneCountInString(compatibility.Value) > 500 {
			return "compatibility must contain 1 to 500 characters"
		}
	}
	for _, field := range []string{"license", "allowed-tools"} {
		if fieldValue := value(metadata, field); fieldValue != nil && (fieldValue.Kind != yaml.ScalarNode || fieldValue.Tag != "!!str") {
			return field + " must be a string"
		}
	}
	if extra := value(metadata, "metadata"); extra != nil {
		if extra.Kind != yaml.MappingNode {
			return "metadata must be a string map"
		}
		for i := 0; i+1 < len(extra.Content); i += 2 {
			if extra.Content[i].Kind != yaml.ScalarNode || extra.Content[i].Tag != "!!str" || extra.Content[i+1].Kind != yaml.ScalarNode || extra.Content[i+1].Tag != "!!str" {
				return "metadata must be a string map"
			}
		}
	}
	return ""
}

func validSkillName(name string) bool {
	if len(name) == 0 || len(name) > 64 || name[0] == '-' || name[len(name)-1] == '-' || strings.Contains(name, "--") {
		return false
	}
	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
			return false
		}
	}
	return true
}

func isExternalSkillCandidate(path string) bool {
	for _, family := range []string{".github", ".claude", ".codex", ".cursor", ".gemini"} {
		prefix := filepath.ToSlash(filepath.Join(family, "skills")) + "/"
		if strings.HasPrefix(path, prefix) && strings.TrimPrefix(path, prefix) != "SKILL.md" && len(discovery.Classify(path)) > 0 {
			return true
		}
	}
	return false
}

func hasSkillAncestor(root, directory string) bool {
	for parent := filepath.ToSlash(filepath.Dir(filepath.FromSlash(directory))); parent != "."; parent = filepath.ToSlash(filepath.Dir(filepath.FromSlash(parent))) {
		candidate := filepath.ToSlash(filepath.Join(parent, "SKILL.md"))
		if !isExternalSkillCandidate(candidate) {
			continue
		}
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(parent), "SKILL.md"))
		if err == nil && info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0 {
			return true
		}
	}
	return false
}

func addIgnored(ignored *[]Item, uncertain *[]string, path, reason string) {
	for _, item := range *ignored {
		if item.Identity == path && item.Detail == reason {
			return
		}
	}
	*ignored = append(*ignored, Item{Kind: "ignored", Identity: path, Detail: reason})
	*uncertain = append(*uncertain, path)
}

func validate(n *yaml.Node) error {
	skills := value(n, "skills")
	if skills == nil || skills.Kind != yaml.MappingNode {
		return errors.New("cannot safely update darp.yml: skills must be a mapping")
	}
	for i := 0; i+1 < len(skills.Content); i += 2 {
		key, entry := skills.Content[i], skills.Content[i+1]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
			return errors.New("cannot safely update darp.yml: skill names must be strings")
		}
		if entry.Kind != yaml.ScalarNode || entry.Tag != "!!str" {
			return fmt.Errorf("cannot safely update darp.yml: skill %q path must be a string", key.Value)
		}
	}
	a := value(n, "assets")
	if a != nil && a.Kind != yaml.SequenceNode {
		return errors.New("cannot safely update darp.yml: assets must be a sequence")
	}
	d := value(n, "darp")
	if d != nil && d.Kind != yaml.MappingNode {
		return errors.New("cannot safely update darp.yml: darp must be a mapping")
	}
	s := value(d, "sync")
	if s != nil && s.Kind != yaml.MappingNode {
		return errors.New("cannot safely update darp.yml: darp.sync must be a mapping")
	}
	m := value(s, "managed")
	if m != nil && m.Kind != yaml.MappingNode {
		return errors.New("cannot safely update darp.yml: darp.sync.managed must be a mapping")
	}
	if x := value(m, "skills"); x != nil {
		if x.Kind != yaml.MappingNode {
			return errors.New("cannot safely update darp.yml: managed skills must be a mapping")
		}
		for i := 0; i+1 < len(x.Content); i += 2 {
			if x.Content[i+1].Kind != yaml.ScalarNode {
				return errors.New("cannot safely update darp.yml: managed skill paths must be scalars")
			}
		}
	}
	if x := value(m, "assets"); x != nil {
		if x.Kind != yaml.SequenceNode {
			return errors.New("cannot safely update darp.yml: managed assets must be a sequence")
		}
		for _, entry := range x.Content {
			if assetID(entry) == "" {
				return errors.New("cannot safely update darp.yml: managed asset identities require path, family and type")
			}
		}
	}
	return nil
}

func bootstrapManifest(projectName string) ([]byte, error) {
	root := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	project := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	governance := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	workflows := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	skills := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	set(root, "version", scalar("1.0"))
	set(project, "name", scalar(projectName))
	set(root, "project", project)
	set(governance, "lifecycle", scalar(".darp/lifecycle.md"))
	set(root, "governance", governance)
	set(workflows, "default", scalar("implement"))
	set(root, "workflows", workflows)
	for _, name := range []string{"documentation", "architecture", "testing", "release"} {
		set(skills, name, scalar(filepath.ToSlash(filepath.Join(".agents/skills", name))))
	}
	set(root, "skills", skills)

	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(root); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func value(n *yaml.Node, k string) *yaml.Node {
	if n == nil || n.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == k {
			return n.Content[i+1]
		}
	}
	return nil
}
func set(n *yaml.Node, k string, v *yaml.Node) {
	if x := value(n, k); x != nil {
		*x = *v
		return
	}
	n.Content = append(n.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: k}, v)
}
func scalar(s string) *yaml.Node { return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: s} }
func hasPrefixScope(p string, scopes []string) bool {
	for _, s := range scopes {
		if p == s || strings.HasPrefix(p, strings.TrimSuffix(s, "/")+"/") {
			return true
		}
	}
	return false
}
func reconcileSkills(root *yaml.Node, found map[string]string, r *Result, uncertain []string) bool {
	n := value(root, "skills")
	before := yamlNodeString(n)
	ledger := managed(root, "skills")
	// Adopt legacy entries only when their literal key and canonical path agree.
	for i := 0; i+1 < len(n.Content); i += 2 {
		name, entry := n.Content[i].Value, n.Content[i+1]
		if entry.Kind == yaml.ScalarNode && entry.Tag == "!!str" && canonicalSkillPath(name) == entry.Value {
			ledger[name] = entry.Value
		}
	}
	seen := map[string]bool{}
	keys := make([]string, 0, len(found))
	for k := range found {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, name := range keys {
		path := found[name]
		seen[name] = true
		old := value(n, name)
		if old == nil {
			n.Content = append(n.Content, scalar(name), scalar(path))
			r.Added = append(r.Added, Item{Kind: "skill", Identity: name, Detail: path})
			ledger[name] = path
			continue
		}
		if old.Kind != yaml.ScalarNode || old.Value != path {
			r.Ignored = append(r.Ignored, Item{Kind: "ignored", Identity: name, Detail: "skill path conflict"})
			r.Warnings = append(r.Warnings, Item{Kind: "warning", Identity: name, Detail: "skill path conflict"})
			continue
		}
		r.Kept = append(r.Kept, Item{Kind: "skill", Identity: name, Detail: path})
		ledger[name] = path
	}
	for name, p := range ledger {
		if seen[name] || hasPrefixScope(p, uncertain) {
			continue
		}
		if current := value(n, name); current != nil && current.Kind == yaml.ScalarNode && current.Value == p {
			removeMapKey(n, name)
			r.Removed = append(r.Removed, Item{Kind: "skill", Identity: name, Detail: p})
		} else if current != nil {
			// A user edit relinquishes the old managed value; preserve it.
			r.Ignored = append(r.Ignored, Item{Kind: "ignored", Identity: name, Detail: "managed skill was edited manually"})
			r.Warnings = append(r.Warnings, Item{Kind: "warning", Identity: name, Detail: "managed skill was edited manually"})
		}
		delete(ledger, name)
	}
	writeManaged(root, "skills", ledger)
	return before != yamlNodeString(n)
}

func canonicalSkillPath(name string) string {
	if !validSkillName(name) {
		return ""
	}
	return filepath.ToSlash(filepath.Join(".agents/skills", name))
}

func sortItems(items []Item) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Identity != items[j].Identity {
			return items[i].Identity < items[j].Identity
		}
		return items[i].Detail < items[j].Detail
	})
}
func reconcileAssets(root *yaml.Node, found []asset.Asset, r *Result, uncertain []string) bool {
	n := value(root, "assets")
	before := yamlNodeString(n)
	if n == nil {
		n = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		set(root, "assets", n)
	}
	ledger := managedAssets(root)
	current := map[string]*yaml.Node{}
	for _, x := range n.Content {
		if id := assetID(x); id != "" {
			current[id] = x
			if assetRecognized(x) {
				ledger[id] = true
			}
		}
	}
	seen := map[string]bool{}
	for _, a := range found {
		id := assetKey(a)
		seen[id] = true
		if x := current[id]; x != nil {
			r.Kept = append(r.Kept, Item{Kind: "asset", Identity: assetDisplay(a)})
			continue
		}
		x := assetNode(a)
		n.Content = append(n.Content, x)
		current[id] = x
		ledger[id] = true
		r.Added = append(r.Added, Item{Kind: "asset", Identity: assetDisplay(a)})
	}
	out := n.Content[:0]
	for _, x := range n.Content {
		id := assetID(x)
		if id == "" || seen[id] || !ledger[id] || hasPrefixScope(assetPath(x), uncertain) {
			out = append(out, x)
			continue
		}
		if !assetRecognized(x) {
			out = append(out, x)
			continue
		}
		r.Removed = append(r.Removed, Item{Kind: "asset", Identity: assetDisplay(asset.Asset{Path: nodeValue(x, "path"), Family: asset.Family(nodeValue(x, "family")), Type: asset.Type(nodeValue(x, "type"))})})
		delete(ledger, id)
	}
	n.Content = out
	writeManagedAssets(root, ledger)
	return before != yamlNodeString(n)
}
func managed(root *yaml.Node, kind string) map[string]string {
	out := map[string]string{}
	d := value(root, "darp")
	s := value(d, "sync")
	m := value(s, "managed")
	x := value(m, kind)
	if x != nil && x.Kind == yaml.MappingNode {
		for i := 0; i+1 < len(x.Content); i += 2 {
			out[x.Content[i].Value] = x.Content[i+1].Value
		}
	}
	return out
}
func writeManaged(root *yaml.Node, kind string, m map[string]string) {
	d := value(root, "darp")
	if d == nil {
		d = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		set(root, "darp", d)
	}
	s := value(d, "sync")
	if s == nil {
		s = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		set(d, "sync", s)
	}
	node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		node.Content = append(node.Content, scalar(k), scalar(m[k]))
	}
	managedNode := value(s, "managed")
	if managedNode == nil {
		managedNode = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		set(s, "managed", managedNode)
	}
	set(managedNode, kind, node)
}
func managedAssets(root *yaml.Node) map[string]bool {
	out := map[string]bool{}
	d := value(root, "darp")
	s := value(d, "sync")
	m := value(s, "managed")
	n := value(m, "assets")
	if n != nil && n.Kind == yaml.SequenceNode {
		for _, x := range n.Content {
			if id := assetID(x); id != "" {
				out[id] = true
			}
		}
	}
	return out
}
func writeManagedAssets(root *yaml.Node, m map[string]bool) {
	d := value(root, "darp")
	if d == nil {
		d = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		set(root, "darp", d)
	}
	s := value(d, "sync")
	if s == nil {
		s = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		set(d, "sync", s)
	}
	managedNode := value(s, "managed")
	if managedNode == nil {
		managedNode = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		set(s, "managed", managedNode)
	}
	n := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		parts := strings.Split(k, "\x00")
		if len(parts) != 3 {
			continue
		}
		n.Content = append(n.Content, assetNode(asset.Asset{Path: parts[0], Family: asset.Family(parts[1]), Type: asset.Type(parts[2])}))
	}
	set(managedNode, "assets", n)
}
func assetID(n *yaml.Node) string {
	if n == nil || n.Kind != yaml.MappingNode {
		return ""
	}
	p, f, t := assetPath(n), nodeValue(n, "family"), nodeValue(n, "type")
	if p == "" || f == "" || t == "" {
		return ""
	}
	return p + "\x00" + f + "\x00" + t
}
func assetPath(n *yaml.Node) string { return nodeValue(n, "path") }
func nodeValue(n *yaml.Node, k string) string {
	x := value(n, k)
	if x == nil || x.Kind != yaml.ScalarNode {
		return ""
	}
	return x.Value
}
func assetKey(a asset.Asset) string {
	return a.Path + "\x00" + string(a.Family) + "\x00" + string(a.Type)
}
func assetDisplay(a asset.Asset) string { return fmt.Sprintf("%s (%s/%s)", a.Path, a.Family, a.Type) }
func assetNode(a asset.Asset) *yaml.Node {
	return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: []*yaml.Node{scalar("path"), scalar(a.Path), scalar("family"), scalar(string(a.Family)), scalar("type"), scalar(string(a.Type))}}
}
func assetRecognized(n *yaml.Node) bool {
	p, f, t := assetPath(n), nodeValue(n, "family"), nodeValue(n, "type")
	for _, a := range discovery.Classify(p) {
		if string(a.Family) == f && string(a.Type) == t {
			return true
		}
	}
	return false
}
func removeMapKey(n *yaml.Node, k string) {
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == k {
			n.Content = append(n.Content[:i], n.Content[i+2:]...)
			return
		}
	}
}
func yamlNodeString(n *yaml.Node) string {
	if n == nil {
		return ""
	}
	var b bytes.Buffer
	e := yaml.NewEncoder(&b)
	_ = e.Encode(n)
	_ = e.Close()
	return b.String()
}
func rejectDuplicates(n *yaml.Node) error {
	if n == nil {
		return nil
	}
	if n.Kind == yaml.MappingNode {
		m := map[string]bool{}
		for i := 0; i+1 < len(n.Content); i += 2 {
			k := n.Content[i].Value
			if m[k] {
				return fmt.Errorf("duplicate key %q", k)
			}
			m[k] = true
			if err := rejectDuplicates(n.Content[i+1]); err != nil {
				return err
			}
		}
	} else {
		for _, c := range n.Content {
			if err := rejectDuplicates(c); err != nil {
				return err
			}
		}
	}
	return nil
}

var _ fs.FileMode
