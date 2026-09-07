# Plan 002 - `darp doctor`

## Status

Approved and implemented — validated on 2026-08-27.

## Related Specification

- `spec.md`

## Implementation Strategy

Implement the doctor as a read-only application service made of independent checks.
The CLI layer only dispatches the command and renders its result. Each check reads
the project independently, so a malformed configuration does not prevent the
remaining contracts from being reported.

The supported project contract is `1.x`. YAML parsing uses `gopkg.in/yaml.v3` so
the CLI validates YAML syntax rather than a subset of the format.
Agent skills live under `.agents/skills/<name>/SKILL.md`; supporting resource
directories are optional and are not part of the DARP metadata tree.

## Affected Areas

- `internal/project/doctor/` for the diagnosis model, checks, and service
- `internal/cli/` for command dispatch and output
- `internal/project/init/` so new projects use the official `darp.yml` contract
- unit tests for doctor, CLI behavior, and initializer output
- `README.md` for command usage

## Validation Strategy

- `go test ./...`
- run `darp doctor` against a valid temporary project
- run it against invalid and older-version projects, checking exit codes and
  ensuring the project tree is unchanged

## Assumptions

- The mandatory structure is `.darp/governance`, `.darp/workflows`,
  `.agents/skills`, and `.darp/templates`.
- `workflows.default` names the required workflow file and its declared `name`.
- A workflow step maps to an agent skill directory with the same name under
  `.agents/skills`.
- The project configuration is `darp.yml`, and the lifecycle document is
  `.darp/lifecycle.md`.
