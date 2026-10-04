# Piloto do versionamento central - issue #8

## Realização

A [issue #8](https://github.com/aplicacoesBoilerplate/boilerplate-cli/issues/8)
trata de pipeline, versão e publicação de binários multiplataforma. Este
incremento conecta o versionamento ao MVP central sem assumir que a entrega
de assets GoReleaser também esteja concluída. A issue #14 continua responsável
pela versão incorporada ao executável.

O [workflow único do piloto](../.github/workflows/ci.yml) contém:

| Job | Responsabilidade |
| --- | --- |
| go-ci | Testes, vet, build, lint, help e validação de commits/contrato da aplicação. |
| validate-pr | Chamar prévia central em PR/review com destino develop/master. |
| publish | Após CI e push em master, chamar publicação central protegida. |

Os calleds estão em `aplicacoesBoilerplate/.github`, revisão imutável
`f3e7a0168da9e3882840302add3ac8d162bac81e`, integrada pelo
[PR central #10](https://github.com/aplicacoesBoilerplate/.github/pull/10).
Não usamos arquivos do checkout local do central, `@master` mutável, comandos
locais para escolher versão ou um token pessoal no caller.

Os únicos inputs de versionamento atuais são:

```yaml
with:
  adapter: go-gitsemver
  project_path: .
  target_branch: master
```

Para outra aplicação Go, adaptar caminho e branch, manter o mesmo SHA nos dois
calleds e ligar publicação à sua própria CI. Os comandos específicos de
`boilerplate` pertencem somente ao job de CI deste piloto. O versionamento não
requer Cobra, variáveis de build, GoReleaser, hooks ou nome de executável.
Git, Bash, Node, GitHub CLI e o adaptador são infraestrutura do called, não
dependências funcionais do módulo Go. Consulte os
[scripts centrais](https://github.com/aplicacoesBoilerplate/.github/blob/f3e7a0168da9e3882840302add3ac8d162bac81e/scripts/versioning/README.md).

### PRs e integração

PRs de feature para `release/*` executam CI sem chamar uma prévia cujo fluxo
final não se aplica. Na promoção para develop, o central produz o guia
Markdown/JSON derivado do relatório nativo, sem commitar changelog.
`develop -> master` e `hotfix/<nome> -> master` compartilham todas as regras.
O hotfix parte da master publicada; não existe PATCH forçado no caller.

Milestone ausente aceita o cálculo nativo. Milestone presente deve ter título
estrito `vMAJOR.MINOR.PATCH` e coincidir com versão/incremento. Título inválido
bloqueia sem override. Divergência válida exige label `versioning:override`
por Maintain/Admin e aprovação posterior de outro Maintain/Admin no head atual.
Commit, milestone, label e review mudados reexecutam a avaliação.

A publicação usa `needs: go-ci`, evento push em master e token com contents
write; prévia é somente leitura, incluindo checks. O central revalida branch
padrão, SHA remoto e associação ao único PR integrado. Review se refere ao head
do PR; tag e Release se referem ao SHA integrado. Reexecução idêntica não escreve
novamente; conflito não move tag nem sobrescreve Release.

Environment adicional é opcional. Para exigi-lo, adicionar ao `with` de publish:

```yaml
publication_environment: production
```

As regras do environment e os rulesets/required checks devem ser configurados
no consumidor. O caller não cria essas proteções. Environment não substitui CI
ou autorização de divergência. Não há aprovação por LLM nesta integração.

### Bootstrap do adaptador

O arquivo existente [go-gitsemver.yml](../go-gitsemver.yml) é autodetectado pela
revisão fixada do adaptador. Não adicionar outro `GitVersion.yml` sem revisar
a precedência de busca, pois isso pode ocultar esta configuração.

`next-version: 0.0.1` e `commit-message-incrementing: Disabled` são exceções
temporárias aprovadas para bootstrap. O central usa `0.0.0` como base
matemática na ausência de tags e só aceita primeira versão estável `0.0.1`.
Após a primeira tag, retirar **ambas** as opções em uma alteração rastreável e
verificar increments convencionais antes da próxima publicação. Não manter
`next-version` permanente nem usar override para esconder configuração errada.

### Distribuição da CLI

O central cria tag/Release e suas notas; não gera assets GoReleaser ou injeta
versão em Cobra. A configuração GoReleaser existente permanece separada e
fora dos jobs de versionamento. Sua integração futura deve consumir a versão
publicada, sem recalcular versão ou competir pela criação da mesma Release.

Não depender apenas de push de tag feito com `GITHUB_TOKEN` para disparar um
segundo workflow de assets: esse encadeamento precisa ser projetado e validado
separadamente. Nenhuma distribuição real é autorizada só pelo teste do caller.

## Fontes modificados

- `.github/workflows/ci.yml`: CI inline e chamadas remotas versionadas.
- `tests/versioning/test_caller.py`: contrato YAML e cálculo real sem publicação.
- `go-gitsemver.yml`, `.golangci.yml`, `scripts/validate-commit-msg.sh`:
  preparação existente necessária à CI, incorporada sem mudar seu contrato.
- `core/root/root.go`, `internal/features/auth/command/login.go`: gofmt apenas.
- `README.md` e este guia: adoção e limites.
- `.specs/features/central-versioning-pilot/`: escopo e tarefas TLC locais.

Hooks previamente staged, `.goreleaser.yml` e o auxiliar local de validação de
release foram preservados, sem incorporar essa entrega ao commit do caller.
Não alteramos o worktree da issue #14 nem descartamos o planejamento OpenSpec.

## p/ teste

Pré-requisito do teste de contrato: Python 3.12 com PyYAML 6.0.1. Os testes
Go seguem a versão de `go.mod`. O workflow instala Python e a dependência para
o contrato; o adaptador é instalado pelo central na execução do called.

```bash
python -m pip install PyYAML==6.0.1
python tests/versioning/test_caller.py
go test ./...
go vet ./...
go build ./...
golangci-lint run ./...
```

Para o cálculo nativo real, sem token nem publicação:

```bash
go install github.com/MyCarrier-DevOps/go-gitsemver@680c1c12d9a4f573a8da1b2e3ccebb3571b1cab6
python tests/versioning/test_caller.py --real-go
```

O binário deve estar no PATH. O teste usa histórico real do piloto em clone
temporário, configuração nativa existente, SHA exato, versão `0.0.1`,
explicação e tags inalteradas; descarta o clone ao terminar. Criar nele uma ref
master no SHA avaliado é preparação isolada, não promoção da branch real.
O clone evita a falha de resolução de branch observada na biblioteca nativa
ao executar diretamente neste linked worktree Windows.

Em Windows, git/autocrlf pode introduzir diferenças de terminadores que o
gofmt/lint detecta. Normalizar com gofmt antes do lint; isso não muda regras
funcionais. Os dois diffs versionáveis desta integração são apenas formatação.

### Gate hospedado seguinte

Após autorização de push/PR, a entrega volta para `release/v0.0.1`, preservando
o fluxo da issue. Esse PR testa a CI do piloto, mas não deve publicar nem
chamar a prévia final. Uma promoção controlada para develop permite conferir
o guia e o contrato remoto. Publicação de aplicação exige promoção separada,
aprovada e integrada em master; não fazer merge para produção só para obter
uma evidência de teste.

Contrato estático e cálculo nativo local não provam que o GitHub aceitou o
called, suas permissões ou os artifacts hospedados. O relatório de validação
deve distinguir esses resultados; nenhum resultado remoto está presumido.

## O que há de novo

O piloto utiliza a revisão já publicada da mesma política central, sem copiar
scripts nem manter algoritmo próprio de release. A CI da aplicação continua
substituível, e distribuição de binários pode evoluir sem acoplar o cálculo de
versão a este CLI.
