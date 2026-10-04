# boilerplate-cli

CLI em Go para iniciar e manter aplicações Java e Vue integradas ao ecossistema `aplicacoesBoilerplate`.

O repositório e os templates de origem permanecem públicos. Projetos gerados para outras organizações devem ser privados por padrão. A leitura de artefatos no GitHub Packages continua exigindo autenticação, mesmo quando o repositório de origem é público.

## Estado atual

O projeto está no milestone [`v0.0.1`](https://github.com/aplicacoesBoilerplate/boilerplate-cli/milestone/1). O código atual é somente o scaffold inicial; os comandos ainda exibem mensagens `TODO` e não devem ser usados para modificar projetos reais.

- A [issue #2](https://github.com/aplicacoesBoilerplate/boilerplate-cli/issues/2) define `auth`, `init`, `new`, `add`, `update`, `doctor` e `audit`, incluindo `--dry-run`, idempotência, códigos de saída estáveis e suporte a monorepos.
- A [issue #3](https://github.com/aplicacoesBoilerplate/boilerplate-cli/issues/3) implementa autenticação por `gh auth token` e edição segura das configurações Maven/npm, sem registrar credenciais.
- A [issue #1](https://github.com/aplicacoesBoilerplate/boilerplate-cli/issues/1) acompanha a entrega completa da primeira versão.

## Fluxo planejado

Para pessoas desenvolvedoras que já possuem acesso aos packages da organização:

```text
boilerplate auth login
boilerplate auth status
boilerplate init
boilerplate new [java|vue] <nome>
boilerplate add [java|vue] <package>@<versao>
boilerplate update
boilerplate doctor
boilerplate audit
```

O fluxo de autenticação reutilizará a sessão do GitHub CLI. Não haverá OAuth Device Flow próprio nem armazenamento de uma cópia do token pelo `boilerplate-cli`.

## Arquitetura

A CLI é organizada em Vertical Slices por recurso. Cada slice concentra seus contratos, camada de aplicação, adapters de comando e, quando necessários, domínio e infraestrutura.

```text
internal/features/<recurso>/
  application/    # Orquestra os casos de uso.
  command/        # Adapter Cobra: flags, argumentos e saída.
  contracts/      # Contratos exclusivos do recurso.
  domain/         # Regras de domínio, quando existirem.
  infrastructure/ # Integrações externas, quando existirem.
```

`core` contém somente o runtime, o comando raiz e o registro da árvore Cobra. Contratos realmente compartilhados entre recursos ficam em `internal/shared`.

## Desenvolvimento

Requisitos:

- Go 1.26;
- acesso aos repositórios públicos da organização;
- GitHub CLI autenticado para os futuros fluxos que consultam ou consomem GitHub Packages.

```bash
go test ./...
go vet ./...
go build ./...
```

A CLI usa [Cobra](https://github.com/spf13/cobra). Viper não é dependência do módulo enquanto não existir uma necessidade concreta de configuração adicional.

## Integrações relacionadas

- [`PackagesJava#1`](https://github.com/aplicacoesBoilerplate/PackagesJava/issues/1): estrutura base dos packages Maven;
- [`PackagesJava#5`](https://github.com/aplicacoesBoilerplate/PackagesJava/issues/5): governança e observabilidade;
- [`PackagesJava#6`](https://github.com/aplicacoesBoilerplate/PackagesJava/issues/6): pipeline e publicação Maven;
- [`PackagesJava#9`](https://github.com/aplicacoesBoilerplate/PackagesJava/issues/9): Dependency Graph e Dependabot;
- [`PackagesVue#1`](https://github.com/aplicacoesBoilerplate/PackagesVue/issues/1): estrutura base dos packages npm.

A distribuição multiplataforma do CLI e a publicação de `v0.0.1` serão habilitadas somente depois da implementação e verificação das issues do milestone.

## Preparação local de versão e release

Antes de contribuir, instale `golangci-lint` v2.13.2 e ative os hooks versionados neste checkout com `git config --local core.hooksPath .githooks`. O hook `pre-commit` confere `gofmt` e lint Go nas alterações, e `commit-msg` valida assuntos como `feat(cli): exibir versão`, `fix: ...` ou `release: ...`; commits incompatíveis também são rejeitados na CI. Os hooks locais não substituem os checks remotos de lint, `go test ./...`, `go vet ./...` e `go build ./...`.

O `go-gitsemver` calcula SemVer sobre o histórico e as tags Git (`fetch-depth: 0` em CI). Para reproduzir localmente o cálculo use a mesma revisão fixa do adapter compartilhado:

```bash
go install github.com/MyCarrier-DevOps/go-gitsemver@680c1c12d9a4f573a8da1b2e3ccebb3571b1cab6
go-gitsemver --show-variable SemVer
```

No primeiro release, `go-gitsemver.yml` usa `base-version: 0.0.0`, `next-version: 0.0.1` e `commit-message-incrementing: Disabled`. Esta exceção de bootstrap mantém `v0.0.1` apesar dos commits históricos `feat`; após a primeira tag, remover **ambas** as opções temporárias (`next-version` e `Disabled`) para voltar a incrementar por Conventional Commits. Se a versão calculada não corresponder à versão aprovada para publicação, interromper a release em vez de mover uma tag existente.

`bash scripts/validate-release-version.sh v0.0.1 0.0.1` é um auxiliar local anterior que verifica igualdade com a tag aprovada e rejeita tags locais duplicadas. Ele não é chamado pelo versionamento central: a policy compartilhada valida milestone/override e reconcilia tag e Release remotas, aceitando reexecução idêntica sem sobrescrita.

Antes de publicar, confira a distribuição localmente sem criar release:

```bash
goreleaser check
goreleaser release --snapshot --clean
```

O snapshot gera arquivos para Windows, Linux e macOS em `dist/`, inclusive `checksums.txt`. Um repositório derivado deve consumir o workflow reutilizável com sua própria identidade, histórico e permissões, sem acesso de escrita à organização que mantém este template. Snapshots não publicam uma GitHub Release.

## Caller de versionamento centralizado

O workflow `.github/workflows/ci.yml` reúne `go-ci`, `validate-pr` e `publish`.
Os dois últimos consomem o central publicado em
`f3e7a0168da9e3882840302add3ac8d162bac81e`, revisão integrada à master do
repositório `aplicacoesBoilerplate/.github`. Enviam somente `go-gitsemver`,
caminho do projeto e branch alvo; não dependem de Cobra, GoReleaser ou hooks.
A publicação aguarda a CI da aplicação e só é elegível após push em master
associado a um PR integrado permitido pelo central.

Veja [o guia do piloto](docs/versioning-pilot.md) para fluxo, testes seguros,
environment opcional, bootstrap e limites de publicação. Esta integração local
não conclui a distribuição de binários prevista na issue #8 nem publica a
primeira versão da aplicação.
