# Contributing to DARP CLI

Thank you for contributing to DARP CLI. The project aims to make reusable AI assets versionable, shareable, and installable as packages.

## Before You Start

Read these documents before proposing a change:

- [README](README.md)
- [AGENTS.md](AGENTS.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)

Specifications, plans, ADRs, and internal context are maintained in the private `darp-cli-specs` repository. For a functional contribution, ask a maintainer to confirm the approved scope; do not include copies of private documents in this public repository.

For a significant change, look for an existing issue or open a discussion before implementing it. The project follows Specification-Driven Development; features must not be implemented without an approved specification.

## Contribution Types

Contributions are welcome in the following areas:

- Go code and bug fixes;
- automated tests;
- documentation and examples;
- skills, prompts, and other AI assets;
- governance and development-tool improvements;
- bug reports and feature proposals.

## Local Setup

Requirements:

- Git;
- Go compatible with the version declared in [go.mod](go.mod).

Clone the repository and run the available checks:

```bash
git clone https://github.com/darpbr/darp-cli.git
cd darp-cli
make verify
```

Useful commands:

| Command | Purpose |
| --- | --- |
| `make build` | Builds the CLI for the current platform. |
| `make test` | Runs the Go tests. |
| `make lint` | Runs `golangci-lint` or falls back to `go vet`. |
| `make fmt` | Formats the Go files. |
| `make verify` | Runs formatting, tests, and lint. |
| `make run` | Runs the CLI locally. |

## Development Workflow

1. Create a descriptive branch from `main`.
2. For functional changes, confirm the approved specification in the private `darp-cli-specs` repository with a maintainer; contributors without access must follow an approved public scope.
3. Keep the change small, focused, and compatible with the existing architecture.
4. Add or update tests when behavior changes.
5. Update the affected documentation. Write all public user documentation and published specifications in English.
6. Run `make verify` before opening a pull request.
7. Open a pull request using the repository template.

The project workflow is:

`Vision → Constitution → Specification → Plan → Tasks → Implementation → Review`

## Implementation Standards

- Prefer idiomatic Go and the standard library when sufficient.
- Keep functions focused, packages cohesive, and coupling low.
- Preserve existing behavior and avoid unnecessary dependencies.
- Do not add secrets, credentials, tokens, or personal data to the repository.
- Do not document features that do not exist yet.
- Explain decisions affecting architecture, compatibility, or security in the pull request.

## AI-Assisted Contributions

AI tools may be used, but the author remains responsible for all submitted content. Critically review the output, verify sources, run tests, and make sure you understand the changes and their impact.

## Issues and Vulnerabilities

Use issues for reproducible bugs, project questions, and improvement proposals. Search existing issues before opening a new one, and include context, reproduction steps, and the expected result when applicable.

Do not disclose vulnerabilities, credentials, or sensitive information in public issues or pull requests. For security or privacy concerns, contact the maintainers through the [@darpbr](https://github.com/darpbr) profile and request a private channel.

## Pull Requests

Pull requests should:

- describe the problem and the solution;
- reference the related issue or specification;
- report tests and commands that were run;
- record risks, limitations, and potential breaking changes;
- contain only changes related to the stated goal;
- be ready for review or marked as Draft while still in progress.

Maintainers may request changes to code, tests, documentation, architecture, or scope. Reply to comments in the pull request to keep the public history traceable.

## Review and Merge

A pull request may be merged when its scope is clear, applicable checks pass, review comments are addressed, and documentation is aligned. Final approval is the responsibility of the maintainers defined in [CODEOWNERS](.github/CODEOWNERS).

## Code of Conduct

Participation in the project is subject to the [Code of Conduct](CODE_OF_CONDUCT.md). Treat conduct concerns respectfully and report them privately to the maintainers when needed.
