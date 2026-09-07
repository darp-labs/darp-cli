# Tasks 002 - `darp doctor`

## Status

Completed — implementation and validation completed on 2026-08-27.

## Related Plan

- `plan.md`

## Task List

1. Add YAML parsing support and the doctor result/check model.
2. Implement independent configuration, structure, workflow, skill, template,
   governance, and version checks.
3. Add CLI dispatch, stable output, summary, and exit-code handling.
4. Update `darp init` to create the official project contract.
5. Add automated tests for passing, warning, failing, read-only, and CLI flows.
6. Update command documentation.

Additional hardening:

7. Validate agent-skill paths, regular `SKILL.md` files, optional resource
   directories, traversal attempts, symlink containment, and deterministic
   diagnostics.
