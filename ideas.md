# Project ideas

# Processing

- Process ideas one at a time, in priority order: P1 → P2 → P3 (file order within each priority).
- For each idea, dispatch a single sub-agent (no parallel sub-agents).
- The sub-agent must first read the current source and validate the idea. Reject if the issue no longer exists, is a false positive, the fix would break the public API, or it's too large/risky for one focused PR.
- If valid: create a `feature/<slug>` branch off `main`, implement the fix, run `make all` in that repo, commit, push, and open a PR (base `main`, assignee `pierrre`).
- If rejected: no branch, no PR. Record the reason in the removal commit message (see below).
- After each idea (implemented or rejected): remove its bullet from this file, commit on `go-projects` `main` with a message like `Remove <idea>: implemented` or `Remove <idea>: rejected: <reason>`, and push. No branch or PR for this edit.
- Return the target repo to `main` before starting the next idea.

Generated: 2026-08-09
Scope: /home/pierre/Git/pierrre workspace (Go 1.26.0)
Method: direct review of each project's key files (golang-review + golang-new skills applied).

Legend:
  Tags:     [bug] [feat] [refactor] [improve] [test] [docs]
  Priority: P1 high  P2 medium  P3 low

## assert — github.com/pierrre/assert

Go test assertion library using generics (no reflection), with auto-updating snapshot assertions (`assertauto`).

- [refactor] (P2) `ValueStringer.Load()` is called multiple times per failure (e.g. twice in `Equal`); hoist to a local `vs := ValueStringer.Load()` once per Fail. The double-load pattern repeats across every assertion file. Refs: equal.go:21, slice.go:109,145, map.go:55,109, deep_equal.go:43
- [improve] (P2) `bytesWriterPool` uses `MaxCap: -1` (unbounded); pooled buffers can retain arbitrarily large stack-trace strings, pinning memory across runs. Cap `MaxCap` to bound retention. Refs: assert.go:77-79
- [improve] (P3) `Fail` allocates `args := []any{msg}` then spreads it into `o.report(tb, args...)`; since there is exactly one arg, call `o.report(tb, msg)` directly to drop the slice allocation. Refs: assert.go:73-74
- [feat] (P2) `assertauto` documents that concurrent calls sharing a test name are unsafe and produce flaky results; a per-test-name `sync.Mutex` (keyed by test name) would make this safe automatically. Refs: assertauto/assertauto.go:17-20
- [feat] (P2) No unordered slice comparison (testify's `ElementsMatch` equivalent); `SliceEqual` is order-sensitive. A `SliceElementsMatch[S ~[]E, E comparable]` would fill the gap. Refs: slice.go
- [feat] (P3) `chan.go` TODO is unresolved: no support for receive/send-only channels (`<-chan T`, `chan<- T`). Refs: chan.go:8
- [improve] (P3) `AllocsPerRun` silently passes (returns true) under `-race`, which can hide allocation regressions in race-enabled CI; consider `tb.Log`-ing that the check was skipped so it's visible. Refs: alloc.go:17-18
- [docs] (P3) README documents core assertions but has no section for `assertauto` (auto-updating snapshot assertions), a notable differentiating feature. Refs: README.md

## errors — github.com/pierrre/errors

Errors library with message wrapping, stack traces, tags, values, verbose output, and a drop-in std `errors` replacement; composable sub-packages (errbase, errmsg, errstack, errtag, errval, errtmp, errignore, erriter).

- [improve] (P2) `errors.go` imports `testing` and calls `testing.Testing()` in `init()` to decide whether to panic on global-init detection. This pulls the `testing` package into every consumer's production binary (binary bloat). The escape hatch (`errbase.New`) already exists; consider documenting the tradeoff or gating detection behind a build tag / an opt-in atomic. Refs: errors.go:10,70-78
- [docs] (P3) `errtmp.Is` returns `true` by default (errors without a `Temporary()` method are treated as temporary), which inverts the stdlib convention where absence means "not temporary"; document this explicitly or reconsider the default. Refs: errtmp/errtmp.go:46-54
- [refactor] (P3) `errtag.Get` and `errval.Get` duplicate a "first key wins" map-build loop over `All(...)`; extract a shared helper in `erriter` to reduce duplication. Refs: errtag/errtag.go:89-98, errval/errval.go:84-93
- [feat] (P3) `errval` only exposes `Get` (returns a full `map[string]any`); a `GetValue(err, key) (any, bool)` single-value helper would be more convenient and allocate less. Refs: errval/errval.go:84
- [feat] (P3) `errtag` ships typed helpers only for int/int64/float64/bool; consider a generic `WrapTag[T constraints.Integer|Float|~bool]` or drop the helpers in favor of `errval` for non-string values. Refs: errtag/errtag.go:29-47
- [test] (P3) `erriter.iterFunc` interleaves linear (`Unwrap() error`) and multi (`Unwrap() []error`) recursion; add a test that locks the traversal order for mixed Join+Wrap chains. Refs: erriter/erriter.go:17-33
- [docs] (P3) README "Extend" lists sub-packages but omits `erriter` (public error-tree iteration API) and `errverbose`. Refs: README.md:73-81

## pretty — github.com/pierrre/pretty

Reflection-based pretty printer with a modular `ValueWriter` chain, cycle/recursion detection, max depth, hex dumps, iter.Seq/Seq2 support, and protobuf extensions.

- [improve] (P2) `RecursionWriter.checkRecursion` does `slices.Contains(st.Visited, e)` for every pointer/map/slice node, so printing a deeply nested graph (e.g. a long linked list) is O(n²). Replace the slice with a hash set (`map[VisitedEntry]struct{}` or a specialized set) for O(1) membership. Refs: recursion.go:50
- [improve] (P3) `State.release()` resets `Writer` but not `Visited`; a pooled `*State` that once printed a large graph retains the oversized `Visited` backing array indefinitely. Truncate to `nil` (or cap) on release. Refs: state.go:49-52,38
- [improve] (P3) `common.go` imports `testing` and calls `testing.Testing()` in `init()` (same pattern as the errors package), pulling the `testing` package into consumer production binaries. Refs: common.go:6,15-19
- [improve] (P3) `reflectValuePools` (`syncutil.Map[reflect.Type, *Pool]`) in map.go grows unbounded — every distinct map key/elem type creates a permanent pool that is never evicted, leaking memory in processes that encounter many map types. Refs: map.go:94-109
- [feat] (P3) No global output-size cap: only per-collection `MaxLen` and `MaxDepth` exist. A `MaxBytes`/`MaxRunes` limit on the `State` would make it safer for logging very large values. Refs: state.go, common.go
- [docs] (P3) README is 27 lines for a feature-rich library; doesn't mention `RecursionWriter`, `Filter`, `UnwrapInterface`, `WeakPointer`, `Iter`/`Range` writers, or the `ext/` packages. Refs: README.md

## vld — github.com/pierrre/vld

Generics-based, allocation-free validation library with composable validators (And/Or/If/IfElse/Switch), built-in checks, localized messages (en/fr), and error-path extraction.

- [improve] (P2) `MessageValidator.Validate` discards the underlying error and returns a `MessageError` with no `Unwrap`, so `GetErrorPath`, `errors.Is`, and `errors.As` lose all path/context when `Message` wraps a path-bearing validator (e.g. `FieldValidator`/`SliceEach`). Consider wrapping the original error so the chain survives while the displayed message is overridden. Refs: transform.go:200
- [improve] (P3) `GetErrorPath` follows a single chain (`errors.AsType` then `err = pErr.Err`); for joined errors (multi-branch, as produced by `SliceEach`/`ErrorJoin`) it returns only the first branch's path and silently drops the others, which can mislead users debugging nested-slice/map errors. Refs: path.go:112
- [docs] (P3) Constructors `If`, `IfElse`, `Case`, `Parse`, `Get`, `Field`, `SliceUniqueBy` document "It panics if X is nil" but perform no explicit nil check, so the panic actually fires at `Validate` time as a nil-function-call panic, far from construction and harder to diagnose — either add an explicit check at construction or reword the doc to say "panics at Validate time." Refs: condition.go:9
- [docs] (P3) `SwitchValidator.Validate` silently returns nil when no `Case` condition matches (no default), so an exhaustive switch over a finite value set can let invalid values pass unnoticed; document this explicitly or add a `Default`/required-match option. Refs: condition.go:111
- [refactor] (P3) `SliceUniqueValidator.Validate` and `SliceUniqueByValidator.Validate` duplicate the same seen-map/append-`PathElemError` logic; extract a shared `validateUniqueByKey(s, getKey, makeErr)` helper to halve the code and keep the two paths in sync. Refs: slice.go:321

## go-libs — github.com/pierrre/go-libs

Shared utility sub-packages (bytesutil, chansutil, goroutine, reflectutil, singleflight, syncutil/atomicutil, unsafeio, weakutil, etc.) used across the author's projects.

- [refactor] (P3) In `singleflight.doCall`, the block `if !c.doneInitialized { c.doneInitialized = true }` is dead code: the key is deleted under `mu` earlier, so no future `getOrCreateCall` can observe this flag for this call, and the doer's close-decision checks `c.done != nil` not the flag. Refs: singleflight/singleflight.go:105
- [improve] (P3) In `singleflight.waitCall`, `if c.done != nil` is always true and the nil-channel fall-through is unreachable: the first waiter creates `c.done` before any waiter enters `waitCall`, and the doer deletes the key before setting `doneInitialized`, so no waiter can observe `doneInitialized==true && done==nil`. Refs: singleflight/singleflight.go:81

## geohash — github.com/pierrre/geohash

Geohash encode/decode library with a CLI front-end (`cmd/geohash`).

- [bug] (P2) `Encode` only clamps precision above `encodeMaxPrecision`, so a negative `precision` (public API) panics on `buf[:precision]`. Refs: geohash.go:53-55,81
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

- [improve] (P2) `io.ReadAll(req.Body)` reads untrusted bodies with no size limit before signature verification, enabling memory-exhaustion DoS; wrap with `http.MaxBytesReader`. Refs: githubhook.go:86
- [feat] (P2) Only the legacy SHA1 `X-Hub-Signature` is verified; GitHub now recommends SHA-256 via `X-Hub-Signature-256`. Refs: githubhook.go:116
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

- [bug] (P2) Passing `-min-size 0` is silently ignored: the `fl.minSize != 0` guard skips `WithMinSize`, so the library default (1) applies and zero-byte files can never be scanned. Refs: cmd/file-duplicate/main.go:53
- [improve] (P2) `ctx` is threaded into the walk/hash functions but never consulted, so long scans can't be cancelled — the CLI uses `context.Background`, so Ctrl+C only works via signal kill. Refs: file_duplicate.go:133, file_duplicate.go:168
- [improve] (P3) No validation that roots are non-empty: invoking the CLI with no path args silently emits nothing and exits 0. Refs: cmd/file-duplicate/flags.go:22
- [improve] (P3) `-v` only has effect together with `-continue-on-error` (it gates the per-error logger); on its own it's a no-op, which is surprising. Refs: cmd/file-duplicate/main.go:56
- [improve] (P3) `filesBySize` holds every qualifying path in memory and each same-sized file is hashed in full, with no partial-hash/early-exit for large groups. Refs: file_duplicate.go:163, file_duplicate.go:192

## file-random — github.com/pierrre/file-random

CLI library `filerandom` walks one or more filesystems, collects regular files above a min size, and `rand.Intn`-picks uniformly; the `cmd/file-random` binary loops printing/opening picks.

- [bug] (P2) Root "/" is rewritten to "" before `os.DirFS`, which then scans the process CWD (not `/`): DirFS opens via `filepath.Join(dir, name)`, so an empty dir resolves relative to cwd, and the printed path is relative too. Refs: cmd/file-random/main.go:74-77
- [improve] (P3) `-min-size` accepts negatives with no validation; `WithMinSize(<0)` silently disables the filter since `size < negative` is never true. Refs: cmd/file-random/flags.go:22, file_random.go:121
- [improve] (P3) Omitting roots produces the generic "no file" error instead of a clear "no roots specified", deferring the real cause past the flag layer. Refs: cmd/file-random/flags.go:26, cmd/file-random/main.go:43
- [refactor] (P3) `fl.minSize != 0` uses 0 as an implicit "no filter" sentinel while `newFlags` defaults to 1, making the option's tri-state (unset/0/positive) semantics implicit and surprising. Refs: cmd/file-random/main.go:80
- [docs] (P3) README lists no CLI flags beyond a `-h` hint, so users have no reference for `-min-size`, `-open`, `-loop`, `-continue-on-error`, `-v`. Refs: README.md

## mandelbrot — github.com/pierrre/mandelbrot

Mandelbrot fractal renderer: core set computation with specialized per-power iterators, an image package (sequential/parallel render, colorizers), and a small CLI helper for PNG output.

- [bug] (P2) `ColorsIterColorizer` indexes `cols[(res.Iter+shift)%len(cols)]` — panics with integer divide by zero when `cols` is empty, and panics with negative index when `shift<0` since Go's `%` keeps the dividend's sign. Refs: image/colorizer.go:22
- [refactor] (P3) The 19 specialized `newPowN` functions (lines 88–556) duplicate an identical loop body differing only in the power computation; a generator or shared helper would remove ~470 lines of near-copy-paste. Refs: mandelbrot.go:88
- [docs] (P3) README.md is only 5 lines (title + pkg.go.dev badge) with no install, usage, CLI, or example sections to orient new users. Refs: README.md:1
- [improve] (P3) `cmd.Save` panics on `png.Encode`/`os.WriteFile` errors instead of returning them, forcing callers to `recover` if they want to handle I/O failures gracefully. Refs: cmd/cmd.go:16

## langton — github.com/pierrre/langton

Langton's ant cellular automaton library with a core grid/ant/rules engine and two CLI frontends (termbox graphical, text stdout).

- [bug] (P2) The termbox `PollEvent` goroutine ignores `ctx` and never exits; on keypress the deferred `termbox.Close()` runs concurrently with an active `PollEvent`, then the goroutine blocks forever on the unbuffered `evQueue`. Refs: cmd/termbox/termbox.go:21-25
- [bug] (P3) `Game.Step` indexes `g.Rules[v]` where `v ∈ [0, States)` with no check that `len(Rules) >= States`; a multi-state grid (e.g. `States=4`) paired with `RulesBasic` (2 rules) panics on slice bounds after the first wrap. Refs: langton.go:169
- [improve] (P3) `cmd/text` runs an infinite `for` loop with no exit condition, step limit, or signal handling, so it can only be stopped by killing the process (unlike `cmd/termbox` which exits on keypress). Refs: cmd/text/text.go:25
- [docs] (P3) README is 5 lines with no usage, build, or run instructions and doesn't mention the `cmd/termbox` or `cmd/text` demo commands. Refs: README.md

## cellauto — github.com/pierrre/cellauto

Small cellular automaton library exposing Game of Life and Wireworld rules with two CLI frontends (interactive termbox viewer and PNG-output runner).

- [improve] (P3) `Step(ctx)` accepts a context but never selects on it, so neither library Step nor the unbounded `run` loop in cmd/wireworld can be cancelled. Refs: cellauto/cellauto.go:104, wireworld/wireworld.go:76
- [improve] (P3) cmd/gameoflife's render/step `for` loop has no sleep or frame-rate cap, busy-spinning and pegging one CPU core. Refs: cmd/gameoflife/gameoflife.go:41
- [docs] (P3) README is a 5-line stub with no usage, subcommand listing, or description of the `.wi` file format consumed by cmd/wireworld. Refs: README.md
- [refactor] (P3) `wireworld.Game` re-implements the `tmpGrid` alloc + swap pattern already present in `cellauto.Game`; it could embed or delegate to the shared Game. Refs: wireworld/wireworld.go:69

## go-projects — github.com/pierrre/go-projects

Stub Go package with a project name list and a Makefile that fans out commands across sibling repos.

- [feat] (P2) Expose project metadata (name + repo URL pattern) as a Go API, mirroring the GIT_REPOSITORY_PATTERN logic already encoded in the Makefile. Refs: Makefile:13
- [docs] (P1) Expand README.md beyond its single title line to describe purpose, the projects.txt format, and Makefile usage. Refs: README.md:1
- [improve] (P2) Makefile parses projects.txt via `cat | xargs -I {} sh -c` (fragile: no quoting, fails on blanks/special chars); prefer a Go-based iterator or `while read` loop. Refs: Makefile:5
- [improve] (P3) `all-copy-common` blindly overwrites shared config (.golangci.yml, .github, etc.) in sibling repos with no diff guard, risking clobbered local edits. Refs: Makefile:22

## agents — (opencode agent config repo, not Go)

Personal opencode agent-config repo: global AGENTS.md rules, opencode.jsonc provider/model config, and a Makefile that symlinks the repo into ~/.config/opencode.

- [improve] (P3) Commented-out `model` line is dead config with no default model set; either set it or remove it. Refs: opencode/opencode.jsonc:70
- [improve] (P3) `$(CURDIR)` is unquoted in symlink commands; paths with spaces would break installation. Refs: Makefile:3,5
- [improve] (P3) Scaleway baseURL hardcodes a project UUID, coupling config to one project and leaking org identity. Refs: opencode/opencode.jsonc:45
- [refactor] (P3) glm-5.2 model block (limits/cost/variants) is duplicated verbatim across Codix-Iliad and Scaleway providers; factor to avoid drift. Refs: opencode/opencode.jsonc:23,49
- [improve] (P3) `goimports` runs `golang.org/x/tools/cmd/goimports@latest`, unpinned and network-dependent on cache miss. Refs: opencode/opencode.jsonc:102

## di — github.com/pierrre/di

Small generic dependency-injection container with named services, lazy singleton init, and panic-aware builders.

- [bug] (P2) `serviceWrapper.close` does not recover panics (unlike `ensureInitialized` which defers `recoverPanicToError`), so a panicking `Close` aborts `Container.Close` mid-loop and leaks all remaining services. Refs: service.go:89
- [improve] (P3) `Container.Close` closes services in alphabetical key order, not reverse-build/dependency order; a Close that touches a dependency finds it already closed, and a subsequent Get rebuilds it outside the snapshot (leak). Refs: container.go:54
- [docs] (P3) `Container.Close` doc omits close-order semantics and does not warn that Close callbacks must not call Get on other services. Refs: container.go:47

## compare — github.com/pierrre/compare

Generics-supported, reflect-based recursive comparator producing path-annotated Differences with cycle detection, custom Func dispatch, and cached .Equal()/.Cmp() method handling.

- [improve] (P3) `comparePointer` skips `compareNil`, so nil-vs-non-nil `*T` reports "only one is valid" (via `Elem` on the nil pointer) instead of "only one is nil" like chan/func/interface — inconsistent and misleading. Refs: compare.go:311-320
- [improve] (P3) Identifier typo: `cmdMethodFuncs`/`cmdMethodFuncsLock` should be `cmp...` (the sibling func/msg consts all use `Cmp`); hurts grepability, copy-paste from the Equal variant. Refs: compare.go:665-666
- [refactor] (P3) `getMethodEqualFunc` and `getMethodCmpFunc` are ~30-line near-duplicates (nil-sentinel cache lookup + signature validation); collapse into one shared helper parameterized by name/outType. Refs: compare.go:619-691
