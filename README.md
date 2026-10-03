# DARP CLI

> **Developer AI Resource Platform**

DARP is an open platform for composing, governing, distributing, and evaluating reusable capabilities for agentic software engineering.

DARP provides a vendor-neutral control layer for AI engineering environments.

It integrates open standards such as:

* Agent Skills
* Agent Plugins
* MCP

while providing DARP-level abstractions for:

* Harness
* Registry
* Policy
* Compatibility
* Evaluation

The goal is to make agentic engineering capabilities as portable, versionable, reproducible, and governable as software dependencies are today.

---

## Vision

DARP aims to become an open platform for agentic software engineering capabilities.

Instead of creating another proprietary ecosystem of prompts, Skills, plugins or tools, DARP is designed to compose and govern capabilities from multiple ecosystems.

Conceptually:

```text
Agent Plugin
    ↓
Portable capability

DARP Harness
    ↓
Project-level composition and governance

Agent Client
    ↓
Execution
```

DARP is designed to work with different coding agents and AI providers.

Examples include:

* OpenAI Codex
* GitHub Copilot
* Claude Code
* Cursor
* Gemini-based agents
* future compatible agents

DARP does not replace these products.

---

## DARP Ecosystem

```text
                         DARP
                          │
              ┌───────────┴───────────┐
              │                       │
          Registry                 CLI / SDK
              │                       │
              └───────────┬───────────┘
                          │
                       Harness
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
    Plugins             Assets            Policy
        │                 │                 │
   Agent Plugins       Skills/etc.       Security
        │
   ┌────┴────┐
   │         │
 Skills     MCP
        │
        ▼
┌───────────────────────────────────────────┐
│            Agent Clients                  │
│                                           │
│ Codex │ Copilot │ Claude │ Cursor │ etc. │
└───────────────────────────────────────────┘
```

---

## Core Concepts

### Asset

A reusable AI-oriented artifact or capability.

Examples include:

* Prompts
* Instructions
* Skills
* Workflows
* Templates
* Policies
* Context packages
* MCP-related resources

---

### Agent Plugin

DARP understands interoperable Agent Plugins.

DARP may:

* discover;
* inspect;
* validate;
* catalog;
* install;
* version;
* verify;
* evaluate;
* compose plugins.

DARP does not modify or redefine the Agent Plugin format.

---

### Agent Skill

DARP consumes and manages interoperable Agent Skills.

DARP does not create a competing Skill format.

---

### MCP

DARP integrates with MCP as the tool and contextual interoperability layer.

DARP may discover, catalog, validate and reference MCP resources.

DARP does not replace MCP.

---

### Harness

The DARP Harness is a declarative project-level composition and governance contract.

The canonical manifest is:

```text
.darp/harness.yaml
```

A Harness may reference:

* Agent Plugins
* Skills
* MCP
* Assets
* Providers
* Policies
* Workflows
* Evaluation Contracts
* Quality Gates
* Compatibility requirements

The Harness is a DARP source of truth.

Agent-specific configuration may be generated or projected from the Harness through future adapters.

---

## Important Boundary

The existence of:

```text
.darp/harness.yaml
```

does not mean that every coding agent automatically understands it.

Different agents may require different native configuration mechanisms.

DARP therefore follows an adapter/projection architecture:

```text
DARP Harness
     │
     ├── AGENTS.md
     ├── Copilot configuration
     ├── Claude configuration
     └── other agent projections
```

These projections should not become independent sources of truth.

---

## Current Commands

```bash
darp init
darp doctor
darp --help
darp --version
```

The current implementation provides the DARP project baseline, governance contracts, asset discovery and Harness manifest validation.

---

## Current Scope

The project currently focuses on:

* DARP project initialization;
* repository diagnosis;
* governance;
* AI asset discovery;
* Harness contracts;
* deterministic validation;
* interoperability foundations.

---

## Non Goals

DARP is not intended to:

* replace coding agents;
* replace IDEs;
* replace LLM providers;
* replace Agent Skills;
* replace Agent Plugins;
* replace MCP;
* become an LLM runtime;
* automatically execute discovered capabilities.

Execution requires an explicit future runtime or integration.

---

## Development Methodology

DARP uses Specification-Driven Development.

Every feature follows:

```text
Vision
→ Architecture
→ ADR
→ Specification
→ Plan
→ Tasks
→ Implementation
→ Tests
→ Review
→ Release
```

No feature should be implemented before its specification is approved.

---

## Project Structure

```text
docs/
    PROJECT_CONTEXT.md
    ROADMAP.md
    AGENTIC_HARNESS.md
    ADR/

.spec/
    constitution.md
    specs/

.darp/
    lifecycle.md
    governance/
    harness.yaml       # project Harness, when initialized

.agents/
    skills/

cmd/
internal/
pkg/
test/
```

---

## Principles

* Open standards first
* Provider agnostic
* Agent agnostic
* Interoperability over replacement
* Specification-Driven Development
* AI-first development
* Deterministic behavior
* Reproducibility
* Declarative composition
* Explicit activation
* Least privilege
* Evaluability
* Extensibility

---

## Project Status

🚧 Early development

Current milestone:

* Repository foundation
* Development methodology
* Governance
* Asset discovery
* Harness Manifest foundation
* Agent interoperability architecture

The DARP ecosystem is being developed incrementally.
