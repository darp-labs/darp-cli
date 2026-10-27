# `.github/prompts/bootstrap-repository.prompt.md`

````md
# Bootstrap DARP Repository

## Purpose

Bootstrap the DARP CLI repository.

This prompt MUST be executed only once.

Its sole responsibility is preparing the repository structure for future development.

No business functionality should be implemented.

---

# Context

You are the lead software architect responsible for initializing the repository.

The project vision, architecture and methodology already exist.

Your responsibility is to transform that vision into a well-organized repository.

You are NOT designing the product.

You are preparing the engineering foundation.

---

# Source of Truth

Before making ANY change, read completely:

1. README.md
2. AGENTS.md
3. Confirm access to the private `darp-cli-specs` repository and read its project context and constitution for authorized work.

If the private repository is unavailable, do not recreate or copy internal documents into this public repository. Ask a maintainer for an approved public scope before proceeding.

Before creating or changing documentation, record its classification, owner, author/source, license or attribution, and audience. Treat `docs/` as public and publish there only after the review and approval process in the private publication policy. Never add credentials or restricted material to the repository.

These documents have precedence over this prompt.

If any contradiction exists:

STOP.

Explain the conflict.

Do not make assumptions.

---

# Repository Development Lifecycle (DDL)

This project follows the DARP Development Lifecycle.

Every future feature will follow:

Vision

↓

Constitution

↓

Project Context

↓

Specification

↓

Architecture Decision Record (ADR)

↓

Implementation Plan

↓

Tasks

↓

Implementation

↓

Review

↓

Release

This bootstrap prepares the repository for that lifecycle.

It does NOT execute any feature.

---

# Mission

Prepare the repository.

Nothing more.

Do NOT implement:

- CLI commands
- Providers
- Registries
- Package manager
- MCP support
- Asset loading
- Business logic
- APIs
- Networking
- Configuration loading

No executable feature should exist after this prompt.

---

# Execution Protocol

## Phase 1

Inspect the repository.

Understand the existing documents.

Identify missing repository components.

---

## Phase 2

Produce a short Bootstrap Plan describing:

- what will be created
- why it is needed
- expected repository layout

Proceed automatically unless a major architectural conflict is found.

---

## Phase 3

Bootstrap the repository.

You may create:

directories

placeholder files

README files

templates

tooling configuration

documentation skeletons

workflow skeletons

development configuration

Nothing else.

---

## Phase 4

Validate the repository.

Verify consistency with:

README

AGENTS

PROJECT_CONTEXT

Constitution

---

## Phase 5

Produce a Bootstrap Report summarizing:

Directories created

Templates created

Configuration files created

Documentation created

Recommendations

Future work

---

# Expected Repository Structure

The final repository should resemble:

```text
docs/
    PROJECT_CONTEXT.md
    ROADMAP.md
    architecture/

.spec/
    constitution.md
    templates/
    specs/
        <number>-<short-name>/
            spec.md
            plan.md
            tasks.md
    archive/

.darp/
    governance/
    workflows/
    templates/

.agents/
    skills/

.github/
    workflows/
    prompts/

.codex/

assets/
    copilot/
    codex/

cmd/

internal/

pkg/

test/
```

Create missing directories only when appropriate.

Do not overwrite existing documentation unless necessary.

---

# Templates

Create reusable templates for:

Specification

Architecture Decision Record (ADR)

Implementation Plan

Tasks

Every template should explain:

Purpose

When to use

Expected sections

Completion criteria

Do NOT create example implementations.

---

# Tooling

Prepare the repository for modern Go development.

Examples include:

.editorconfig

.gitignore

Makefile

golangci configuration

GitHub workflow skeletons

Keep tooling lightweight.

Prefer standard Go tooling.

---

# Documentation

Every architectural directory should contain a short README explaining:

Purpose

Contents

How it participates in the development lifecycle

Documentation should be concise.

Avoid unnecessary verbosity.

---

# AI Optimization

Optimize the repository for:

GitHub Copilot

OpenAI Codex

Future AI agents

Repository structure should be:

predictable

discoverable

easy to navigate

easy to index

---

# Architectural Principles

Always preserve:

Low coupling

High cohesion

Deterministic behavior

Explicit dependencies

Provider agnostic architecture

Small reusable components

Incremental evolution

Readable documentation

Avoid:

Premature abstraction

Speculative features

Unnecessary dependencies

Hidden behavior

Magic configuration

---

# Validation Checklist

Before finishing verify:

✓ No CLI implementation exists.

✓ No business logic exists.

✓ No providers exist.

✓ No registry implementation exists.

✓ No MCP implementation exists.

✓ Templates are reusable.

✓ Documentation is internally consistent.

✓ Repository follows the Constitution.

✓ Repository is ready for the first Specification.

---

# Expected Deliverable

The repository should resemble a mature open-source project that has not yet implemented its first feature.

After completion the next activity should be writing the first Specification.

Nothing beyond repository preparation should be implemented.
````
