# Quality Gates

These gates are deterministic review criteria for a change. They are not
executed automatically by the DARP CLI in this specification.

For every gate, record assumptions and classify the result as `PASS`, `WARNING`
or `BLOCKED`/`FAIL`. An accepted warning must identify who accepted it and why.

## Build

- Objective: confirm that the project can be built with its documented tooling.
- Verifiable criteria: the documented build command completes successfully;
  compiler or packaging errors are absent.
- Reviewer: implementation owner, with the project's normal build tooling.
- Expected evidence: exact command, exit status and relevant output or artifact.
- PASS: the documented build completes successfully.
- BLOCKED/FAIL: the command cannot run, fails, or the project has no usable
  documented build procedure.

## Tests

- Objective: confirm expected behavior, failure paths and regression safety.
- Verifiable criteria: required tests pass; changed behavior has coverage or an
  explicit rationale; failure and compatibility scenarios are considered.
- Reviewer: `testing` skill or designated test reviewer.
- Expected evidence: exact commands, exit statuses, test summary and recorded
  assumptions.
- PASS: required tests pass and no unresolved test gap blocks acceptance.
- BLOCKED/FAIL: a required test fails, cannot run, or an unexplained coverage
  or regression risk remains.

## Documentation

- Objective: ensure documentation accurately describes the implemented behavior.
- Verifiable criteria: README and technical documentation match the code;
  examples and internal links are valid; nonexistent behavior is not claimed.
- Reviewer: `documentation` skill.
- Expected evidence: changed documentation paths, link checks and findings.
- PASS: documentation is aligned and links resolve.
- BLOCKED/FAIL: stale, contradictory or unverifiable documentation remains.

## Architecture

- Objective: verify cohesion, coupling, responsibilities, duplication,
  extensibility and constitutional alignment.
- Verifiable criteria: responsibilities are clear; unnecessary duplication and
  coupling are absent; relevant ADRs and principles are respected.
- Reviewer: `architecture` skill.
- Expected evidence: review findings, assumptions, recommendations and ADR
  references where applicable.
- PASS: no unresolved architectural blocker exists.
- BLOCKED/FAIL: a constitutional conflict, unclear ownership or material design
  risk remains unresolved.

## Compatibility

- Objective: verify that existing supported behavior and contracts remain valid.
- Verifiable criteria: compatibility-sensitive tests pass; public commands,
  configuration and documented paths remain valid unless explicitly approved.
- Reviewer: implementation owner and project reviewer.
- Expected evidence: regression results, compatibility comparison and affected
  paths.
- PASS: supported existing behavior remains compatible or the impact is
  explicitly approved.
- BLOCKED/FAIL: an unapproved breaking change or unresolved regression exists.

## Release Notes

- Objective: make the change and its release impact understandable and
  traceable.
- Verifiable criteria: features, fixes, breaking changes, compatibility impact
  and decisions needed before release are recorded.
- Reviewer: `release` skill.
- Expected evidence: release-notes entry or an explicit not-applicable record,
  plus links to the specification and review decision.
- PASS: release impact is complete and traceable.
- BLOCKED/FAIL: release impact is unknown, breaking changes are undisclosed,
  or a required decision is unresolved.
