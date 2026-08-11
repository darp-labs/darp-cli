---
name: architecture
description: Review system design for cohesion, coupling, responsibilities, duplication, extensibility, and alignment with project principles.
---

# Architecture Review Skill

## When to Use

Use during Architecture Review or when a change affects system boundaries,
responsibilities, dependencies, or significant technical decisions.

## Responsibilities

- Read the relevant project context, Constitution, specification, plan and ADRs.
- Evaluate cohesion, coupling, ownership, duplication and extensibility.
- Identify conflicts with approved principles and record actionable recommendations.
- State assumptions and distinguish `PASS`, `WARNING` and `BLOCKED`.

## Workflow

1. Identify the affected components and their responsibilities.
2. Trace dependencies and data or control flow across the change.
3. Check for duplication, unnecessary coupling and constitutional conflicts.
4. Record findings with evidence, assumptions and recommendations.
5. Assign the final review state: `PASS`, `WARNING` or `BLOCKED`.

## Expected Output

Produce a concise review containing scope, assumptions, evidence, findings,
recommendations and a final `PASS`, `WARNING` or `BLOCKED` state.

## Boundaries

Do not invent requirements, approve unrecorded scope changes, or modify source
code. Do not claim that the review or quality gates execute automatically.
