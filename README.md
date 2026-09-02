# boilerplate-cli

CLI em Go para DX dos packages privados `aplicacoesBoilerplate/PackagesJava` (Maven) e `PackagesVue` (npm).

## Conceito
Para **devs internos** (ja membros da org) instalarem packages em **projetos ja existentes ou novos**. Reusa credencial do `gh CLI` (`gh auth token` com `read:packages`).

## Comandos (planejados - issue #10)
- `boilerplate auth login` - le `gh auth token` -> `~/.m2/settings.xml` + `~/.npmrc`
- `boilerplate auth status`
- `boilerplate init` - bootstrap em projeto existente
- `boilerplate new [java|vue] <nome>`
- `boilerplate add [java|vue] <pkg>@<ver>`
- `boilerplate update` / `doctor` / `audit`

## Stack
Go + cobra + viper, distribuicao via GoReleaser + GitHub Releases (privado).

## Milestone
v0.0.1-beta - refs aplicacoesBoilerplate/PackagesJava#10, #14, #15
