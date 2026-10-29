# Versionamento centralizado: etapa 1

O caller de produção está em .github/workflows/release.yml. Ele chama o workflow central e seus scripts na mesma revisão imutável 1d874d946e5f5ad233a6b1f6bf278fef78f3b46d. A aplicação mantém sua CI em ci.yml, reutilizável via workflow_call. Não copiamos scripts nem implementamos cálculo de versão neste CLI.

## Dois momentos

1. PR funcional (develop ou hotfix para master), com milestone SemVer, revisão e merge humano: a CI do commit integrado roda primeiro; depois o called seleciona a estratégia Go do Release Please e prepara o PR técnico com a versão da milestone e o changelog.
2. PR técnico release-please para master: após revisão/CI/merge humano, a pipeline cria tag e GitHub Release automaticamente. Não exige outra milestone nem repete o override.

Só PR closed com merged=true e destino master é elegível. Fechar sem merge, PRs para develop/release e pushes isolados não publicam. O job validate tem leitura; release depende dele e recebe contents/pull-requests/issues write. Não há secrets de registry herdados. Os dois momentos continuam sujeitos às políticas de Actions do consumidor.

## Migração e distribuição

A prévia e o publisher antigos go-gitsemver foram retirados da CI ativa. go-gitsemver.yml e scripts/validate-release-version.sh permanecem históricos, sem execução pela pipeline. As regras antigas de milestone opcional, label de override e environment da implementação Shell não fazem parte deste novo fluxo Release Please.

publish-artifacts=false impede GoReleaser na etapa 1. A configuração .goreleaser.yml, hooks e código funcional permanecem preservados para trabalho posterior. Versão embutida em Cobra, assets/binários, checksums e publicação em registry pertencem às issues da aplicação e à etapa 2; não foram implementados aqui.

## Verificação T6

A suíte Python existente mantém nove testes de contrato, agora alinhados ao fluxo aprovado: SHA igual, tipo Go/flag false, eventos/merge, CI reutilizável somente leitura, dependência de CI/permissões, aposentadoria dos escritores antigos, branches, comandos Go e separação de responsabilidades. Não replica o cálculo do Release Please; esse motor é validado na suíte central TS.

Gates locais aprovados: nove contratos Python, go test ./..., go vet ./..., go build ./..., golangci-lint run ./... (zero issues), build e --help. O antigo modo opcional --real-go executava go-gitsemver e foi aposentado com a integração correspondente, sem reduzir os nove casos de contrato. Os testes Go de aplicação permanecem intactos.

O lint local identificou apenas CRLF introduzido pelo checkout Windows. Foi comprovada igualdade dos fontes normalizados com os blobs HEAD e normalizado esse checkout; git diff --exit-code -- '*.go' terminou sem alteração. Nenhum fonte Go será incluído no commit. T6 tem implementação/gates concluídos e é persistida antes da verificação final.

A revisão independente encontrou um desvio no checkout legado: no evento closed/merged a CI ainda selecionava head.sha do PR. Corrigido para usar merge_commit_sha após merge e preservar head.sha em PR aberto. O contrato existente foi fortalecido, falhou antes da correção e passou depois; continuam nove casos. Assim CI e writer validam/executam o mesmo commit integrado.

Integração em codex/release-please-stage1 baseada em release/v0.0.1. Nenhuma feature foi promovida a master nem alterada. O ciclo real do consumidor continua pendente da revisão e dos merges humanos: levar a integração pelo Git Flow, realizar entrega fix:/feat: com milestone e inspecionar o PR técnico antes de integrar. Conferir tag/Release no SHA desse segundo merge, sem assets publicados.

[Contrato central](https://github.com/aplicacoesBoilerplate/.github/blob/1d874d946e5f5ad233a6b1f6bf278fef78f3b46d/docs/pipelines/versionamento-release-please.md).
