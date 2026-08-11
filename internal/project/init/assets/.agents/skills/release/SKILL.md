---
name: release
description: Review release notes, features, breaking changes, compatibility impact, and decisions required before release.
---

# Release Review Skill

## When to Use

Use during Release Notes, Review and Release preparation after implementation,
tests and required quality gates have evidence.

## Responsibilities

- Review release notes for features, fixes, breaking changes and compatibility.
- Verify that user-visible changes are recorded in `CHANGELOG.md` under
  `Unreleased` during development, and in the corresponding versioned section
  when a release is prepared.
- Check traceability to the specification, plan and review decision.
- Identify decisions, warnings and blockers that must be resolved before release.
- State assumptions and distinguish `PASS`, `WARNING` and `BLOCKED`.

## Workflow

1. Read the approved scope and the evidence from implementation and reviews.
2. Summarize user-visible changes and compatibility impact.
3. Check for undisclosed breaking changes and unresolved release decisions.
4. Record assumptions, evidence and required follow-up.
5. Assign the final state: `PASS`, `WARNING` or `BLOCKED`.

## Expected Output

Produce a release review containing release notes status, feature and breaking
change summary, compatibility impact, decisions, assumptions and a final
`PASS`, `WARNING` or `BLOCKED` state.

## Boundaries

Do not invent release requirements, publish packages or artifacts, alter source
code, or imply that the CLI executes this review automatically.
