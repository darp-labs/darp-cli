# DARP CLI — Roadmap

> **Status:** Documento de direção do projeto
> **Última atualização:** 2026-09-09
> **Objetivo:** Registrar a evolução arquitetural e estratégica do DARP CLI e do ecossistema DARP.

---

# 1. Visão

DARP é uma plataforma aberta para compor, governar, distribuir e avaliar capacidades reutilizáveis para engenharia de software assistida por agentes.

O DARP deve permanecer:

* agent agnostic;
* provider agnostic;
* interoperável;
* determinístico;
* extensível;
* baseado em padrões abertos.

A arquitetura de longo prazo é:

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

# 2. Strategic Positioning

DARP is:

> **An open platform for composing, governing, distributing, and evaluating reusable capabilities for agentic software engineering.**

DARP does not replace:

* coding agents;
* LLM providers;
* Agent Skills;
* Agent Plugins;
* MCP;
* IDEs.

DARP provides a project-level control and composition layer above these capabilities.

---

# 3. Completed Foundation

## 001 — Init

Status: Completed.

Provides the initial DARP project structure.

---

## 002 — Doctor

Status: Completed.

Provides deterministic, read-only project diagnosis.

---

## 003 — Governance

Status: Completed.

Establishes DARP governance contracts.

---

## 004 — Governance Skills

Status: Completed.

Provides reusable governance skills for documentation, architecture, testing and release.

---

## 005 — Asset Discovery

Status: Completed.

Discovers and registers supported AI assets without silently modifying discovered resources.

---

## 006 — Harness Manifest

Status: Completed.

Provides:

```text
.darp/harness.yaml
```

as the canonical Harness manifest path and implements structural parsing and validation.

The current implementation does not execute Harness resources.

---

# 4. Current Architectural Decision

## 007 — DARP Ecosystem Architecture

Status: Architectural Decision Accepted.

Defined by:

```text
docs/ADR/ADR-002-darp-ecosystem-architecture-and-agent-interoperability.md
```

The decision establishes:

* DARP positioning;
* Agent Plugin interoperability;
* Agent Skill interoperability;
* MCP interoperability;
* Harness boundaries;
* Registry direction;
* Runtime boundaries;
* Policy direction;
* Evaluation direction;
* adapter/projection architecture;
* vendor neutrality.

---

# 5. Next Development Phases

## 008 — Harness Lifecycle

Objective:

Define how a Harness is:

* inspected;
* initialized;
* validated;
* updated;
* maintained.

Expected future command:

```bash
darp harness init
```

Potential commands:

```bash
darp harness inspect
darp harness validate
```

Exact commands require a future specification.

Important principle:

The Harness should not be silently created or activated as a side effect of unrelated commands.

---

# 6. Agent Plugin Interoperability

Objective:

Allow DARP to understand the Agent Plugin standard.

Potential capabilities:

```text
discover
inspect
validate
catalog
install
version
verify
evaluate
```

DARP must not modify the Agent Plugin format.

DARP must not create a competing plugin standard.

---

# 7. Harness Composition

Objective:

Define how Harnesses compose:

```text
Agent Plugins
Skills
MCP
Assets
Providers
Policies
Workflows
Evaluation
Compatibility
```

The Harness composes existing capabilities.

It does not redefine them.

---

# 8. Agent Adapters / Projections

Objective:

Allow DARP Harness semantics to be consumed by different agent environments.

Potential targets include:

```text
AGENTS.md
Copilot instructions
Claude configuration
other agent-specific configuration
```

The Harness remains the DARP source of truth.

Adapters are projections.

Future specifications must define:

* generation;
* synchronization;
* ownership;
* conflict resolution;
* precedence;
* compatibility.

---

# 9. Policy and Trust

Objective:

Establish a security and governance model for agentic capabilities.

Potential concepts:

```text
Trust
Provenance
Publisher
Signature
Checksum
Permissions
Network
Filesystem
Secrets
Execution
Audit
```

The system should follow least-privilege principles.

---

# 10. Evaluation Contract

Objective:

Define reproducible evaluation of:

* Assets;
* Skills;
* Agent Plugins;
* Harnesses;
* agent/provider combinations.

Potential dimensions:

```text
Task
Fixture
Expected behavior
Tests
Quality gates
Security
Cost
Time
Result
```

---

# 11. Registry Foundation

Objective:

Define the future DARP Registry.

Potential resources:

```text
Assets
Agent Plugins
Skills
MCP
Harnesses
Policies
Templates
```

Potential capabilities:

* discovery;
* versioning;
* dependency resolution;
* distribution;
* provenance;
* trust;
* compatibility;
* evaluation metadata.

Registry implementation should come only after the local composition model has matured.

---

# 12. DARP SDK

Objective:

Provide reusable libraries for applications integrating with DARP.

Potential consumers:

* CLI;
* Registry;
* IDE integrations;
* enterprise tooling;
* future runtime;
* automation systems.

The SDK architecture should remain independent of a particular coding agent or provider.

---

# 13. DARP Registry

After the Registry Foundation specification, implement the Registry service.

Potential architecture:

```text
                  DARP Registry
                       │
          ┌────────────┼────────────┐
          │            │            │
       Catalog      Metadata      Artifacts
          │            │            │
          └────────────┼────────────┘
                       │
                    DARP CLI
```

Detailed architecture will be defined by future ADRs.

---

# 14. Future Runtime

A DARP Runtime may eventually consume Harness definitions.

Potential responsibilities:

* agent execution;
* provider selection;
* tool execution;
* policy enforcement;
* sandboxing;
* observability;
* evaluation;
* execution state.

The Runtime is intentionally separated from the DARP CLI.

It requires a separate architectural decision.

---

# 15. Ecosystem Evolution

The long-term ecosystem is:

```text
                     DARP Ecosystem
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
     DARP CLI          DARP SDK         DARP Registry
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
                        Harness
                           │
            ┌──────────────┼──────────────┐
            │              │              │
         Plugins         Assets         Policy
            │              │              │
       Agent Plugins     Skills        Security
            │
       ┌────┴────┐
       │         │
     Skills     MCP
       │
       ▼
  Agent Clients
       │
 ┌─────┼─────────────┐
 ▼     ▼             ▼
Codex Copilot      Claude
       │
       ▼
  Software Delivery
```

---

# 16. Architectural Priorities

The project should prioritize:

1. Interoperability
2. Security
3. Determinism
4. Reproducibility
5. Vendor neutrality
6. Enterprise governance
7. Evaluation
8. Developer experience
9. Extensibility
10. Ecosystem adoption

---

# 17. What DARP Should Avoid

DARP should avoid:

* creating competing open standards without strong justification;
* coupling itself to a single coding agent;
* coupling itself to a single LLM provider;
* silently executing discovered capabilities;
* turning the CLI into a hidden runtime;
* creating unnecessary abstractions;
* building a marketplace before the underlying package and trust model is mature.

---

# 18. Strategic Principle

The DARP ecosystem should follow:

> **Adopt where standards exist. Compose where standards meet. Govern where organizations need control. Define only where the ecosystem has a genuine gap.**

This principle should guide future architecture decisions.
