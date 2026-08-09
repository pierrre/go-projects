# Project ideas

# Processing

- Keep local clones of all repos up to date: before working in a repo, `git checkout main && git pull`. Periodically clean up merged or deleted branches with `git branch --merged main | grep -v '^\*\|main' | xargs -r git branch -d` and `git remote prune origin`.
- Keep feature branches you own up to date by rebasing them on top of `main`: check out the branch, `git fetch origin && git rebase origin/main`, force-push (`git push --force-with-lease`), then return to `main`. Do this during the PR scan step for every open feature branch, even if there are no review comments to address. Skip branches that are not yours or that have been merged.
- Before processing each idea, dispatch a sub-agent to scan open PRs for unaddressed review comments from `pierrre`. List PRs with `gh search prs --owner pierrre --state open --json repository,number,title,url`. Skip PRs in archived repositories (check with `gh repo view <repo> --json isArchived`). For each non-archived PR, fetch issue comments (`gh api repos/<owner>/<repo>/issues/<n>/comments`) and review comments (`gh api repos/<owner>/<repo>/pulls/<n>/comments`), merge and sort by timestamp. The last comment starting with `Done:` marks the cutoff — comments after it (or all if none) are unaddressed. For each unaddressed actionable comment: check out the PR branch, implement the requested change, run `make all`, commit, push, reply to the comment individually (never batch multiple responses into one reply), and return the repo to `main`. Reply mechanism depends on comment type: for review comments (line-level), reply in-thread via `gh api repos/<owner>/<repo>/pulls/<n>/comments/<comment_id>/replies -f body="Done: <summary>"`; for issue comments (general), use `gh pr comment <n> --repo <repo> --body "Done: <summary>"`. Skip non-actionable comments (approvals, "LGTM", emoji). If a comment requests a change you disagree with, reply with `Done: not applied — <reason>`. This step is mandatory before every idea — never skip it, even if the previous idea just opened a PR (the reviewer may have commented in the meantime).
- Process ideas one at a time, in priority order: P1 → P2 → P3 (file order within each priority).
- For each idea, dispatch a single sub-agent (no parallel sub-agents).
- The sub-agent must first read the current source and validate the idea. Reject if the issue no longer exists, is a false positive, the fix would break the public API, or it's too large/risky for one focused PR.
- If valid: ensure the target repo's `main` is up to date (`git checkout main && git pull`), create a `feature/<slug>` branch off `main`, implement the fix, run `make all` in that repo, commit, push, and open a PR (base `main`, assignee `pierrre`, request review from `pierrre`). All new and modified code must be covered by tests — verify with `go test -coverprofile=cover.out ./... && go tool cover -func=cover.out` and add cases until coverage is complete. If a part genuinely cannot be tested, note it in the PR body. Remove the `cover.out` file after checking.
- If rejected: no branch, no PR. Record the reason in the removal commit message (see below).
- After each idea (implemented or rejected): remove its bullet from this file, commit on `go-projects` `main` with a message like `Remove <idea>: implemented` or `Remove <idea>: rejected: <reason>`, and push. No branch or PR for this edit.
- Return the target repo to `main` and pull before starting the next idea.

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

- [bug] (P3) `processStdin` never emits a trailing newline, unlike `processArgs` (which uses `Fprintln`), producing inconsistent/missing final output. Refs: cmd/geohash/geohash.go:107-113
- [improve] (P3) `Decode("")` returns the whole-world `defaultBox` with nil error rather than an "empty geohash" error, surprising callers. Refs: geohash.go:85-108
- [refactor] (P3) `flag.Parse()` runs in `init()`; moving it into `main()` is conventional and would make the CLI testable. Refs: cmd/geohash/geohash.go:60-64
- [docs] (P3) `Encode` doc claims max precision 32 but is silent on `precision <= 0` behavior (panics for negative, returns "" for zero). Refs: geohash.go:49-55

## unlimited-channel — github.com/pierrre/unlimited-channel

Unbounded channel that bridges an input and output channel through a single goroutine and an in-memory linked-list queue, with `Close()` draining both ends.

- [bug] (P3) `Channel.Close` is unsafe for concurrent callers: two goroutines can both hit `default: close(c.in)` before either close lands, double-closing and panicking; the contract is undocumented. Refs: unlimited_channel.go:64-76
- [improve] (P3) `WithContext` passes a context that has no cancellation effect on the goroutine — it only exits via `Close()` or closing the input — so the parameter implies context-driven shutdown that does not exist. Refs: unlimited_channel.go:30-32,166-174
- [improve] (P3) `queue.enqueue` sets `newElem.value` but never `newElem.next`, so list correctness silently depends on `dequeue` nilling `next` before returning elements to the pool; any future change to `dequeue` corrupts the tail. Refs: queue.go:14-27,29-45
- [refactor] (P3) `run()` is a four-flag state machine (`inOpen/inReceived/outValueOK/outSent`) across three cascading selects with multiple `continue` paths, acknowledged by the `gocyclo` suppression and hard to audit for correctness. Refs: unlimited_channel.go:78-145

## githubhook — github.com/pierrre/githubhook

GitHub webhook HTTP handler that parses and HMAC-verifies untrusted payloads, dispatching decoded events via callbacks.

- [improve] (P3) Empty `Secret` silently skips signature verification, so a misconfigured handler accepts forged payloads; fail closed or document loudly. Refs: githubhook.go:113
- [improve] (P3) Distinct error messages ("format" vs "doesn't match secret") are reflected in the HTTP response body via `http.Error`, leaking the verification stage to an attacker. Refs: githubhook.go:124,142-144
- [docs] (P3) README lists "Secret validation" as a feature without noting that omitting `Secret` disables it entirely. Refs: README.md:10

## go-stuff — github.com/pierrre/go-stuff

Shared-module collection of small, self-contained Go experiment sub-packages (five fibonacci variants, a reflection-based function-mocking helper, and a password-entropy estimator).

- [improve] (P3) fibonacci v1/v2/v3 `fibonacciString` keeps a dead local `s` (`s := ""; s += sb.String(); return s`) that just copies the builder's output; return `sb.String()` directly. Refs: fibonacci/v1/fibonacci.go:24,30 (also v2:25,32; v3:28,37)
- [docs] (P3) README.md is a single header line ("# Go stuff") with no overview of the sub-packages or their structure, so a new reader has no map of the repo. Refs: README.md:1
- [refactor] (P3) `call` uses an `else` after a `return` (`if ... { return ... } else { return ... }`), which is redundant; drop the `else` branch. Refs: func-mock/func_mock.go:94

## file-duplicate — github.com/pierrre/file-duplicate

CLI/library that finds duplicate files by walking filesystems, grouping files by identical size, then SHA256-hashing each same-sized group.

- [improve] (P3) No validation that roots are non-empty: invoking the CLI with no path args silently emits nothing and exits 0. Refs: cmd/file-duplicate/flags.go:22
- [improve] (P3) `-v` only has effect together with `-continue-on-error` (it gates the per-error logger); on its own it's a no-op, which is surprising. Refs: cmd/file-duplicate/main.go:56
- [improve] (P3) `filesBySize` holds every qualifying path in memory and each same-sized file is hashed in full, with no partial-hash/early-exit for large groups. Refs: file_duplicate.go:163, file_duplicate.go:192

## file-random — github.com/pierrre/file-random

CLI library `filerandom` walks one or more filesystems, collects regular files above a min size, and `rand.Intn`-picks uniformly; the `cmd/file-random` binary loops printing/opening picks.

- [improve] (P3) `-min-size` accepts negatives with no validation; `WithMinSize(<0)` silently disables the filter since `size < negative` is never true. Refs: cmd/file-random/flags.go:22, file_random.go:121
- [improve] (P3) Omitting roots produces the generic "no file" error instead of a clear "no roots specified", deferring the real cause past the flag layer. Refs: cmd/file-random/flags.go:26, cmd/file-random/main.go:43
- [refactor] (P3) `fl.minSize != 0` uses 0 as an implicit "no filter" sentinel while `newFlags` defaults to 1, making the option's tri-state (unset/0/positive) semantics implicit and surprising. Refs: cmd/file-random/main.go:80
- [docs] (P3) README lists no CLI flags beyond a `-h` hint, so users have no reference for `-min-size`, `-open`, `-loop`, `-continue-on-error`, `-v`. Refs: README.md

## mandelbrot — github.com/pierrre/mandelbrot

Mandelbrot fractal renderer: core set computation with specialized per-power iterators, an image package (sequential/parallel render, colorizers), and a small CLI helper for PNG output.

- [refactor] (P3) The 19 specialized `newPowN` functions (lines 88–556) duplicate an identical loop body differing only in the power computation; a generator or shared helper would remove ~470 lines of near-copy-paste. Refs: mandelbrot.go:88
- [docs] (P3) README.md is only 5 lines (title + pkg.go.dev badge) with no install, usage, CLI, or example sections to orient new users. Refs: README.md:1
- [improve] (P3) `cmd.Save` panics on `png.Encode`/`os.WriteFile` errors instead of returning them, forcing callers to `recover` if they want to handle I/O failures gracefully. Refs: cmd/cmd.go:16

## langton — github.com/pierrre/langton

Langton's ant cellular automaton library with a core grid/ant/rules engine and two CLI frontends (termbox graphical, text stdout).


- [bug] (P3) `Game.Step` indexes `g.Rules[v]` where `v ∈ [0, States)` with no check that `len(Rules) >= States`; a multi-state grid (e.g. `States=4`) paired with `RulesBasic` (2 rules) panics on slice bounds after the first wrap. Refs: langton.go:169
- [improve] (P3) `cmd/text` runs an infinite `for` loop with no exit condition, step limit, or signal handling, so it can only be stopped by killing the process (unlike `cmd/termbox` which exits on keypress). Refs: cmd/text/text.go:25
- [docs] (P3) README is 5 lines with no usage, build, or run instructions and doesn't mention the `cmd/termbox` or `cmd/text` demo commands. Refs: README.md

## cellauto — github.com/pierrre/cellauto

Small cellular automaton library exposing Game of Life and Wireworld rules with two CLI frontends (interactive termbox viewer and PNG-output runner).

- [improve] (P3) `Step(ctx)` accepts a context but never selects on it, so neither library Step nor the unbounded `run` loop in cmd/wireworld can be cancelled. Refs: cellauto/cellauto.go:104, wireworld/wireworld.go:76
- [improve] (P3) cmd/gameoflife's render/step `for` loop has no sleep or frame-rate cap, busy-spinning and pegging one CPU core. Refs: cmd/gameoflife/gameoflife.go:41
- [docs] (P3) README is a 5-line stub with no usage, subcommand listing, or description of the `.wi` file format consumed by cmd/wireworld. Refs: README.md
- [refactor] (P3) `wireworld.Game` re-implements the `tmpGrid` alloc + swap pattern already present in `cellauto.Game`; it could embed or delegate to the shared Game. Refs: wireworld/wireworld.go:69

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
