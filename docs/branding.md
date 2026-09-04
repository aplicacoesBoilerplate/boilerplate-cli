# Branding da CLI

Este repositório é um template. O nome `boilerplate` é apenas o padrão inicial e não deve ser tratado como contrato público por organizações que clonarem a CLI.

A configuração efetiva pode vir de `.boilerplate-cli.json` na raiz do projeto ou de variáveis de ambiente.

## Arquivo de configuração

Copie `.boilerplate-cli.example.json` para `.boilerplate-cli.json` e ajuste os valores:

```json
{
  "name": "empresa",
  "binary": "empresa",
  "npmScope": "@empresa",
  "mavenRepo": "github-empresa",
  "githubOrg": "empresa",
  "mavenGroupId": "br.com.empresa"
}
```

## Variáveis de ambiente

```bash
export BOILERPLATE_CLI_NAME=empresa
export BOILERPLATE_CLI_BINARY=empresa
export BOILERPLATE_NPM_SCOPE=@empresa
export BOILERPLATE_MAVEN_REPO=github-empresa
export BOILERPLATE_GITHUB_ORG=empresa
export BOILERPLATE_MAVEN_GROUP_ID=br.com.empresa
```

## Inspeção

```bash
<cli> config show
```

Saída esperada:

```text
Name: empresa
Binary: empresa
NPM scope: @empresa
Maven repo: github-empresa
GitHub org: empresa
Maven groupId: br.com.empresa
```

## Regra de produto

Documentações e exemplos devem usar `<cli>` quando estiverem descrevendo comportamento genérico. Use o nome real apenas em documentações específicas de uma organização.
