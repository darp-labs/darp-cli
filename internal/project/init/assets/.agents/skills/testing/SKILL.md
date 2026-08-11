---
name: testing
description: Review test coverage, regressions, failure scenarios, compatibility, and the evidence used to validate a change.
---

# Testing Review Skill

## When to Use

Use during Tests and Quality Gates, or when a change may introduce behavioral,
failure-path, regression or compatibility risk.

## Responsibilities

- Read the specification, acceptance criteria, implementation and existing tests.
- Identify missing coverage, regressions, error scenarios and inconsistencies between tests and implementation.
- Discover the project's documented commands before recommending validation.
- Record every command, exit status and relevant evidence.
- State assumptions and distinguish `PASS`, `WARNING` and `BLOCKED`.

## Workflow

1. Map acceptance criteria and changed behavior to existing or required tests.
2. Inspect happy paths, failure paths, regressions and compatibility cases.
3. Run or review the documented validation commands when authorized and available.
4. Record commands, evidence, gaps and assumptions.
5. Assign the final state: `PASS`, `WARNING` or `BLOCKED`.

## Expected Output

Produce a test review with scope, commands, results, evidence, uncovered risks,
assumptions and a final `PASS`, `WARNING` or `BLOCKED` state.

## Boundaries

Do not invent requirements, silently redefine acceptance criteria, or claim that
tests passed without evidence. Do not modify source code or publish artifacts.
