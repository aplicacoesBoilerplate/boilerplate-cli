# Contrato de comandos v0.0.1

Este documento define a interface estável da primeira versão do `boilerplate-cli`. Os comandos usam Cobra com `RunE`, validação antes de qualquer efeito e dependências injetadas. A raiz ativa `SilenceUsage` e `SilenceErrors`, portanto cada falha é impressa uma vez e o help aparece somente quando solicitado.

## Opções comuns

```text
--root <diretorio>  raiz do projeto ou monorepo; padrão: diretório atual
--dry-run           planeja sem gravar arquivos, criar repositórios ou instalar/atualizar dependências
```

O modo dry-run pode ler o filesystem e consultar metadados remotos necessários para montar o plano. Ele não grava credenciais e nunca as inclui em argumentos de processo, logs ou relatórios.

## Códigos de saída

| Código | Significado |
| --- | --- |
| `0` | operação concluída |
| `2` | comando, flag ou argumento inválido |
| `3` | configuração, ferramenta ou ambiente indisponível |
| `4` | autenticação ou autorização necessária |
| `5` | falha de rede ou serviço externo |
| `6` | conflito com estado local/remoto que não pode ser sobrescrito com segurança |
| `10` | falha interna inesperada |

Mensagens públicas são sanitizadas. A causa interna não é impressa automaticamente.

## Descoberta de projetos

Comandos de workspace reconhecem `pom.xml` e `package.json` na raiz e nos subdiretórios. Os resultados são caminhos absolutos, deduplicados e ordenados. A busca ignora `.git`, `.idea`, `node_modules`, `target`, `build` e `dist`, não atravessa symlinks de diretório e trata outro `.git` como fronteira de repositório.

Isso permite um projeto Java, Vue ou um monorepo misto sem depender da ordem retornada pelo filesystem.

## `auth`

```text
boilerplate auth login
boilerplate auth logout
boilerplate auth status
```

Os subcomandos não aceitam argumentos. `login` reutiliza a conta ativa do GitHub CLI por `gh auth token --hostname github.com` e configura:

- o servidor `github-boilerplate` em `~/.m2/settings.xml`;
- o registry `@aplicacoesBoilerplate` e sua credencial em `~/.npmrc`.

Os arquivos são atualizados atomicamente com permissão `0600` onde o sistema operacional oferece permissões POSIX. Antes de cada alteração, a versão anterior recebe o sufixo `.boilerplate-cli.bak`. Se uma etapa falhar, arquivos e backups já alterados são restaurados; configurações XML/INI desconhecidas permanecem intactas. Blocos gerenciados incompletos, duplicados ou fora de ordem retornam conflito (`6`) sem escrita.

`status` exige sessão válida do `gh`, capacidade de leitura de packages (`read:packages` ou o escopo superior `write:packages`) e entradas Maven/npm completas, mas nunca imprime a credencial. O escopo é lido do JSON de `gh auth status` sem usar `--show-token`; quando estiver ausente, pode ser adicionado com `gh auth refresh --hostname github.com --scopes read:packages`. `logout` remove somente as entradas gerenciadas pela CLI e mantém a sessão do `gh`. `login` e `logout` aceitam o `--dry-run` global e não criam diretórios nesse modo.

Elementos Maven autocontidos (`<settings/>` ou `<servers/>`) são recusados como conflito antes de qualquer escrita, pois inserir filhos sem expandi-los invalidaria o XML.

## `init`

```text
boilerplate init
```

Detecta todos os módulos compatíveis sob `--root` e aplica somente configurações ausentes. Repetir o comando sobre o estado desejado não reescreve arquivos nem executa instalações desnecessárias.

## `new`

```text
boilerplate new <java|vue> <nome> \
  [--owner <organizacao-ou-usuario>] \
  [--directory <destino>] \
  [--visibility <private|internal|public>]
```

`private` é o padrão. O destino local padrão é `<root>/<nome>`. Nome inválido ou destino que não possa ser usado sem sobrescrita retorna erro antes de clonar/gerar conteúdo.

A decisão pendente é se v0.0.1 cria o repositório remoto imediatamente ou se gera somente o diretório local até um publish explícito.

## `add`

```text
boilerplate add <java|vue> <pacote>@<versao>
```

A versão é obrigatória. Packages npm com escopo são aceitos, por exemplo `@aplicacoesBoilerplate/ui@1.2.3`. O catálogo de aliases/coordinates Java ainda precisa ser ratificado; a implementação não deve adivinhar `groupId`.

Em monorepo, o comando recebe todos os módulos compatíveis da plataforma escolhida. Repetir a mesma versão não altera arquivos.

## `update`

```text
boilerplate update [java|vue|all]
```

O padrão é `all`. Somente dependências pertencentes ao catálogo do ecossistema podem ser atualizadas; dependências de terceiros ficam fora do escopo. Se tudo já estiver na versão desejada, a operação termina sem escrita.

## `doctor`

```text
boilerplate doctor
```

Opera somente em leitura. Verifica descoberta dos módulos, ferramentas necessárias, sessão do GitHub CLI e configurações Maven/npm. Cada verificação tem resultado determinístico e não corrige arquivos automaticamente.

## `audit`

```text
boilerplate audit [--format <text|json>] [--output <arquivo>]
```

Opera somente em leitura e relata, por módulo/package, versão atual, versão desejada e estado. `text` no stdout é o padrão; `--output` grava somente o relatório final, de forma atômica, quando a implementação de governança estiver disponível.

O schema da API de governança de `PackagesJava#5` ainda precisa ser estabilizado antes do handler remoto.

## Fronteira de implementação

O parser e o scanner não executam processos nem acessam credenciais. Cada `RunE` cria um request tipado e delega a uma interface de serviço. Isso permite testar o contrato sem rede e injetar, na issue #3, filesystem e processos seguros para Windows, Linux e macOS.

Enquanto um handler real não estiver configurado, a CLI retorna `3`. Nenhum placeholder imprime `TODO` e retorna sucesso.
