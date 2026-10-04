# Central versioning pilot validation

**Verdict: PASS** for the nine local integration acceptance criteria. Hosted pilot execution, binary distribution, real publication and completion of issue #8 are not established by this report.

**Date**: 2026-10-04
**Verifier**: independent sub-agent, author != verifier
**Spec**: `.specs/features/central-versioning-pilot/spec.md`
**Diff range**: `7e2febf0df2edc4f0fad4f49facc977216f27953..da0d9c60c7e1a9d2908f4e3c9cb53900f5659827`
**Central caller pin**: `f3e7a0168da9e3882840302add3ac8d162bac81e`

## Spec-anchored acceptance criteria

| Criterion | Exact expected outcome and evidence | Result |
| --- | --- | --- |
| PILOT-01.1 | `tests/versioning/test_caller.py:33` asserts both remote workflow names at the literal approved SHA; `:40` asserts the complete input mapping equals adapter `go-gitsemver`, path `.`, target `master`. | PASS |
| PILOT-01.2 | `tests/versioning/test_caller.py:45` asserts the complete PR/review and develop/master guard; `:48` asserts all four read permissions; `:68` preserves application CI for release and pilot PR bases. Read `.github/workflows/ci.yml:15`: application CI excludes review events only. | PASS, local YAML contract |
| PILOT-01.3 | `tests/versioning/test_caller.py:54` asserts push/master guard, `:55` asserts `needs == go-ci`, and `:56` asserts contents write plus PR/issues read. The job has no `always()` override, so its declared dependency requires the CI success path. | PASS, local YAML contract |
| PILOT-01.4 | `tests/versioning/test_caller.py:61` asserts all eight PR event types exactly; `:64` asserts submitted/dismissed reviews. | PASS |
| PILOT-02.1 | `tests/versioning/test_caller.py:72` rejects a reusable go-ci, `:76` asserts go.mod runtime selection, `:79` checks test/vet/build/help/contract commands, `:81` checks lint action, `:82` checks commit validation; `:90` rejects a separate go-ci workflow. Workflow read confirms the inline command loop and Python dependency installation. | PASS, configured CI plus local gates |
| PILOT-02.2 | `tests/versioning/test_caller.py:88` asserts exact allowed job keys, allowing needs only in publish; `:40` restricts inputs. No steps, custom secrets/token, GoReleaser or application build commands exist inside either versioning job. | PASS |
| PILOT-02.3 | Commit-range `git diff --name-only ... -- .githooks .goreleaser.yml scripts/validate-release-version.sh` returns no paths. The existing staged hook entries and untracked release files remain in the real porcelain baseline. Before/after binary index/worktree diffs and SHA256 hashes of the two untracked files match exactly. See preservation record below and documentation `docs/versioning-pilot.md:106`. | PASS, direct Git/hash evidence |
| PILOT-03.1 | `tests/versioning/test_caller.py:104` clones actual pilot history; `:105` selects the requested SHA; `:106` creates master only in scratch; `:107` supplies current native config; `:111` asserts literal SemVer `0.0.1` and exact SHA; `:113` requires explanation; `:115` asserts unchanged tags. Actual result: `0.0.1` at `da0d9c60c7e1a9d2908f4e3c9cb53900f5659827`, explanation present, tags unchanged. | PASS, real native execution |
| PILOT-03.2 | Direct review of `docs/versioning-pilot.md:20` identifies published revision; `:61` describes optional environment; `:78` and `:81` identify removal of both bootstrap options after the first tag; `:85` separates binary distribution; `:152` explicitly limits static/native proof and hosted claims. README links the guide and states the same delivery limitation. All three relative guide links resolve to present files. | PASS, documentation review |

All nine criteria have evidence. No spec-precision gap within this approved local scope. These static contracts assert exact policy configuration; they do not execute GitHub Actions or prove that GitHub accepts the called workflows, permissions or hosted artifacts.

## Build gate and test integrity

All commands were run independently against the real worktree at the cited HEAD:

| Command | Result |
| --- | --- |
| `python tests/versioning/test_caller.py` | 9 tests passed, 0 failures, 0 skips |
| `python tests/versioning/test_caller.py --real-go` | Real pinned adapter returned stable 0.0.1/exact HEAD; explanation present; tags unchanged |
| `go test ./...` | PASS |
| `go vet ./...` | PASS |
| `go build ./...` | PASS |
| `golangci-lint run ./...` | PASS, 0 issues |
| `git diff --check` | PASS; Git reports existing LF/CRLF conversion notices |

The existing Go suite contains one top-level test and three auth command subtests, all passing. Seven packages have no test files; these are not disabled tests. `git grep` on the base and HEAD finds the same top-level test at `core/runtime/execute_test.go:12`; no Go test was removed or weakened in the commit range. Caller coverage increases from zero to nine tests, plus one native integration check. The real-go mode intentionally runs the native check instead of the YAML suite, so both invocations are required and were executed.

## Discrimination sensor

Each mutation was injected independently into isolated copies, with the prior mutation restored before the next. Tests were unchanged.

| Mutation | Location | Observed rejection | Result |
| --- | --- | --- | --- |
| Publication dependency `go-ci` changed to `unrelated-ci` | `.github/workflows/ci.yml:64` | Contract fails at `tests/versioning/test_caller.py:55` | Killed |
| Preview checks permission changed from read to write | `.github/workflows/ci.yml:55` | Exact permission mapping fails at `tests/versioning/test_caller.py:48` | Killed |
| Publish pin changed from approved SHA to mutable master | `.github/workflows/ci.yml:69` | Exact remote workflow reference fails at `tests/versioning/test_caller.py:33` | Killed |
| Native bootstrap next-version changed from 0.0.1 to 0.0.2 | `go-gitsemver.yml:4` | Real native adapter returns 0.0.2 and the literal version assertion fails at `tests/versioning/test_caller.py:112` | Killed |

**Sensor result**: 4 injected, 4 killed, 0 survived. Restoring the copies returns the nine contracts and native execution to PASS. The native clone has no versionable diff after restoration.

Scratch removal with native PowerShell, an explicitly resolved and validated absolute temporary path, was rejected by the automatic tool policy with `blocked by policy`. No alternate destructive operation was attempted. All mutations were restored. The inert scratch remains at `C:/Users/gerso/AppData/Local/Temp/pilot-verifier-aad6441a-ec06-4f07-a4f1-047fab938b74` for later authorized cleanup. Each clone created internally by the real-go test was removed by its TemporaryDirectory context.

## Preservation record

Real-tree porcelain before the sensor and after restoration matched exactly. Before/after `git diff --binary`, `git diff --cached --binary`, and SHA256 hashes for `.goreleaser.yml` and `scripts/validate-release-version.sh` also matched exactly. The baseline includes staged `.githooks/commit-msg`, `.githooks/pre-commit`, and `scripts/validate-commit-msg.sh`; untracked `.goreleaser.yml` and `scripts/validate-release-version.sh`; and fourteen existing Go working-copy entries caused by line-ending normalization. These entries were neither staged nor edited by this verifier. The only real-tree write is this report. No stash, commit, push, PR, tag, Release, merge or other worktree modification occurred.

## Native fixture adequacy and limits

The fixture uses actual repository history and current config, not synthetic commits or a mocked native result. The temporary master ref at the tested HEAD is explicit controlled setup and does not prove promotion eligibility or hosted branch policy. Changing the real config to 0.0.2 kills the native test, demonstrating that its expected version is independent of the configured value and that config is actually consumed. The SHA assertion is literal requested HEAD equality. The no-tags assertion compares the whole tag list before and after calculation.

The binary metadata inspected independently identifies `github.com/MyCarrier-DevOps/go-gitsemver` with pseudo-version ending `680c1c12d9a4`. The test's metadata guard checks that short revision substring; it is a provenance sanity check, not cryptographic binary attestation. Nonempty stderr establishes explanation presence only, not explanation semantics. No Release API is invoked by the local native calculation, and no GitHub write was performed; this does not test hosted Release authorization or publication behavior. These limits do not contradict the local criterion.

## Quality and completion

The two committed Go diffs only align struct fields/remove spacing before a struct literal. No functional CLI behavior changed. The workflow keeps existing application CI inline and adds minimal caller jobs. Tests each map to a scoped requirement; documentation clearly separates future distribution and hosted delivery. No repository AGENTS quality guide was found by the author; the supplied project TLC criteria and existing application gates were applied. The OpenSpec worktree is outside this verification and untouched.

T1, T2 and T3 deliverables are independently verified for this local scope. T3's verification checkbox and requirement traceability updates are left to the author, because the verifier was authorized to write only this file. There are no ranked implementation gaps and no grounded failure lesson to record. Cleanup of the retained inert scratch is an operational follow-up.

**Overall**: PASS, ready for the subsequent authorized hosted pilot gate. This report does not close issue #8 or issue #14 and does not approve a publication or merge.
