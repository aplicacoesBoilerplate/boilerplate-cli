# STATE

## Decisions

### AD-001
- **Decision**: Caller Go reutiliza os calleds centrais publicados; a CI da aplicação é um gate anterior, não uma dependência funcional da regra de versão.
- **Reason**: O contrato precisa atender outras aplicações Go sem Cobra, GoReleaser ou lógica local de release.
- **Scope**: Integração piloto na issue #8; não encerra a distribuição de binários.
- **Date**: 2026-10-04
- **Status**: active

## Handoff

- **Feature**: central-versioning-pilot, incremento da issue #8
- **Phase / Task**: T1-T3 complete; independent local verification PASS
- **Completed**: nine criteria, nine caller contracts, real pinned Go 0.0.1/exact SHA, tests/vet/build/lint, four killed mutations; central post-merge CI PASS at f3e7a01
- **In-progress**: none locally
- **Next step**: obtain explicit pilot branch push/PR authorization; feature/issue-8 -> release/v0.0.1 for CI, then separately approved promotion to develop for hosted preview
- **Blockers**: no hosted pilot execution, app merge or Release authorized; local tests do not prove remote called execution
- **Uncommitted files**: prior staged hooks, prior staged version of commit validator, GoReleaser config and local release validator retained; local Go line endings normalized, no remaining unstaged semantic Go diff
- **Branch**: feature/issue-8
- **Central revision**: f3e7a0168da9e3882840302add3ac8d162bac81e (GitHub master)
- **Sensor**: mutations restored and real-tree baseline preserved; temporary directory removal blocked by policy, path recorded in validation.md
