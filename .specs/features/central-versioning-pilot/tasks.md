# Central versioning pilot Tasks

**Status**: Complete - local verification PASS; hosted gate pending authorization

## Test Coverage Matrix

| Layer | Required test | Coverage | Command |
| --- | --- | --- | --- |
| Caller boundary | Parsed YAML contract | Events, refs, inputs, permissions, CI dependency and isolation | python tests/versioning/test_caller.py |
| Native configuration | Real pinned adapter in temporary clone | Stable bootstrap, exact SHA, no tags | python tests/versioning/test_caller.py --real-go |
| Go formatting/CI | Existing application suite and lint | No functional source change | go test ./...; go vet ./...; go build ./...; golangci-lint run ./... |
| Documentation | none | Independent review and relative links | git diff --check |

## Gate Check Commands

Quick: `python tests/versioning/test_caller.py`
Build: Quick, `--real-go`, `go test ./...`, `go vet ./...`, `go build ./...`, `golangci-lint run ./...`.

## Execution Plan

### Phase 1

```text
T1 → T2 → T3
```

## Task Breakdown

### T1: Restore existing formatting gate

**Status**: Complete
**What**: Normalize Go working-copy formatting; Windows CRLF produced successive baseline gofmt findings. Only two files have a versionable formatting diff.
**Where**: `core/root/root.go`, `internal/features/auth/command/login.go`
**Depends on**: none
**Requirement**: PILOT-02.1
**Done when**:

- [x] Existing Go tests, vet, build and lint pass; diff is formatting-only.

**Tests**: existing Go suite and golangci-lint
**Gate**: Build, existing Go suite plus vet/build/lint (caller tests start at T2)
**Commit**: `style(go): restore formatting gate for release pipeline`

**Evidence**: baseline Go tests/vet/build passed; lint failed gofmt. After normalization lint reports zero issues; versionable diff is whitespace alignment and spacing before a composite literal only. All other source blobs unchanged; working-copy line endings were normalized.

### T2: Integrate published calleds with inline Go CI

**Status**: Complete - independent local verification PASS
**What**: Add preview/publish jobs to existing workflow; preserve inline CI and native config, with contract and real bootstrap tests.
**Where**: `.github/workflows/ci.yml`, `go-gitsemver.yml`, `.golangci.yml`, `scripts/validate-commit-msg.sh`, `tests/versioning/test_caller.py`
**Depends on**: T1
**Requirement**: PILOT-01, PILOT-02, PILOT-03.1
**Done when**:

- [x] Parsed contract rejects wrong revision/events/permissions, missing CI dependency and application-coupled versioning.
- [x] Existing configuration gives stable 0.0.1 with exact SHA in a clean clone; no GitHub writes.
- [x] Existing hooks/GoReleaser/local script remain outside this commit and all build gates pass.

**Tests**: caller boundary and native real adapter
**Gate**: Build
**Commit**: `ci(versioning): consume centralized Go release workflows`

**Evidence / adequacy**: nine caller contracts PASS after eight missing-configuration failures before implementation; native real pinned Go PASS at local HEAD, stable 0.0.1/exact SHA/explanation/unchanged tags. Go tests, vet, build and lint PASS. Contract assertions cover revision/inputs (PILOT-01.1), phase guard/read permissions (01.2), push/needs/write permissions (01.3), event lists (01.4), inline CI commands/runtime (02.1), exact allowed versioning job keys (02.2). Native literal checks cover PILOT-03.1. Existing staged hooks and untracked GoReleaser/local validator are excluded via path-specific commit (02.3). Independent mutation review follows T3.

### T3: Document pilot adoption and remaining release scope

**Status**: Complete - independent local verification PASS
**What**: Document caller, native bootstrap and future distribution separately; preserve prior README work.
**Where**: `README.md`, `docs/versioning-pilot.md`
**Depends on**: T2
**Requirement**: PILOT-03.2
**Done when**:

- [x] Reader can distinguish central versioning, application CI and future GoReleaser assets; no hosted/release claim without evidence.
- [x] Independent verifier reports scoped criteria, gates and isolated mutations: nine local ACs PASS, four mutants killed, build gates PASS; no hosted result claimed.

**Tests**: none (human documentation layer)
**Gate**: Build plus link validation
**Commit**: `docs(versioning): explain pilot caller and delivery gates`

## Local delivery result

Independent report `validation.md`: 9/9 local criteria, nine contract tests,
real pinned native adapter, Go tests/vet/build/lint and 4/4 killed mutations.
Remote pilot push/PR and hosted execution remain a separate authorization gate.
The existing staged hooks and untracked release distribution work remain outside
these commits. A restored temporary sensor directory was retained after policy
blocked cleanup; its exact path is in the report.
