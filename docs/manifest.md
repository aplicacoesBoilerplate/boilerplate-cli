# Manifesto DX dos packages

O manifesto DX descreve os artefatos que um package publica para a CLI aplicar em uma aplicacao consumidora: presets de configuracao, snippets, metadados de IDE e documentacao.

A CLI deve tratar o package como fonte da verdade. O nome final da CLI e definido por cada organizacao que clonar este template; por isso, os exemplos usam `<cli>`.

## Locais oficiais

- Packages Vue/npm: `.boilerplate/package.manifest.json`
- Packages Java/Maven: `META-INF/boilerplate/package.manifest.json`

## Exemplo Vue

```json
{
  "schemaVersion": "1.0",
  "package": "@empresa/ui",
  "language": "vue",
  "aliases": ["ui", "vuetify"],
  "minimumCliVersion": "0.1.0",
  "docs": {
    "remote": "https://docs.empresa.local/ui",
    "local": "apps/docs"
  },
  "ide": {
    "vscode": {
      "snippets": [".boilerplate/vscode/ui.code-snippets"]
    },
    "jetbrains": {
      "webTypes": "web-types.json",
      "liveTemplates": [".boilerplate/jetbrains/ui.xml"]
    }
  },
  "presets": {
    "vuetify": {
      "description": "Configura plugin Vuetify e tema base",
      "templates": [".boilerplate/templates/vuetify/plugin.ts"]
    }
  }
}
```

## Exemplo Java

```json
{
  "schemaVersion": "1.0",
  "package": "boilerplate-starter",
  "language": "java",
  "groupId": "br.com.empresa",
  "artifactId": "boilerplate-starter",
  "minimumCliVersion": "0.1.0",
  "docs": {
    "remote": "https://docs.empresa.local/java/boilerplate-starter",
    "javadoc": "https://docs.empresa.local/javadoc/boilerplate-starter"
  },
  "ide": {
    "vscode": {
      "snippets": ["META-INF/boilerplate/vscode/java.code-snippets"]
    },
    "jetbrains": {
      "liveTemplates": ["META-INF/boilerplate/jetbrains/live-templates.xml"]
    }
  },
  "presets": {
    "yml": {
      "templates": ["META-INF/boilerplate/templates/application.yml"]
    },
    "properties": {
      "templates": ["META-INF/boilerplate/templates/application.properties"]
    }
  }
}
```

## Inspecao

```bash
<cli> manifest inspect .boilerplate/package.manifest.json
```

Saida esperada:

```text
Package: @empresa/ui
Language: vue
Minimum CLI: 0.1.0
Aliases: ui, vuetify
Presets: vuetify
IDE: vscode, jetbrains
Docs: remote, local
```

## Regras

- `schemaVersion` deve ser `1.0`.
- `language` deve ser `java` ou `vue`.
- `package` e `minimumCliVersion` sao obrigatorios.
- Caminhos de snippets, templates e metadados de IDE devem ser relativos ao package.
- Caminhos absolutos ou com escape para fora da raiz do package devem ser recusados.

## Fora de escopo

Os packages nao devem receber testes automatizados internos como parte deste contrato. Validacoes com Playwright, k6, Kafka ou ferramentas similares devem ficar em uma aplicacao externa de QA.
