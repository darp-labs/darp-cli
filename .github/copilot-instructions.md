# GitHub Copilot Instructions

This repository follows Specification-Driven Development.

Before implementing code:

- read AGENTS.md
- read README.md
- understand the project vision
- verify with a maintainer whether an approved specification exists in the private `darp-cli-specs` repository under `specs/internal/darp-cli/.spec/specs/`
- read the matching `spec.md`, `plan.md`, and `tasks.md` only when authorized to access that repository

Never implement functionality directly from a prompt if no approved specification exists. Contributors without access to the private specs must ask a maintainer for an approved public scope; do not copy private documents into this repository.

Prefer:

- small pull requests
- deterministic implementations
- idiomatic Go
- reusable components
- low coupling

Avoid:

- speculative abstractions
- premature optimization
- unnecessary dependencies

When generating code:

- follow Effective Go
- use the standard library whenever possible
- keep functions focused
- write meaningful comments only when needed

Always preserve architectural consistency.

## Documentation and Publication

- Follow the `Documentation Publication` rules in `AGENTS.md` for documentation changes.
- Do not copy private source specifications or planning documents into this public repository, prompts, skills, or pull requests.
- For a public document, confirm its classification, source, rights, sanitization, review, and maintainer approval before adding it under `docs/`.
- Write all public user documentation and published specifications in English. Internal documents may use another language unless they are approved for public release.
- Treat AI scans as assistance only; they cannot confirm ownership or approve publication.

Feature planning artifacts belong together in the private
`darp-cli-specs/specs/internal/darp-cli/.spec/specs/<number>-<short-name>/`.
Do not create active specs, plans, or tasks in this public repository. The
`.darp/` directory is reserved for public runtime contracts.
