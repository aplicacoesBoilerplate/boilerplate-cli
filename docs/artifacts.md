# Artefatos DX

Artefatos DX sao arquivos publicados por um package para melhorar a experiencia no projeto consumidor:

- snippets VS Code;
- live templates JetBrains;
- `web-types.json` para autocomplete Vue em JetBrains;
- templates de configuracao Java ou Vue;
- arquivos SCSS de identidade visual.

## Planejamento

```bash
<cli> artifacts plan fixtures/manifests/vue-ui.package.manifest.json --preset vuetify
```

Exemplo de saida:

```text
Package: @empresa/ui
Language: vue
Preset: vuetify
Artifacts:
- template: .boilerplate/templates/vuetify/plugin.ts -> src/plugins/plugin.ts
- vscode-snippet: .boilerplate/vscode/ui.code-snippets -> .vscode/boilerplate.code-snippets
- web-types: web-types.json -> .idea/webTypes/@empresa/ui.json
- jetbrains-live-template: .boilerplate/jetbrains/ui.xml -> .idea/templates/@empresa/ui.xml
```

## Regras

- Esta etapa apenas calcula origem e destino.
- A aplicacao real deve criar backup antes de alterar arquivos.
- Snippets pessoais e configuracoes de IDE existentes devem ser preservados.
- Packages nao devem receber testes internos por causa deste fluxo.
