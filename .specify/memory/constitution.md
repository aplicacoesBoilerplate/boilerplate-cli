# boilerplate-cli Constitution

## Core Principles

### I. Template independente de organização

Este repositório MUST servir como base para que organizações criem suas próprias CLIs.
Nomes, owners, escopos de packages, repositórios, URLs e credenciais da
`aplicacoesBoilerplate` MUST ser configuráveis nos fluxos entregues; valores dessa
organização MAY existir como padrão da instância original, mas MUST NOT impedir o uso por
outras organizações. Repositórios derivados e forks MUST ser tratados como instâncias
independentes: a CLI MUST NOT presumir acesso ou permissão de escrita no repositório de
origem. Isso mantém o template reutilizável sem transferir identidade ou autoridade da
organização mantenedora para quem o adota.

### II. Origem de projetos e packages controlada

A CLI MUST contemplar a inicialização de projetos a partir de repositórios template
públicos ou de templates privados da organização proprietária da instância da CLI,
além da manutenção de projetos já existentes. A seleção da origem e do owner de destino
MUST ser explícita ou derivada de configuração verificável; a CLI MUST NOT assumir que
todo template é público ou que um fork pertence à organização de origem. Templates de
aplicação e de package da mesma organização MUST poder participar do fluxo, respeitando
seus contratos próprios. Esta delimitação evita inicializações a partir de uma origem
incorreta e preserva a autonomia das organizações derivadas.

### III. Dependências e packages sem duplicação

A CLI MUST abranger adição, remoção e atualização de dependências em projetos Java/Maven
e Vue/npm, incluindo os packages mantidos na organização da instância. Alterações em
`pom.xml` e `package.json` MUST reconhecer as dependências existentes, evitar entradas
duplicadas e preservar dados não relacionados. Operações repetidas MUST ser idempotentes,
e atualizações MUST respeitar as versões e a compatibilidade declaradas pelos packages.
O mesmo contrato MUST atender projetos individuais e monorepos quando ambos os tipos
de manifesto estiverem presentes. Isso permite manter aplicações e packages derivados
sem corromper os projetos consumidores.

### IV. Credenciais e privacidade por padrão

A CLI MUST reutilizar a autenticação existente do GitHub CLI para acessar templates
privados e GitHub Packages; MUST NOT implementar um OAuth Device Flow próprio nem
persistir cópia do token. Configurações locais de Maven/npm MUST ser alteradas sem
exibir credenciais em logs, erros ou saídas. Repositórios criados para organizações
destinatárias MUST ser privados por padrão, mesmo quando o template de origem for
público; publicação explícita exige escolha consciente do operador. A leitura de
GitHub Packages MUST exigir a autenticação apropriada independentemente da visibilidade
do repositório template. Essas regras protegem os projetos gerados e as credenciais.

### V. Contrato CLI previsível e verificável

Os comandos `auth`, `init`, `new`, `add`, `update`, `doctor` e `audit` MUST manter
responsabilidades distinguíveis e saídas com erros e códigos de saída estáveis.
Comandos que alteram arquivos ou repositórios MUST oferecer `--dry-run` e indicar
o que será alterado antes da execução efetiva. Mudanças em operações de templates,
dependências ou autenticação MUST ter validação automatizada para casos de sucesso,
falha e reexecução, sem acessar dados reais de consumidores nos testes. Contratos
explícitos e simulações permitem revisar alterações antes de aplicá-las em projetos.

## Limites técnicos e arquiteturais

A implementação da CLI MUST usar Go e Cobra e organizar recursos em Vertical Slices:
`internal/features/<recurso>/` concentra aplicação, adapters de comando, contratos e,
quando necessários, domínio e infraestrutura. `core` MUST conter somente runtime,
comando raiz e registro da árvore Cobra; contratos compartilhados MUST ficar em
`internal/shared`. Dependências de configuração adicionais, como Viper, MUST ter
necessidade concreta antes de sua inclusão. A CLI MUST distinguir o repositório
template de origem, a organização proprietária da instância e o repositório gerado
ou atualizado, inclusive quando houver forks.

## Fluxo de desenvolvimento e revisão

README e issues do repositório `aplicacoesBoilerplate/boilerplate-cli` MUST orientar o
escopo da instância original; decisões novas ou incompatibilidades entre fontes MUST
ser resolvidas em issue antes de alterar contratos. Cada mudança MUST ser revisada
quanto aos cinco princípios, incluindo uso por organização derivada, privacidade,
idempotência e impacto nos consumidores Java/Vue. Mudanças executáveis MUST passar por
`go test ./...`, `go vet ./...` e `go build ./...` antes da entrega. Comandos ainda
representados apenas por mensagens `TODO` MUST NOT ser apresentados como prontos para
modificar projetos reais. Funcionalidades previstas em issues não passam a ser
consideradas implementadas somente por constarem nesta constituição.

## Governance

Esta constituição prevalece sobre orientações de implementação conflitantes. Uma
emenda MUST registrar a motivação, os princípios afetados e eventuais impactos em
projetos derivados; a revisão MUST confirmar aderência ao README, às issues aplicáveis
e aos contratos existentes antes da aprovação. Alterações incompatíveis em princípios
ou sua remoção exigem incremento MAJOR; novos princípios/seções ou orientação
materialmente ampliada exigem MINOR; esclarecimentos sem mudança normativa exigem
PATCH. A primeira ratificação usa 1.0.0. Emendas MUST atualizar a versão e a data de
última alteração; a data de ratificação original permanece fixa. Cada revisão de
mudança de produto MUST verificar e documentar conformidade com os princípios e
justificar explicitamente qualquer exceção antes de sua aprovação.

**Version**: 1.0.0 | **Ratified**: 2026-09-23 | **Last Amended**: 2026-09-23
