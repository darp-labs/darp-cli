package init

import (
	"embed"
	"fmt"
)

//go:embed assets/.darp/lifecycle.md assets/.darp/governance/quality-gates.md assets/.darp/workflows/implement.yaml assets/.agents/skills/documentation/SKILL.md assets/.agents/skills/architecture/SKILL.md assets/.agents/skills/testing/SKILL.md assets/.agents/skills/release/SKILL.md
var embeddedAssets embed.FS

var assetFiles = loadAssetFiles()

var assetPaths = []struct{ path, target string }{
	{"assets/.darp/lifecycle.md", ".darp/lifecycle.md"},
	{"assets/.darp/governance/quality-gates.md", ".darp/governance/quality-gates.md"},
	{"assets/.darp/workflows/implement.yaml", ".darp/workflows/implement.yaml"},
	{"assets/.agents/skills/documentation/SKILL.md", ".agents/skills/documentation/SKILL.md"},
	{"assets/.agents/skills/architecture/SKILL.md", ".agents/skills/architecture/SKILL.md"},
	{"assets/.agents/skills/testing/SKILL.md", ".agents/skills/testing/SKILL.md"},
	{"assets/.agents/skills/release/SKILL.md", ".agents/skills/release/SKILL.md"},
}

var historicalPlaceholders = map[string]string{
	".darp/governance/quality-gates.md":     "# Quality Gates\n",
	".agents/skills/documentation/SKILL.md": "---\nname: documentation\ndescription: Keep project documentation aligned with the implementation.\n---\n\n# Documentation Skill\n",
	".darp/workflows/implement.yaml":        "name: implement\n\nsteps:\n  - documentation\n",
}

func loadAssetFiles() map[string]string {
	assets := make(map[string]string, len(assetPaths))
	for _, asset := range assetPaths {
		path := asset.path
		content, err := embeddedAssets.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("read embedded init asset %s: %v", path, err))
		}
		assets[asset.target] = string(content)
	}
	return assets
}
