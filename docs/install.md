# Instalacao de packages

O comando `install` e o ponto de entrada para instalar packages Java ou Vue e aplicar os artefatos DX publicados no manifesto do proprio package.

Nesta primeira base, o comando monta e imprime o plano de instalacao. A execucao real de npm, Maven e aplicacao de templates deve ser adicionada nas proximas etapas.

## Vue

```bash
<cli> install vue @empresa/ui --preset vuetify --dry-run
<cli> install vue @empresa/ui --preset prime --dry-run
<cli> install vue @empresa/ui --preset theme-scss --dry-run
```

Fluxo esperado:

1. Detectar `package.json`.
2. Instalar package npm.
3. Ler manifesto em `node_modules/<package>/.boilerplate/package.manifest.json`.
4. Aplicar preset declarado pelo manifesto.
5. Sincronizar snippets VS Code e artefatos JetBrains.

## Java

```bash
<cli> install java boilerplate-starter --preset yml --dry-run
<cli> install java boilerplate-starter --preset properties --dry-run
```

Fluxo esperado:

1. Detectar `pom.xml`.
2. Adicionar dependency Maven.
3. Resolver dependencias.
4. Extrair manifesto do jar em `META-INF/boilerplate/package.manifest.json`.
5. Aplicar preset em `application.yml` ou `application.properties`.
6. Sincronizar snippets VS Code e live templates JetBrains.

## Regras de seguranca

- Instalacoes reais devem criar backup antes de alterar arquivos existentes.
- Templates devem ser aplicados com merge seguro.
- Arquivos pessoais de IDE nao devem ser sobrescritos.
- O comando deve aceitar `--dry-run` em todas as etapas destrutivas ou potencialmente mutantes.
- Packages nao recebem testes automatizados internos neste fluxo; validacoes pesadas ficam em uma aplicacao externa de QA.
