# DARP Project Context

## Project

**DARP — Developer AI Resource Platform**

---

## Vision

DARP is an open platform for composing, governing, distributing, and evaluating reusable capabilities for agentic software engineering.

DARP provides a vendor-neutral control layer for AI engineering environments.

The platform integrates open standards such as:

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

## Strategic Positioning

DARP is not a coding agent.

DARP is not an LLM provider.

DARP is not an IDE.

DARP is not an MCP replacement.

DARP is not an Agent Plugin replacement.

DARP is not an Agent Skill replacement.

DARP is the layer that composes and governs these capabilities within an engineering environment.

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

---

## Long-Term Goals

* Discover agentic capabilities
* Install reusable capabilities
* Version capabilities
* Compose capabilities into Harnesses
* Govern capabilities
* Validate compatibility
* Evaluate capabilities
* Distribute capabilities
* Support Agent Skills
* Support Agent Plugins
* Support MCP
* Support multiple AI providers
* Support multiple coding agents
* Provide a DARP Registry
* Provide a DARP SDK
* Support future agent runtimes without coupling the CLI to runtime execution

---

## Core Concepts

### Asset

A reusable AI-oriented artifact or capability.

Examples may include:

* Prompt
* Instruction
* Skill
* Workflow
* Template
* Policy
* Context package
* MCP-related resource

The exact Asset taxonomy is defined progressively by future specifications.

---

### Agent Skill

A reusable agent capability following an interoperable Skill format.

DARP consumes, validates, catalogs, installs and composes Skills.

DARP does not define a competing Skill format.

---

### Agent Plugin

A distributable package following an interoperable Agent Plugin standard.

DARP can:

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

### MCP

An interoperability protocol for tools and contextual capabilities.

DARP can discover, catalog, validate and reference MCP resources.

DARP does not redefine MCP.

---

### Harness

A declarative project-level composition and governance contract.

A Harness defines how capabilities are selected, composed and governed for an engineering environment.

A Harness may reference:

* Agent Plugins
* Skills
* MCP
* local Assets
* Providers
* Policies
* Workflows
* Evaluation Contracts
* Quality Gates
* Compatibility requirements

The canonical Harness manifest is:

```text
.darp/harness.yaml
```

---

### Registry

A future central or federated catalog for DARP resources.

Potential resources include:

* Assets
* Skills
* Agent Plugins
* MCP resources
* Harnesses
* Policies
* Templates

Registry architecture is intentionally deferred.

---

### Provider

A model or AI provider used by an agent or runtime.

Examples:

* OpenAI
* Anthropic
* Google
* DeepSeek
* OpenRouter
* Local model runtimes

DARP remains provider agnostic.

---

### Agent Client

A coding agent responsible for executing development work.

Examples:

* Codex
* GitHub Copilot
* Claude Code
* Cursor
* Gemini-based agents

DARP does not replace these clients.

---

### Adapter / Projection

A mechanism that makes DARP Harness information consumable by an agent-specific environment.

Conceptually:

```text
DARP Harness
     │
     ├── AGENTS.md
     ├── Copilot configuration
     ├── Claude configuration
     └── other agent projections
```

Adapters are not independent sources of truth.

---

### Evaluation

A reproducible mechanism for determining whether a capability or Harness produces the expected engineering result.

Evaluation may include:

* tests;
* quality gates;
* security checks;
* expected behavior;
* cost;
* time;
* task completion;
* compatibility.

---

## Source of Truth

The DARP Harness is the canonical DARP-level source of truth for project-level agentic composition and governance.

Agent-specific configuration files may exist because different agents consume different formats.

These files should be treated as projections or adapters whenever possible.

DARP MUST avoid creating multiple independent sources of truth.

---

## Detection vs Activation

DARP distinguishes between discovering a capability and activating it.

Detection means:

> DARP found evidence that a capability exists.

Activation means:

> The project explicitly configured the capability for use.

Discovery MUST NOT automatically imply:

* execution;
* installation;
* network access;
* credential access;
* MCP connection;
* plugin activation.

This distinction is especially important for enterprise environments.

---

## Architecture

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

## Technology

### DARP CLI

* Language: Go
* Version control: Git
* Hosting: GitHub

### Development

* VS Code
* GitHub Copilot
* OpenAI Codex

The architecture remains independent of these tools.

---

## Current Scope

Current implemented capabilities include:

* DARP project initialization
* project diagnosis
* governance contracts
* governance skills
* AI asset discovery
* Harness manifest parsing
* Harness structural validation

The Harness implementation currently remains read-only and declarative.

It does not execute agents, providers, models, MCP servers, tools or workflows.

---

## Current Non Goals

DARP currently does not:

* execute AI models;
* execute coding agents;
* operate as an IDE;
* replace Agent Plugins;
* replace Agent Skills;
* replace MCP;
* provide a production agent runtime;
* automatically activate discovered capabilities;
* provide a full marketplace;
* provide a centralized Registry.

These may be addressed by future projects or specifications.

---

## Architectural Direction

DARP follows these principles:

1. Open standards first.
2. Provider agnostic.
3. Agent agnostic.
4. Interoperability over replacement.
5. Declarative composition.
6. Explicit activation.
7. Deterministic behavior.
8. Reproducibility.
9. Least privilege.
10. Evaluability.
11. Extensibility.
12. Source-of-truth discipline.

---

## Future Ecosystem

```text
                         DARP Ecosystem
                              │
              ┌───────────────┼───────────────┐
              │               │               │
           DARP CLI       DARP SDK       DARP Registry
              │               │               │
              └───────────────┼───────────────┘
                              │
                           Harness
                              │
             ┌────────────────┼────────────────┐
             │                │                │
          Plugins           Assets           Policy
             │                │                │
      Agent Plugins         Skills          Security
             │
        ┌────┴────┐
        │         │
      Skills     MCP
        │
        ▼
   Agent Clients
        │
 ┌──────┼────────┬──────────┐
 ▼      ▼        ▼          ▼
Codex Copilot Claude      Cursor
```

A future DARP Runtime may consume Harness definitions, but runtime execution requires a separate architectural decision.

---

## Development Methodology

DARP follows Specification-Driven Development.

The expected lifecycle is:

```text
Vision
  ↓
Architecture
  ↓
ADR
  ↓
Specification
  ↓
Plan
  ↓
Tasks
  ↓
Implementation
  ↓
Tests
  ↓
Review
  ↓
Release
```

Architectural decisions must be documented before implementation.

---

## Architectural Vocabulary

| Term          | Meaning                                                        |
| ------------- | -------------------------------------------------------------- |
| Asset         | Reusable AI-oriented artifact or capability                    |
| Skill         | Reusable agent capability                                      |
| Agent Plugin  | Portable package of agent capabilities                         |
| MCP           | Tool/context interoperability protocol                         |
| Harness       | Project-level composition and governance contract              |
| Registry      | Future catalog/distribution infrastructure                     |
| Provider      | AI/model provider                                              |
| Agent Client  | Coding agent responsible for execution                         |
| Adapter       | Mechanism for projecting Harness semantics to an agent         |
| Policy        | Governance rules for capabilities and execution                |
| Evaluation    | Reproducible validation of engineering outcomes                |
| Compatibility | Relationship between resources, agents, providers and runtimes |
