---
name: create-pr
description: Create a focused Git branch and commit, then push and open a GitHub pull request with the repository's PR template using the gh CLI. Use when the user asks to prepare or open a PR.
---

# Create a Branch, Commit, and Pull Request

Use this skill when the user asks you to prepare a contribution or open a pull
request. Follow repository instructions and conventions first. Keep the change
focused, understandable, and ready for review.

## Workflow

1. **Inspect before changing Git state.** Read applicable `AGENTS.md` files and
   contribution guidance. Check `git status`, the current branch, the diff, and
   recent commit and branch naming conventions. Identify the repository's
   documented validation commands. Check `gh auth status` and repository
   metadata when needed. If GitHub CLI is unavailable or unauthenticated, report
   the blocker without switching to another publishing method.
2. **Protect existing work.** Determine which changes belong to the user's
   request. Preserve unrelated staged and unstaged changes. Never reset, clean,
   amend, or overwrite existing work to make the workflow convenient. If
   requested changes cannot be separated safely, stop and explain the issue.
3. **Validate and review.** Run the relevant checks required by the repository
   and task. Review the complete diff, including untracked files, for scope,
   correctness, accidental generated files, credentials, tokens, and other
   sensitive data. Do not claim checks passed unless they ran successfully.
4. **Create a focused branch.** Use a descriptive, short, lowercase
   `type/summary` name when compatible with repository conventions (for
   example, `docs/add-install-guide` or `fix/cache-expiry`). Check that it does
   not already exist locally or remotely. Do not switch away from a branch with
   work in progress until that work is safely accounted for. Never reuse the
   default branch for the contribution.
5. **Stage only intended files.** Stage explicit paths or carefully selected
   hunks; do not use a blanket `git add .` when unrelated files may be present.
   Inspect `git diff --cached` and `git status` before committing. Keep the
   commit cohesive. Follow the repository's commit convention; if none exists,
   use a concise imperative subject, optionally with a Conventional Commit
   type, and avoid bundling unrelated changes.
6. **Commit and verify.** Commit only after the staged diff is reviewed. Confirm
   the commit, branch, and remaining working tree state. Do not amend an
   existing commit unless the user specifically asks.
7. **Push and open the PR with `gh`.** The user's explicit request to create or
   open a PR authorizes the necessary push and PR creation. Do not ask again for
   that same authorization. Push the feature branch normally (never force-push)
   and use `gh pr create` to open the pull request. Do not merge, approve, or
   enable auto-merge unless separately requested.
8. **Use the repository's PR template.** Before creation, look for GitHub pull
   request templates in recognized locations, including
   `.github/pull_request_template.md`, `.github/PULL_REQUEST_TEMPLATE.md`, and
   files under `.github/PULL_REQUEST_TEMPLATE/`. Follow any repository-specific
   guidance if it defines a different location. Select the applicable template,
   preserve its useful sections and checklist, and replace every instructional
   comment and placeholder with accurate, concise content. Remove a section
   only when it is genuinely inapplicable, marking it `N/A` when that best
   preserves the template. Fill in references, change type, validation,
   impact, risks, and review notes from evidence; never invent an issue number
   or claim an unchecked validation passed. If no template exists, write a
   concise body covering context, changes, validation, and relevant risks.
9. **Confirm the result.** Verify creation with `gh pr view` (or equivalent),
   and report the branch, commit, PR URL, checks run, and any remaining
   limitations. If a command fails, report the exact stage and useful error;
   do not retry a remote mutation blindly or create a duplicate PR.

## Quality bar

- Keep the branch, commit, and PR title specific and consistent with the
  change.
- Make the PR easy to review: explain why the change is needed and what changed,
  link only verified references, and state validation accurately.
- Prefer small, self-contained commits and avoid unrelated cleanup.
- Preserve repository conventions when they are stronger than these defaults.
- Do not expose secrets in commit content, PR text, command output, or logs.
