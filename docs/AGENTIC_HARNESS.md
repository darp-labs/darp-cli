# DARP Agentic Harness

> Arquitetura estratégica do DARP para composição e governança de ambientes de engenharia de software assistida por agentes.

**Status:** Architectural Direction

**Updated:** 2026-09-09

---

# 1. Purpose

The DARP Harness defines how a project composes and governs its agentic engineering environment.

The Harness is not a coding agent.

It is not an LLM.

It is not an Agent Plugin.

It is not an Agent Skill.

It is not MCP.

It is the project-level composition and governance layer that brings these capabilities together.

---

# 2. Strategic Model

The DARP architecture is:

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

The DARP Harness is therefore a higher-level abstraction than an individual Skill, Plugin or MCP server.

---

# 3. Relationship with Open Standards

DARP follows an open-standards-first strategy.

Whenever an established standard already solves a problem, DARP should consume and integrate that standard rather than create a competing format.

Relevant standards include:

* Agent Skills
* Agent Plugins
* MCP

DARP provides orchestration and governance around these standards.

---

# 4. Agent Plugins

Agent Plugins provide a portable packaging mechanism for agent capabilities.

DARP should:

* discover plugins;
* inspect plugins;
* validate plugins;
* catalog plugins;
* install plugins;
* version plugins;
* verify provenance;
* evaluate plugins;
* reference plugins from Harnesses.

DARP MUST NOT modify the Agent Plugin format.

DARP MUST NOT create a competing plugin format.

---

# 5. Agent Skills

Agent Skills represent reusable agent capabilities.

DARP should consume the established Skill format.

A Skill may be:

```text
local
remote
part of a plugin
registered in a repository
```

The Harness references Skills rather than redefining them.

---

# 6. MCP

MCP provides an interoperability layer for tools and contextual capabilities.

A Harness may declare MCP dependencies or requirements.

Conceptually:

```text
Harness
   │
   ├── GitHub MCP
   ├── Documentation MCP
   ├── Database MCP
   └── Security MCP
```

DARP manages the declarative relationship.

The actual MCP execution remains the responsibility of the consuming agent or runtime.

---

# 7. Harness Responsibilities

The Harness may describe:

* capabilities;
* plugins;
* Skills;
* MCP;
* providers;
* policies;
* workflows;
* compatibility;
* quality gates;
* evaluation;
* governance.

The Harness answers:

> **How should this project's agentic environment be composed and governed?**

---

# 8. Harness Source of Truth

The canonical DARP Harness manifest is:

```text
.darp/harness.yaml
```

The Harness is the DARP-level source of truth.

However, coding agents do not necessarily understand the DARP Harness natively.

For this reason DARP will eventually provide adapters/projections.

Conceptually:

```text
                         DARP Harness
                              │
                       Source of Truth
                              │
             ┌────────────────┼────────────────┐
             │                │                │
             ▼                ▼                ▼
         AGENTS.md         Copilot          Claude
         projection       projection       projection
             │                │                │
             ▼                ▼                ▼
           Agent            Agent            Agent
```

Agent-specific files should be treated as projections whenever possible.

---

# 9. Detection vs Activation

DARP must distinguish between:

```text
Detected
```

and:

```text
Activated
```

Finding a plugin, Skill, MCP server or agent configuration does not automatically authorize its use.

For example:

```text
DARP discovers:
  .claude/
  .github/
  AGENTS.md
  mcp.json
```

This means:

```text
Detected = yes
```

It does not mean:

```text
Enabled = yes
```

Activation must be explicit.

This is required for predictable and secure behavior.

---

# 10. Harness Architecture

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

# 11. Agent Clients

DARP remains independent of the agent responsible for execution.

Examples include:

* Codex;
* GitHub Copilot;
* Claude Code;
* Cursor;
* Gemini-based agents;
* future agents.

The DARP Harness describes the environment.

The agent executes the work.

---

# 12. Adapters and Projections

Different agents consume different configuration mechanisms.

DARP therefore needs a future adapter architecture.

Potential outputs include:

```text
AGENTS.md
CLAUDE.md
Copilot instructions
agent-specific configuration
Skill declarations
MCP configuration
```

The adapter architecture must preserve the Harness as the source of truth.

Future specifications must define:

* supported adapters;
* projection rules;
* conflict detection;
* synchronization;
* generated vs user-owned files;
* update behavior;
* precedence;
* compatibility.

---

# 13. Generic Harness

A project without known technology should be able to use a generic Harness.

The generic Harness MUST NOT assume:

* language;
* framework;
* database;
* CI/CD platform;
* coding agent;
* LLM provider.

The generic Harness establishes only the structural contract required for future configuration.

---

# 14. Existing Repositories

For existing repositories, DARP may discover:

```text
source code
documentation
build systems
test systems
CI/CD
agent instructions
Skills
Agent Plugins
MCP
tool configuration
```

Discovery is local and deterministic whenever possible.

Detection does not automatically activate discovered capabilities.

---

# 15. Security

Agentic capabilities must be treated as software supply-chain dependencies.

Potential security dimensions include:

* provenance;
* publisher identity;
* checksums;
* signatures;
* trusted sources;
* permissions;
* network access;
* filesystem access;
* secrets;
* execution capabilities;
* auditability;
* security scanning.

A future DARP trust model should define these mechanisms explicitly.

---

# 16. Evaluation

Evaluation is part of the Harness architecture.

The objective is to measure engineering outcomes rather than merely determine whether a package can be installed.

Potential evaluation dimensions:

```text
Task
Repository fixture
Expected behavior
Tests
Quality gates
Security checks
Cost
Time
Result
```

This may eventually enable comparison between:

```text
Codex + Harness A
Claude + Harness A
Copilot + Harness A
Local Model + Harness A
```

using equivalent tasks and evaluation criteria.

---

# 17. Registry

The DARP Registry is a future component.

Its conceptual responsibility is:

```text
DARP Registry
    │
    ├── Assets
    ├── Agent Plugins
    ├── Skills
    ├── MCP
    ├── Harnesses
    ├── Policies
    └── Metadata
```

Detailed Registry architecture is intentionally deferred.

---

# 18. Runtime

DARP CLI is not initially an agent runtime.

It does not execute:

* models;
* agents;
* plugins;
* MCP servers;
* arbitrary workflows.

A future DARP Runtime may consume Harness definitions.

Such a runtime requires a separate architectural decision.

---

# 19. Design Principles

The Harness follows:

1. Declarative composition
2. Vendor neutrality
3. Agent neutrality
4. Open standards
5. Explicit activation
6. Least privilege
7. Reproducibility
8. Deterministic validation
9. Compatibility awareness
10. Evaluability
11. Extensibility
12. Source-of-truth discipline

---

# 20. Final Model

The DARP Harness should be understood as:

> **The declarative, project-level composition and governance contract for agentic software engineering.**

It composes capabilities.

It does not redefine them.

It governs them.

It does not execute them.

It provides a common project-level abstraction.

It does not replace the agents that consume that environment.
