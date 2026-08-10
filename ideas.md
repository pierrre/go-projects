# Project ideas

# Processing

## Conventions

- Workspace: `/home/pierre/Git/pierrre/`. Each repo is a subdirectory named after the repo.
- Before working in any repo: `git checkout main && git pull`. After any work in a repo: return to `main`.
- Periodically: `git branch --merged main | grep -v '^\*\|main' | xargs -r git branch -d` and `git remote prune origin`.
- `make all` runs `build test lint`. Outside CI, lint runs `--fix` and may modify files — commit only after `make all` completes.
- Steps can be executed directly or via sub-agents. No parallel sub-agents.

## PR scan

Mandatory before each idea (the reviewer may have commented meanwhile), or standalone to clear unaddressed comments.

### 1. List and filter PRs

1. `gh search prs --owner pierrre --state open --json repository,number,author --limit 500` (`--limit 500` is mandatory: the default 30 silently truncates).
2. Archived repos (one call): `gh repo list pierrre --json name,isArchived --limit 500`. Skip PRs in archived repos.
3. Skip PRs where `author.login != pierrre`.
4. `gh search prs` doesn't return branch names. For each repo with remaining PRs: `gh pr list --repo pierrre/<repo> --state open --json number,headRefName`.
5. For non-archived repos with own PRs not cloned locally: clone first.

### 2. Fetch and classify comments

For each remaining PR:
1. Issue comments: `gh api repos/pierrre/<repo>/issues/<n>/comments`.
2. Review comments (line-level): `gh api repos/pierrre/<repo>/pulls/<n>/comments`.
3. Merge both, sort by `created_at`.
4. `Done:` cutoff: the last comment whose body starts with `Done:` marks the cutoff; comments after it (or all if none) are unaddressed.
5. Actionable: requests a code change, asks a question, or reports a bug. Skip approvals, emoji, acknowledgments.
6. If you disagree with a requested change: reply `Done: not applied — <reason>`.

Note: `gh api --paginate` concatenates JSON arrays — omit it (most PRs fit one page) or pipe through `jq -s 'add'`.

### 3. Address comments per PR

For each PR with actionable comments:
1. `git checkout <branch>`.
2. Implement all changes — one commit per PR, listing each change in the message.
3. Run `make all`.
4. `git add -A && git commit -m "..." && git push`.
5. Reply to each comment individually (never batch):
   - Review comments: `gh api repos/pierrre/<repo>/pulls/<n>/comments/<comment_id>/replies -f body="Done: <summary>"`.
   - Issue comments: `gh pr comment <n> --repo pierrre/<repo> --body "Done: <summary>"`.
6. Return to `main`.

### 4. Rebase own open feature branches

For every own open feature branch (including those from Step 3):
1. `git checkout <branch> && git fetch origin`.
2. Skip if up to date: `git merge-base --is-ancestor origin/main HEAD` (exit 0).
3. `git rebase origin/main && git push --force-with-lease`.
4. Return to `main`.

## Idea processing

- Process one idea at a time, priority order P1 → P2 → P3 (file order within each).
- Before each idea, run the PR scan above.
- Validate the idea against current source. Reject if the issue no longer exists, is a false positive, would break the public API, or is too large/risky for one focused PR.
- If valid: ensure `main` is up to date, create `feature/<slug>` off `main`, implement, `make all`, commit, push, open a PR (base `main`, assignee `pierrre`, review `pierrre`). Cover all new/modified code with tests — `go test -coverprofile=cover.out ./... && go tool cover -func=cover.out`, add cases until complete. Note untestable parts in the PR body. Remove `cover.out`.
- If rejected: no branch, no PR. Record the reason in the removal commit.
- After each idea (implemented or rejected): remove its bullet from this file, commit on `go-projects` `main` (`Remove <idea>: implemented` or `Remove <idea>: rejected: <reason>`), and push. No branch or PR for this edit.

Generated: 2026-08-09
Scope: /home/pierre/Git/pierrre workspace (Go 1.26.0)
Method: direct review of each project's key files (golang-review + golang-new skills applied).

Legend:
  Tags:     [bug] [feat] [refactor] [improve] [test] [docs]
  Priority: P1 high  P2 medium  P3 low


## di — github.com/pierrre/di

Small generic dependency-injection container with named services, lazy singleton init, and panic-aware builders.

- [improve] (P3) `Container.Close` closes services in alphabetical key order, not reverse-build/dependency order; a Close that touches a dependency finds it already closed, and a subsequent Get rebuilds it outside the snapshot (leak). Refs: container.go:54
- [docs] (P3) `Container.Close` doc omits close-order semantics and does not warn that Close callbacks must not call Get on other services. Refs: container.go:47
