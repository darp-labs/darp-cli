# DARP Lifecycle

## Purpose

The DARP Lifecycle defines how work is planned, designed, implemented,
validated and released. It is the operational detail of the principles in
`.spec/constitution.md` and is intended for humans and AI agents.

The lifecycle is provider-agnostic and independent of programming language,
framework, repository layout and AI model.

## Principles

Every change should be intentional, traceable, incremental, explainable and
reviewable. Planning precedes coding, and implementation must not redefine an
approved specification.

## Lifecycle

```text
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
Tests
  ↓
Architecture Review
  ↓
Documentation Review
  ↓
Quality Gates
  ↓
Review
  ↓
Release
  ↓
Completed
```

## Phase descriptions

### Vision

Defines why the project exists and the long-term direction. It changes rarely.

- Input: the project's mission and long-term direction.
- Output: a stable statement of intent.
- Pass when: the change can be traced to the project's stated direction.

### Constitution

Defines permanent engineering principles. It has precedence over
specifications, plans and implementation decisions.

- Input: the project's enduring engineering principles.
- Output: the principles and decision hierarchy contributors must follow.
- Pass when: the work does not contradict a constitutional principle.

### Project Context

Defines terminology, constraints, current architecture and goals needed to
understand the work.

- Input: repository structure, domain terminology and current constraints.
- Output: shared context sufficient to evaluate the proposed change.
- Pass when: contributors can identify the affected system boundaries and
  constraints.

### Specification

Defines what will be built, its boundaries, requirements and acceptance
criteria. It must be approved before implementation begins.

- Input: the project context and the problem to be solved.
- Output: an approved specification with scope, requirements and acceptance
  criteria.
- Pass when: scope and acceptance criteria are explicit and the specification
  is approved.

### Architecture Decision Record (ADR)

Records significant technical decisions and their trade-offs when the
specification requires them.

- Input: an approved specification and unresolved architectural choices.
- Output: an ADR, or an explicit record that no ADR is required.
- Pass when: significant decisions and their consequences are recorded.

### Implementation Plan

Defines how the approved specification will be implemented, including affected
areas, dependencies, milestones and validation.

- Input: the approved specification and any applicable ADR.
- Output: an approved implementation plan.
- Pass when: affected areas, dependencies and validation strategy are defined.

### Tasks

Breaks the plan into small, independently understandable and verifiable
checkbox tasks. Dependencies and evidence must be explicit.

- Input: the approved plan.
- Output: an approved task list with dependencies and validation evidence for
  each task.
- Pass when: every implementation task is actionable, checkable and validated
  before it is marked complete.

### Implementation

Executes only approved tasks. The implementation must remain within the
specification and plan; discovered scope changes require a document update and
approval before coding continues.

- Input: approved specification, plan and tasks.
- Output: the implementation and its recorded changes.
- Pass when: the change stays within approved scope and all implementation
  tasks have supporting evidence.

### Tests

Validates expected behavior, failure paths, regressions and compatibility using
the project's own documented toolchain. The `testing` skill supports this phase.

- Input: the implementation, acceptance criteria and documented toolchain.
- Output: test results and evidence for expected and failure behavior.
- Pass when: required tests pass or any warning/blocker is explicitly recorded.

### Architecture Review

Checks responsibilities, coupling, cohesion, duplication, extensibility and
alignment with the Constitution. The `architecture` skill produces findings
and recommendations.

- Input: the implementation, project context, Constitution and ADRs.
- Output: architecture findings, assumptions and recommendations.
- Pass when: no unresolved architectural blocker remains.

### Documentation Review

Checks that README files, technical documentation, examples, commands and
project structure describe the implemented behavior accurately. The
`documentation` skill supports this phase.

- Input: the implementation, repository structure and existing documentation.
- Output: documentation findings and any required documentation updates.
- Pass when: documented behavior matches the implementation and internal links
  resolve.

### Quality Gates

Checks Build, Tests, Documentation, Architecture, Compatibility and Release
Notes according to `.darp/governance/quality-gates.md`. Gates are review
criteria; this file does not imply automatic CLI execution.

- Input: implementation artifacts and evidence from the preceding phases.
- Output: a result for each documented quality gate.
- Pass when: every required gate is PASS or has an explicitly accepted
  WARNING; BLOCKED/FAIL results prevent completion.

### Review

Combines technical, architectural and documentation findings. Unresolved
BLOCKED findings prevent completion.

- Input: test results, review findings and quality-gate results.
- Output: a review decision with assumptions and unresolved findings.
- Pass when: no unresolved BLOCKED finding remains and the change is approved.

### Release

Makes the approved change available for consumption and records its release
notes, compatibility impact and traceability to the specification.
The `release` skill supports the release-notes review.

- Input: an approved review and passing quality gates.
- Output: release notes and a traceable release decision.
- Pass when: release notes, compatibility impact and required decisions are
  recorded; publication itself is outside these governance documents.

### Completed

The change is complete when all required tasks are checked, validation evidence
exists, quality gates pass or have an explicitly accepted warning, and no
unresolved blocker remains.

- Input: the release decision and all lifecycle evidence.
- Output: a completed, auditable change record.
- Pass when: tasks, evidence, gates and blockers satisfy the completion criteria.

## Skills

The governance skills in `.agents/skills/` support reviews but do not replace
the lifecycle and do not execute automatically in the CLI. Their instructions
must remain agnostic to language, framework and project type.

## Document hierarchy

When documents disagree, follow this order:

1. Vision
2. Constitution
3. Project Context
4. Approved Specification
5. Architecture Decision Record
6. Implementation Plan
7. Tasks

Lower-level documents must not contradict higher-level documents. This
lifecycle details the shorter workflow summarized in the Constitution.

## AI collaboration

Agents must read the relevant project context before proposing changes, explain
assumptions, preserve existing assets, avoid speculative behavior and request
clarification when requirements conflict.
