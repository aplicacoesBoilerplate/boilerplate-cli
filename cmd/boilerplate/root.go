package boilerplate

import (
	"os"

	"github.com/spf13/cobra"
)

// Execute inicia a CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "boilerplate",
	Short: "CLI DX para packages Java (Maven) + Vue (npm) da org aplicacoesBoilerplate",
	Long: `boilerplate facilita a DX de devs internos:
- auth login/status (le gh auth token e escreve ~/.m2/settings.xml e ~/.npmrc)
- init (bootstrap em projeto existente - detecta pom.xml/package.json)
- new [java|vue] <nome> (scaffolding)
- add [java|vue] <pkg>@<ver>
- update / doctor / audit`,
}
