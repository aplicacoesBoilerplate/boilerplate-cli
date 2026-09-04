# Documentacao dos packages

O comando `doc` deve abrir a documentacao declarada pelo manifesto DX do package instalado.

Nesta primeira base, o comando apenas resolve onde a CLI deve buscar a informacao. A abertura de navegador ou servidor local entra em etapa posterior.

## Vue

```bash
<cli> doc vue @empresa/ui
<cli> doc vue @empresa/ui --local
```

Origem esperada:

- `docs.remote` para documentacao publicada.
- `docs.local` para catalogo local, como VitePress.

## Java

```bash
<cli> doc java boilerplate-starter
<cli> doc java boilerplate-starter --local
```

Origem esperada:

- `docs.remote` para documentacao publicada.
- `docs.javadoc` para Javadoc publicado.

## Regra

A CLI nao deve ter URLs hardcoded de packages. Toda resolucao deve vir do manifesto DX publicado pelo proprio package.
