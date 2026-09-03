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
