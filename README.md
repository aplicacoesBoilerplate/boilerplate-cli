# boilerplate-cli

CLI em Go para iniciar e manter aplicações Java e Vue integradas ao ecossistema `aplicacoesBoilerplate`.

O repositório e os templates de origem permanecem públicos. Projetos gerados para outras organizações devem ser privados por padrão. A leitura de artefatos no GitHub Packages continua exigindo autenticação, mesmo quando o repositório de origem é público.

## Estado atual

O projeto está no milestone [`v0.0.1`](https://github.com/aplicacoesBoilerplate/boilerplate-cli/milestone/1). A árvore de comandos, a validação e a descoberta de monorepos já possuem um contrato executável; os handlers de sistema ainda em desenvolvimento falham com código explícito em vez de simular sucesso.

- A [issue #2](https://github.com/aplicacoesBoilerplate/boilerplate-cli/issues/2) define `auth`, `init`, `new`, `add`, `update`, `doctor` e `audit`, incluindo `--dry-run`, idempotência, códigos de saída estáveis e suporte a monorepos.
- A [issue #3](https://github.com/aplicacoesBoilerplate/boilerplate-cli/issues/3) implementa autenticação por `gh auth token` e edição segura das configurações Maven/npm, sem registrar credenciais.
- A [issue #1](https://github.com/aplicacoesBoilerplate/boilerplate-cli/issues/1) acompanha a entrega completa da primeira versão.

O contrato detalhado, inclusive códigos de saída e garantias de `--dry-run`, está em [`docs/commands.md`](docs/commands.md).

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

O fluxo de autenticação reutiliza a sessão do GitHub CLI. Não há OAuth Device Flow nem cofre de credenciais próprio; o token não aparece em argumentos ou logs e fica restrito às configurações exigidas pelos clientes Maven/npm e aos respectivos backups protegidos.

```text
boilerplate auth login [--dry-run]
boilerplate auth status
boilerplate auth logout [--dry-run]
```

`auth login` lê a conta ativa por `gh auth token --hostname github.com`, configura o servidor Maven `github-boilerplate` em `~/.m2/settings.xml` e o escopo `@aplicacoesBoilerplate` em `~/.npmrc`. XML, comentários, registries e servidores não gerenciados são preservados. Antes de uma escrita, o conteúdo anterior é salvo ao lado do arquivo com o sufixo `.boilerplate-cli.bak`; repetir o comando sobre o mesmo estado não reescreve arquivos nem backups.

`auth logout` remove apenas o servidor Maven e o bloco npm marcados pela CLI; ele não encerra a sessão do GitHub CLI. Configurações ambíguas ou marcadores incompletos resultam em conflito sem truncar o arquivo.

## Desenvolvimento

Requisitos:

- Go 1.26;
- acesso aos repositórios públicos da organização;
- GitHub CLI autenticado para configurar e consumir GitHub Packages.

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
