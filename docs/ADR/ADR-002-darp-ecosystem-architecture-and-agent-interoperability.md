# ADR-002 — DARP Ecosystem Architecture and Agent Interoperability

**Status:** Accepted
**Date:** 2026-09-09
**Decision Type:** Architecture / Ecosystem
**Scope:** DARP Ecosystem / DARP CLI / Harness

---

## 1. Context

The DARP project initially defined itself primarily as a package manager for reusable AI assets.

As the ecosystem of AI-assisted software engineering evolves, several open and vendor-supported standards are emerging around reusable agent capabilities, including Agent Skills, Agent Plugins and MCP.

At the same time, coding agents such as Codex, GitHub Copilot, Claude Code, Cursor and other agent clients are developing their own mechanisms for instructions, skills, tools, plugins and workflows.

This creates an architectural risk for DARP.

DARP could attempt to define proprietary formats for capabilities already covered by emerging standards, resulting in fragmentation and reduced interoperability.

Alternatively, DARP can position itself as a higher-level, vendor-neutral control layer that composes and governs these capabilities without replacing the standards or agents that consume them.

This ADR defines the architectural boundaries between:

* DARP;
* AI Assets;
* Agent Skills;
* Agent Plugins;
* MCP;
* Harness;
* Registry;
* Agent Clients;
* Runtime;
* Policy;
* Evaluation.

---

# 2. Decision

DARP will evolve from a primarily AI asset package-management concept into an **open platform for composing, governing, distributing and evaluating reusable capabilities for agentic software engineering**.

DARP will prioritize interoperability with established open standards and will avoid defining competing formats when an appropriate open standard already exists.

The DARP Harness will become the central declarative abstraction for describing how capabilities are composed and governed within a project.

DARP will not replace coding agents, AI providers, MCP, Agent Skills or Agent Plugins.

---

# 3. New DARP Positioning

The official conceptual positioning of DARP becomes:

> **DARP is an open platform for composing, governing, distributing, and evaluating reusable capabilities for agentic software engineering.**

DARP provides a vendor-neutral control layer for AI engineering environments, integrating open standards such as Agent Skills, Agent Plugins, and MCP while providing its own abstractions for Harness, Registry, Policy, Compatibility, and Evaluation.

The goal is to make agentic engineering capabilities as portable, versionable, reproducible, and governable as software dependencies are today.

---

# 4. Architectural Principles

The following principles are established by this ADR.

## 4.1 Open Standards First

DARP MUST prefer established open standards over proprietary DARP-specific formats when those standards adequately address the problem.

DARP MUST NOT create a competing format merely to provide functionality already covered by an interoperable standard.

---

## 4.2 Vendor Neutrality

DARP MUST remain independent of any individual coding agent, LLM provider, IDE or model runtime.

Examples include:

* OpenAI;
* Anthropic;
* Google;
* GitHub;
* Cursor;
* AWS;
* local model runtimes;
* future providers and agents.

Provider- or agent-specific integrations MUST remain isolated behind compatibility or adapter mechanisms.

---

## 4.3 Interoperability Over Replacement

DARP exists to compose and govern the agentic ecosystem, not to replace its individual components.

DARP MUST NOT attempt to replace:

* Agent Skills;
* Agent Plugins;
* MCP;
* coding agents;
* IDEs;
* LLM providers;
* model runtimes.

---

# 5. Ecosystem Model

The DARP ecosystem is conceptually defined as:

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

The Harness is the DARP abstraction that describes the desired project-level composition and governance.

The actual execution remains the responsibility of the compatible agent or runtime.

---

# 6. Conceptual Boundaries

## 6.1 Asset

An Asset is a reusable AI-oriented artifact or capability.

Examples may include:

* prompts;
* instructions;
* skills;
* workflows;
* templates;
* policies;
* context packages;
* MCP-related resources.

The exact taxonomy of Assets remains subject to future specifications.

---

## 6.2 Agent Skill

An Agent Skill is a capability following an interoperable Skill format supported by compatible agent clients.

DARP MUST support and consume established Agent Skill formats rather than creating a competing Skill format.

DARP MAY:

* discover Skills;
* validate Skills;
* catalog Skills;
* version Skills;
* install Skills;
* evaluate Skills;
* declare Skill dependencies;
* reference Skills from a Harness.

DARP MUST NOT redefine the underlying Agent Skill format.

---

## 6.3 Agent Plugin

An Agent Plugin is a distributable extension following the Agent Plugins standard or another explicitly supported compatible format.

DARP MUST treat Agent Plugins as external interoperable packages.

DARP MAY:

* discover plugins;
* inspect plugins;
* validate plugin structure;
* catalog plugins;
* resolve versions;
* install plugins;
* verify provenance;
* evaluate plugins;
* associate plugins with Harness definitions.

DARP MUST NOT modify the Agent Plugin format.

DARP MUST NOT require a DARP-specific plugin format when the standard Agent Plugin format is sufficient.

---

## 6.4 MCP

MCP is the interoperability layer for exposing tools and contextual capabilities to compatible agents.

DARP MUST treat MCP as an external interoperability standard.

DARP MAY:

* discover MCP servers;
* catalog MCP servers;
* validate declarations;
* reference MCP servers from Harness definitions;
* declare compatibility;
* declare policies;
* evaluate MCP-related capabilities.

DARP MUST NOT redefine the MCP protocol.

Operational execution, credentials and authorization remain outside the initial DARP CLI responsibility unless explicitly specified by a future ADR/specification.

---

# 7. Harness

The DARP Harness is a declarative project-level composition and governance contract.

Its primary question is:

> **How should this project's agentic environment be composed and governed?**

It is not a replacement for individual capabilities.

The conceptual model is:

```text
Agent Plugin
    │
    │ provides capabilities
    ▼
DARP Harness
    │
    ├── selects/composes capabilities
    ├── defines project policy
    ├── defines workflows
    ├── defines compatibility
    ├── defines evaluation requirements
    └── defines governance
    │
    ▼
Agent Client / Runtime
```

The Harness MUST compose existing capabilities instead of redefining them.

A Harness MAY reference:

* Agent Plugins;
* Agent Skills;
* MCP servers;
* local Assets;
* provider profiles;
* workflows;
* policies;
* evaluation contracts;
* quality gates.

---

# 8. Harness Source of Truth

The DARP Harness is the canonical DARP-level source of truth for the project's agentic environment.

The canonical manifest path established by the Harness Manifest specification is:

```text
.darp/harness.yaml
```

However, the existence of this file MUST NOT be interpreted as meaning that every coding agent automatically understands or executes its contents.

Agent clients may only recognize their own native configuration mechanisms and supported interoperability standards.

Therefore DARP MUST support an adapter/projection model.

Conceptually:

```text
                         .darp/harness.yaml
                                  │
                           DARP Source of Truth
                                  │
                 ┌────────────────┼────────────────┐
                 │                │                │
                 ▼                ▼                ▼
             AGENTS.md       Copilot config    Claude config
             projection       projection        projection
                 │                │                │
                 ▼                ▼                ▼
              Agents           Agents           Agents
```

These projections MUST NOT become independent architectural sources of truth.

The exact adapter and projection mechanisms are deferred to future specifications.

---

# 9. Detection vs Activation

DARP MUST distinguish between:

```text
Detected
```

and:

```text
Enabled / Activated
```

For example, discovering:

```text
.claude/
.github/
AGENTS.md
mcp.json
```

does not automatically mean that DARP is authorized to activate, install, execute or connect those resources.

Discovery is evidence.

Activation is an explicit configuration decision.

This distinction is especially important for enterprise security and supply-chain trust.

---

# 10. Agent Client Boundary

Coding agents remain responsible for executing development tasks.

Examples include:

* Codex;
* GitHub Copilot;
* Claude Code;
* Cursor;
* Gemini-based coding agents;
* future compatible agents.

DARP does not become a coding agent merely because it manages the environment consumed by one.

DARP provides declarative project context and governance.

The agent client remains responsible for agent execution.

---

# 11. Runtime Boundary

DARP CLI MUST NOT initially become an LLM or agent runtime.

The following remain outside the initial DARP CLI:

* model inference;
* agent execution loops;
* autonomous task execution;
* model routing;
* runtime orchestration;
* execution sandboxing.

A future DARP Runtime MAY consume Harness definitions.

Such a runtime MUST be defined by a separate specification and architectural decision.

---

# 12. Registry

The DARP Registry is a future ecosystem component.

Its conceptual responsibility is:

```text
DARP Registry
    │
    ├── Assets
    ├── Agent Plugins
    ├── Skills
    ├── MCP resources
    ├── Harnesses
    ├── Policies
    └── Metadata
```

Potential Registry capabilities include:

* discovery;
* versioning;
* dependency resolution;
* distribution;
* provenance;
* trust metadata;
* compatibility metadata;
* evaluation metadata.

Detailed Registry architecture is intentionally deferred.

This ADR establishes only the architectural direction.

---

# 13. Policy and Governance

DARP will eventually provide project-level governance for agentic capabilities.

Potential policy dimensions include:

* allowed plugins;
* allowed MCP servers;
* allowed providers;
* filesystem access;
* network access;
* secrets;
* execution permissions;
* trust levels;
* required evaluations;
* required quality gates.

The policy system MUST follow least-privilege principles.

Policy definition and policy enforcement are separate concerns.

DARP CLI may validate declared policy.

Actual runtime enforcement belongs to the consuming agent/runtime unless a future DARP Runtime is explicitly defined.

---

# 14. Evaluation

Evaluation is a first-class architectural concern.

DARP should eventually be able to describe not only:

> "What capabilities are available?"

but also:

> "Does this capability or Harness produce the expected engineering result?"

Evaluation may include:

* reproducible tasks;
* repository fixtures;
* tests;
* quality gates;
* security checks;
* cost constraints;
* execution time;
* expected behavior;
* result scores.

The detailed Evaluation Contract remains a future specification.

---

# 15. Compatibility

DARP will maintain a compatibility model capable of representing relationships between:

```text
Harness
Plugin
Skill
MCP
Agent
Provider
Runtime
Stack
```

Compatibility MUST be explicit whenever possible.

DARP MUST NOT assume that an asset compatible with one agent is automatically compatible with another.

Compatibility metadata may eventually be used by:

* Registry;
* CLI;
* Harness validation;
* dependency resolution;
* installation;
* evaluation.

---

# 16. Security and Supply Chain

Agent capabilities must be treated as software supply-chain dependencies with behavioral implications.

DARP should eventually support concepts such as:

* publisher identity;
* provenance;
* checksums;
* signatures;
* trusted sources;
* lockfiles;
* permissions;
* security scanning;
* auditability;
* reproducible installation;
* evaluation evidence.

The detailed trust and security model is deferred to future architecture decisions.

---

# 17. DARP CLI Responsibilities

The DARP CLI is responsible for local lifecycle management of DARP resources.

Potential responsibilities include:

```text
discover
inspect
validate
install
catalog
compose
configure
evaluate
```

The CLI MUST remain deterministic and explicit.

The CLI MUST NOT silently activate discovered capabilities.

The CLI MUST NOT execute arbitrary plugin, Skill or MCP behavior merely because it discovers or validates those resources.

---

# 18. Future Adapter Architecture

DARP should eventually support adapters for agent-specific environments.

Conceptually:

```text
DARP Harness
     │
     ├── Generic representation
     │
     ├── Agent Skills
     ├── Agent Plugins
     ├── MCP
     │
     └── Agent Adapters
             │
       ┌─────┼──────┬─────────┐
       ▼     ▼      ▼         ▼
     Codex Copilot Claude    Cursor
```

Adapters may generate or synchronize:

* instructions;
* agent-specific configuration;
* references;
* supported Skills;
* plugin declarations;
* MCP configuration;
* other agent-native metadata.

The exact behavior of adapters requires a future specification.

---

# 19. Consequences

## Positive Consequences

### Interoperability

DARP can benefit from ecosystem standards instead of competing with them.

### Vendor neutrality

Projects can use multiple agents without making DARP dependent on a single provider.

### Enterprise adoption

Organizations can have a project-level governance layer independent of the coding agent selected by individual developers.

### Reduced fragmentation

Skills, plugins and MCP remain compatible with their existing ecosystems.

### Ecosystem leverage

DARP can consume capabilities created by other communities and organizations.

### Future-proof architecture

New agent clients can be integrated through adapters instead of requiring a redesign of the Harness.

---

## Negative Consequences

### Greater architectural complexity

DARP must understand multiple standards and agent-specific compatibility mechanisms.

### Adapter maintenance

Agent-specific projections may need to evolve as clients change.

### Partial interoperability

Not every agent will understand every DARP capability.

### More explicit boundaries

The DARP CLI cannot assume that declaring something in the Harness causes an agent to execute it.

---

# 20. Rejected Alternatives

## 20.1 Create a DARP-specific Skill format

Rejected.

Established interoperable Skill formats should be consumed rather than replaced.

---

## 20.2 Create a DARP-specific Plugin format

Rejected.

DARP should understand and manage Agent Plugins rather than create another competing plugin ecosystem.

---

## 20.3 Make coding agents directly execute `.darp/harness.yaml`

Rejected as an assumption.

The Harness is a DARP contract.

Agent clients may require adapters or projections to consume its semantics.

---

## 20.4 Make DARP a coding agent

Rejected.

DARP should remain complementary to coding agents.

---

## 20.5 Make DARP an LLM runtime

Rejected for the current architecture.

Runtime execution requires a separate architectural decision.

---

## 20.6 Build the Registry immediately

Rejected.

The Registry is strategically important but should be developed after the local composition and governance model has matured.

---

# 21. Impact on Existing Specifications

Spec 006 — Harness Manifest remains valid.

It established the declarative Harness manifest and structural validation contract.

This ADR does not retroactively change Spec 006.

Future specifications MUST build on the boundaries defined here.

The Harness Manifest MUST remain compatible with this ADR.

---

# 22. Impact on `darp init`

`darp init` currently establishes the DARP project baseline.

This ADR does not require `darp init` to automatically create `.darp/harness.yaml`.

Harness materialization MUST be defined by a future lifecycle specification.

The preferred initial direction is an explicit Harness lifecycle command such as:

```bash
darp harness init
```

rather than silently creating or activating a Harness during unrelated initialization.

---

# 23. Impact on `darp doctor`

`darp doctor` may eventually validate:

* Harness structure;
* plugin compatibility;
* Skill compatibility;
* MCP declarations;
* policies;
* agent adapters;
* required capabilities;
* evaluation readiness.

The exact checks remain future specifications.

The command remains read-only and diagnostic.

---

# 24. Target Architecture

The resulting conceptual architecture is:

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

The Harness is the central DARP composition and governance abstraction.

The Agent Plugin standard is an interoperability mechanism.

MCP is a tool/context interoperability mechanism.

Agent Skills are reusable capabilities.

Coding agents execute work.

DARP governs and composes the environment.

---

# 25. Roadmap Consequences

The architecture implies the following evolution:

```text
001 Init                         ✓
002 Doctor                      ✓
003 Governance                  ✓
004 Governance Skills           ✓
005 Asset Discovery             ✓
006 Harness Manifest            ✓

007 Ecosystem Architecture      ← this ADR

008 Harness Lifecycle
009 Agent Plugin Interoperability
010 Registry Foundation
011 Policy & Trust
012 Evaluation Contract
013 Harness Composition
014 Agent Adapters / Projections
015 DARP SDK
016 DARP Registry

Future:
DARP Runtime
```

The order may be adjusted by future ADRs and specifications.

---

# 26. Final Decision

DARP will not compete with the emerging Agent Plugin, Agent Skill or MCP ecosystems.

DARP will integrate with them.

DARP will provide a higher-level project abstraction through the Harness, allowing organizations and development teams to compose, govern, validate and eventually evaluate agentic engineering environments independently of the specific coding agent or model provider being used.

The strategic direction is:

> **DARP makes agentic engineering capabilities portable, versionable, reproducible and governable.**

And the fundamental architectural relationship is:

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
