# Validacao manual da base DX

Este guia valida a base inicial da CLI sem depender de testes dentro dos packages.

## Build local

```bash
go test ./...
go vet ./...
go build ./...
```

## Manifestos validos

```bash
go run . manifest inspect fixtures/manifests/vue-ui.package.manifest.json
go run . manifest inspect fixtures/manifests/java-starter.package.manifest.json
```

Saida esperada para Vue:

```text
Package: @empresa/ui
Language: vue
Minimum CLI: 0.1.0
Aliases: ui, vuetify
Presets: prime, theme-scss, vuetify
IDE: vscode, jetbrains
Docs: remote, local
```

Saida esperada para Java:

```text
Package: boilerplate-starter
Language: java
Minimum CLI: 0.1.0
Aliases: -
Presets: properties, yml
IDE: vscode, jetbrains
Docs: remote, javadoc
```

## Manifesto invalido

```bash
go run . manifest inspect fixtures/manifests/invalid-escape.package.manifest.json
```

Resultado esperado: falha informando que o preset referencia caminho inseguro fora da raiz do package.

## Branding

```bash
go run . config show
```

Para validar override por ambiente:

```bash
BOILERPLATE_CLI_NAME=empresa BOILERPLATE_CLI_BINARY=empresa go run . config show
```

## Install plan-only

Em um projeto consumidor Vue que tenha `package.json`:

```bash
go run /caminho/para/boilerplate-cli install vue @empresa/ui --preset vuetify --dry-run
```

Em um projeto consumidor Java que tenha `pom.xml`:

```bash
go run /caminho/para/boilerplate-cli install java boilerplate-starter --preset yml --dry-run
```

## Doc plan-only

```bash
go run . doc vue @empresa/ui
go run . doc vue @empresa/ui --local
go run . doc java boilerplate-starter
go run . doc java boilerplate-starter --local
```

## Observacao

Os packages continuam sem testes internos. Validacoes com Playwright, k6, Kafka e cenarios de consumo real devem ficar em uma aplicacao externa de QA.
