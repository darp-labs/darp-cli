// Package asset defines the shared domain model for AI assets recognized by
// DARP. Families and types are the contract approved in Spec 005 and ADR 001.
package asset

// Family identifies the tool family that owns an asset convention.
type Family string

const (
	FamilyGitHub Family = "github"
	FamilyClaude Family = "claude"
	FamilyCodex  Family = "codex"
	FamilyCursor Family = "cursor"
	FamilyGemini Family = "gemini"
)

// Type identifies the semantic type of an asset.
type Type string

const (
	TypeInstruction Type = "instruction"
	TypePrompt      Type = "prompt"
	TypeSkill       Type = "skill"
	TypeHook        Type = "hook"
	TypeAgent       Type = "agent"
	TypeCommand     Type = "command"
	TypeRule        Type = "rule"
	TypeConfig      Type = "config"
	TypeExtension   Type = "extension"
)

// Asset is a recognized asset candidate.
type Asset struct {
	// Path is the slash-separated path relative to the project root.
	Path string
	// Family is the tool family that owns the asset convention.
	Family Family
	// Type is the semantic type of the asset.
	Type Type
}

var validFamilies = map[string]bool{
	string(FamilyGitHub): true,
	string(FamilyClaude): true,
	string(FamilyCodex):  true,
	string(FamilyCursor): true,
	string(FamilyGemini): true,
}

var validTypes = map[string]bool{
	string(TypeInstruction): true,
	string(TypePrompt):      true,
	string(TypeSkill):       true,
	string(TypeHook):        true,
	string(TypeAgent):       true,
	string(TypeCommand):     true,
	string(TypeRule):        true,
	string(TypeConfig):      true,
	string(TypeExtension):   true,
}

// ValidFamily reports whether value is a supported family.
func ValidFamily(value string) bool {
	return validFamilies[value]
}

// ValidType reports whether value is a supported asset type.
func ValidType(value string) bool {
	return validTypes[value]
}
