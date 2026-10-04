# Central versioning pilot - issue #8

## Problem Statement

O piloto deve consumir a pipeline aprovada no GitHub, sem copiar suas regras
e sem acoplar versionamento a Cobra, binário boilerplate ou GoReleaser.
Escopo aprovado: continuar no worktree #8, consolidar CI/caller, preservar
trabalho existente e iniciar testes antes de promover uma aplicação a master.

## User Stories

### PILOT-01: Caller mínimo e governado

**Acceptance Criteria**:
1. WHEN the caller loads, it SHALL select both central workflows at published SHA `f3e7a0168da9e3882840302add3ac8d162bac81e`, with only adapter `go-gitsemver`, project_path `.` and target_branch `master` as versioning inputs.
2. WHEN a PR or review targets develop or master, the caller SHALL invoke preview with contents/pull-requests/issues/checks read; other PR bases SHALL keep application CI without invoking unsupported preview transitions.
3. WHEN a push targets master, publication SHALL require successful go-ci and contents write with pull-requests/issues read; PR and review events SHALL not publish.
4. WHEN milestone, label, head or review changes, the caller SHALL subscribe to opened/reopened/synchronize/edited/labeled/unlabeled/milestoned/demilestoned PR events and submitted/dismissed review events.

### PILOT-02: CI e isolamento

**Acceptance Criteria**:
1. WHEN application CI runs, it SHALL use go.mod, run existing lint/test/vet/build/help/commit checks and pilot contract tests in the same workflow, without a separate go-ci reusable file.
2. WHEN versioning jobs load, they SHALL contain no local versioning command, GoReleaser dependency, custom token or application-specific build command; only publication SHALL depend on the application's go-ci gate.
3. WHEN the preparation is committed, it SHALL preserve pre-existing staged hooks and untracked GoReleaser/validation work outside the caller deliverable.

### PILOT-03: Configuração e evidência

**Acceptance Criteria**:
1. WHEN the pinned native adapter evaluates the pilot's actual history in a clean temporary clone using the existing configuration, it SHALL return stable version `0.0.1` and the requested SHA without creating tags or Releases.
2. WHEN documentation is read, it SHALL identify the published revision, bootstrap removal after first release, the optional environment and that binary distribution and hosted publication are not yet validated by this local integration.

## Out of Scope

Publicar ou integrar o piloto em master nesta etapa; criar tag/Release real;
implementar GoReleaser, encerrar #8/#14, alterar código funcional da CLI,
copiar scripts centrais ou substituir o histórico OpenSpec existente.

## Assumptions & Open Questions

| Assumption | Chosen default | Rationale |
| --- | --- | --- |
| Trabalho existente | Preservar hooks, GoReleaser e script local de validação | Reuso autorizado do worktree #8 não implica descartar alterações anteriores. |
| Bootstrap | Manter configuração existente com next-version 0.0.1 e incremento Disabled até primeira release | Ainda não há tags publicadas; retirar ambos em entrega posterior rastreável. |
| Testes remotos | Publicar branch/PR somente com autorização explícita | Validação local não equivale a execução hospedada ou autorização de release. |
| Formatação | Normalizar gofmt/terminadores da cópia local; somente dois arquivos têm diff versionável | Gate de qualidade da issue #8 já prevê lint; nenhuma mudança funcional. |

Open questions: none for local implementation. Hosted pilot execution is a subsequent delivery gate.

## Requirement Traceability

| Requirement | Task | Status |
| --- | --- | --- |
| PILOT-01 | T2 | Implemented - verification pending |
| PILOT-02 | T1, T2 | Implemented - verification pending |
| PILOT-03 | T2, T3 | In progress |
