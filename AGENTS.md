# AGENTS.md

This document defines how AI agents must contribute to this repository.

## Mission

Help build DARP CLI while preserving architectural consistency.

Agents must prioritize quality over speed.

## Project Goal

DARP CLI is a package manager for AI assets.

The CLI manages reusable artifacts instead of source code libraries.

Examples:

- Prompt
- Instruction
- Skill
- Persona
- MCP Server
- Workflow

## Development Methodology

This repository follows Specification-Driven Development (SDD).

Implementation order is always:

1. Vision
2. Constitution
3. Project Context
4. Specification
5. ADR, when required
6. Plan
7. Tasks
8. Implementation
9. Tests
10. Architecture Review
11. Documentation Review
12. Quality Gates
13. Review
14. Release
15. Completed

Never skip steps.

The public runtime lifecycle is defined in `.darp/lifecycle.md`. Internal project context, constitution, specifications, plans, tasks, and ADRs are maintained in the private `darp-cli-specs` repository under `specs/internal/darp-cli/`.

## Specification Layout

Feature planning artifacts must live in the private `darp-cli-specs` repository under `specs/internal/darp-cli/.spec/specs/`.
Create one directory per feature using the format `<number>-<short-name>`:

```text
.spec/specs/002-doctor/
├── spec.md
├── plan.md
├── tasks.md
└── adr/                         # optional, when architecture decisions exist
```

Use the templates in `specs/internal/darp-cli/.spec/templates/` in the private specs repository. Do not create or copy internal specifications, plans, or tasks into this public repository. Contributors without access must obtain an approved public scope from a maintainer before implementing a feature. The `.darp/` directory is reserved for public DARP runtime contracts, while `.agents/skills/` is reserved for skills distributed or discovered by coding agents.

## Documentation Publication

- Classify new or materially changed documents before sharing them. Record owner, author/source, license or attribution, and audience.
- Treat `docs/` as public. Add only content approved under the publication policy in the private `darp-cli-specs` repository.
- Do not copy internal specifications, plans, roadmaps, project context, or ADRs into this repository. For approved publications, keep this public copy canonical and ask a maintainer to update the private publication index.
- Review examples, attachments, filenames, metadata, Git history, and possible external indexing. AI checks are assistance, not approval or proof of rights.
- Never commit credentials, personal data, customer/corporate information, or other restricted content.

## Architecture Principles

- Keep packages cohesive.
- Minimize coupling.
- Prefer composition over inheritance.
- Keep commands independent.
- Avoid global state.
- Favor deterministic behavior.
- Design for extensibility.

## Coding Principles

- Keep functions small.
- Prefer readability over cleverness.
- Document exported APIs.
- Write tests whenever functionality is implemented.
- Never introduce unnecessary abstractions.

## Decision Hierarchy

When conflicts arise, follow this order:

1. Constitution
2. Project Context
3. Approved Specification
4. ADR
5. Plan
6. Tasks

Never override a higher-level document.

## Repository Rules

Do not:

- implement features outside an approved specification;
- create hidden behavior;
- duplicate business logic;
- introduce unnecessary dependencies.

Always:

- update documentation;
- update `CHANGELOG.md` under `Unreleased` for user-visible changes, fixes, breaking changes or deprecations;
- explain assumptions;
- write implementation tasks as unchecked Markdown checkboxes (`- [ ]`);
- mark a task as complete only after its validation passes;
- leave blocked tasks unchecked and record the blocker and evidence;
- preserve backward compatibility whenever possible.

## AI Behaviour

When uncertain:

- ask for clarification;
- do not invent requirements;
- document assumptions.

When implementing:

- prefer incremental changes;
- keep commits focused;
- avoid unrelated refactoring.

## Development commands

- `make build` — build the CLI binary
- `make test` — run Go tests
- `make lint` — run golangci-lint (falls back to go vet)
- `make verify` — fmt + test + lint
- `make run` — run the CLI locally

## Repo layout

- `cmd/` — CLI entrypoints
- `internal/cli/`, `internal/project/` — domain packages
- `pkg/` — public reusable packages
- `docs/` — approved public user guides, references, and specifications
- `.agents/skills/` — agent skills, loaded on demand