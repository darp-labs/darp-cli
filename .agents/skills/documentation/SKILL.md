---
name: documentation
description: Keep project documentation aligned with the implementation. Use when features, commands, configuration, or project structure change.
---

# Documentation Skill

## When to Use

Use during Documentation Review or when implementation changes commands,
configuration, behavior, examples or project structure.

## Responsibilities

- Compare README and technical documentation with the implementation.
- Update documentation when behavior or structure has changed.
- Validate internal links and documented paths.
- Classify new or materially changed documents and record owner, source/authorship, license or attribution, and audience.
- Treat `docs/` in the public repository as approved-publication-only; consult the private publication policy before publishing.
- Review examples, attachments, metadata, history, and external indexing. Escalate uncertain rights or possible sensitive content to a maintainer.
- Treat automated or AI scans as assistance, never as proof of rights or publication approval.
- Remove stale commands, examples and claims.
- State assumptions and distinguish `PASS`, `WARNING` and `BLOCKED`.

## Workflow

1. Read the approved specification and inspect the implementation and structure.
2. Identify documentation affected by the change and classify any new or materially changed document.
3. Check links, paths, commands and examples against existing files and behavior.
4. For public documentation, confirm provenance, rights, sanitization, history/indexing review, and human maintainer approval before publication.
5. Update documentation only when the implementation supports the claim.
6. Record assumptions and assign `PASS`, `WARNING` or `BLOCKED`.

## Expected Output

Produce a documentation review with inspected paths, link results, findings,
assumptions, changes made and a final `PASS`, `WARNING` or `BLOCKED` state.

## Boundaries

Do not invent requirements or unsupported behavior. Do not modify source code,
commands or configuration to make documentation pass. Do not claim automatic
execution of skills or quality gates.
