# Changelog

All notable changes to DARP CLI are documented in this file.

This changelog follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and versions follow [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Documentation

- Move internal specifications, plans, ADRs, and strategy to the private repository; keep only approved public documentation in `docs/`.
- Update contributor and agent guidance to use the private `darp-cli-specs` workflow.
- Require English for public user documentation and published specifications.

### Added

- Add the `darp sync` command to explicitly reconcile skills and assets in
  `darp.yml`, with managed ownership, a report of skipped items, and a
  `--dry-run` option.
- Add a shared skill for creating focused branches and commits and opening pull
  requests with `gh`, filling in the repository template when available;
  published for Codex, GitHub Copilot, and Claude Code.
- Add the `darp init` command to initialize DARP projects with contracts, base
  configuration, and the shared skills structure.
- Add the `darp doctor` command to diagnose DARP projects without modifying
  files, including configuration, structure, workflows, skills, templates,
  governance, and contract compatibility.
- Add the `darp --help` and `darp --version` commands.
- Add a specification-driven development workflow, including a lifecycle,
  governance, and contribution documentation.

### Changed

- `darp init` restores only missing files and directories in partial projects,
  preserving existing files and project specifications.
- `darp init` now distributes all four governance skills and the complete
  lifecycle and quality gate contracts, and updates `darp.yml` additively and
  safely.
- The CLI version is calculated from Git during the build, including a
  development indicator and a modified-worktree indicator when applicable.

### Added

- Automatically discover pre-existing assets during `darp init` and register
  them additively and safely in the new `assets` section of `darp.yml`. Discovery
  is local, deterministic, and read-only; it recognizes patterns from the
  GitHub, Claude, Codex, Cursor, and Gemini families, ignores symlinks, and
  reports already-registered entries and ambiguous paths without overwriting
  anything.
- `darp doctor` now validates registered assets (existence, relative path within
  the project, and expected family and type); a missing asset is reported as a
  `WARNING` and does not block the project.

### Fixed

- `darp sync` now validates Agent Skills frontmatter before registering skills,
  reports invalid structures without removing ownership in unobserved scopes,
  correctly serializes the simulated manifest, and rejects structurally invalid
  skill entries before writing.
- `darp init` now ignores unsupported files in conventional asset directories,
  avoiding discovery failures.
- Discovery now centralizes official asset classification, reports semantic
  conflicts by category and normalized name, and prevents their automatic
  registration.
- `darp doctor` now rejects incompatible path, family, and type combinations,
  as well as manually registered symlinks.
- `darp init` now validates configuration and discovered assets before creating
  contracts, cleans up temporary files when replacement fails, and rolls back
  partial changes after persistence failures.
- `darp doctor` now classifies a directory registered in place of an asset as
  `FAIL`; only a missing asset remains a `WARNING`.

## Maintaining This File

- Record user-facing changes under `Unreleased` during development.
- Use the `Added`, `Changed`, `Fixed`, `Removed`, `Security`, and `Deprecated`
  categories only when needed.
- When publishing a release, turn `Unreleased` into a versioned section dated
  `YYYY-MM-DD` and add a new `Unreleased` section above it.
- Do not record every commit; focus on changes relevant to users, contributors,
  and CLI consumers.

[Unreleased]: https://github.com/darpbr/darp-cli/compare/HEAD...main
