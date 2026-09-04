package boilerplate

import "github.com/aplicacoesBoilerplate/boilerplate-cli/internal/branding"

func init() {
	config := branding.Load()
	rootCmd.Use = config.Binary
	rootCmd.Short = "CLI DX para packages Java (Maven) + Vue (npm)"
	rootCmd.Long = config.Binary + ` facilita a DX de devs internos:
- auth login/status (le gh auth token e escreve ~/.m2/settings.xml e ~/.npmrc)
- init (bootstrap em projeto existente - detecta pom.xml/package.json)
- new [java|vue] <nome> (scaffolding)
- add [java|vue] <pkg>@<ver>
- update / doctor / audit`
}
