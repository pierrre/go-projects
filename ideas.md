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

## assert — github.com/pierrre/assert

Go test assertion library using generics (no reflection), with auto-updating snapshot assertions (`assertauto`).


## errors — github.com/pierrre/errors

Errors library with message wrapping, stack traces, tags, values, verbose output, and a drop-in std `errors` replacement; composable sub-packages (errbase, errmsg, errstack, errtag, errval, errtmp, errignore, erriter).


## pretty — github.com/pierrre/pretty

Reflection-based pretty printer with a modular `ValueWriter` chain, cycle/recursion detection, max depth, hex dumps, iter.Seq/Seq2 support, and protobuf extensions.


## vld — github.com/pierrre/vld

Generics-based, allocation-free validation library with composable validators (And/Or/If/IfElse/Switch), built-in checks, localized messages (en/fr), and error-path extraction.


## go-libs — github.com/pierrre/go-libs

Shared utility sub-packages (bytesutil, chansutil, goroutine, reflectutil, singleflight, syncutil/atomicutil, unsafeio, weakutil, etc.) used across the author's projects.


## geohash — github.com/pierrre/geohash

Geohash encode/decode library with a CLI front-end (`cmd/geohash`).


## unlimited-channel — github.com/pierrre/unlimited-channel

Unbounded channel that bridges an input and output channel through a single goroutine and an in-memory linked-list queue, with `Close()` draining both ends.


## go-stuff — github.com/pierrre/go-stuff

Shared-module collection of small, self-contained Go experiment sub-packages (five fibonacci variants, a reflection-based function-mocking helper, and a password-entropy estimator).


## file-duplicate — github.com/pierrre/file-duplicate

CLI/library that finds duplicate files by walking filesystems, grouping files by identical size, then SHA256-hashing each same-sized group.


## file-random — github.com/pierrre/file-random

CLI library `filerandom` walks one or more filesystems, collects regular files above a min size, and `rand.Intn`-picks uniformly; the `cmd/file-random` binary loops printing/opening picks.


## mandelbrot — github.com/pierrre/mandelbrot

Mandelbrot fractal renderer: core set computation with specialized per-power iterators, an image package (sequential/parallel render, colorizers), and a small CLI helper for PNG output.


## langton — github.com/pierrre/langton

Langton's ant cellular automaton library with a core grid/ant/rules engine and two CLI frontends (termbox graphical, text stdout).



## cellauto — github.com/pierrre/cellauto

Small cellular automaton library exposing Game of Life and Wireworld rules with two CLI frontends (interactive termbox viewer and PNG-output runner).


## agents — (opencode agent config repo, not Go)

Personal opencode agent-config repo: global AGENTS.md rules, opencode.jsonc provider/model config, and a Makefile that symlinks the repo into ~/.config/opencode.

- [improve] (P3) Commented-out `model` line is dead config with no default model set; either set it or remove it. Refs: opencode/opencode.jsonc:70
- [improve] (P3) `$(CURDIR)` is unquoted in symlink commands; paths with spaces would break installation. Refs: Makefile:3,5
- [improve] (P3) Scaleway baseURL hardcodes a project UUID, coupling config to one project and leaking org identity. Refs: opencode/opencode.jsonc:45
- [refactor] (P3) glm-5.2 model block (limits/cost/variants) is duplicated verbatim across Codix-Iliad and Scaleway providers; factor to avoid drift. Refs: opencode/opencode.jsonc:23,49
- [improve] (P3) `goimports` runs `golang.org/x/tools/cmd/goimports@latest`, unpinned and network-dependent on cache miss. Refs: opencode/opencode.jsonc:102

## di — github.com/pierrre/di

Small generic dependency-injection container with named services, lazy singleton init, and panic-aware builders.

- [improve] (P3) `Container.Close` closes services in alphabetical key order, not reverse-build/dependency order; a Close that touches a dependency finds it already closed, and a subsequent Get rebuilds it outside the snapshot (leak). Refs: container.go:54
- [docs] (P3) `Container.Close` doc omits close-order semantics and does not warn that Close callbacks must not call Get on other services. Refs: container.go:47

## compare — github.com/pierrre/compare

Generics-supported, reflect-based recursive comparator producing path-annotated Differences with cycle detection, custom Func dispatch, and cached .Equal()/.Cmp() method handling.

- [improve] (P3) `comparePointer` skips `compareNil`, so nil-vs-non-nil `*T` reports "only one is valid" (via `Elem` on the nil pointer) instead of "only one is nil" like chan/func/interface — inconsistent and misleading. Refs: compare.go:311-320
- [improve] (P3) Identifier typo: `cmdMethodFuncs`/`cmdMethodFuncsLock` should be `cmp...` (the sibling func/msg consts all use `Cmp`); hurts grepability, copy-paste from the Equal variant. Refs: compare.go:665-666
- [refactor] (P3) `getMethodEqualFunc` and `getMethodCmpFunc` are ~30-line near-duplicates (nil-sentinel cache lookup + signature validation); collapse into one shared helper parameterized by name/outType. Refs: compare.go:619-691
