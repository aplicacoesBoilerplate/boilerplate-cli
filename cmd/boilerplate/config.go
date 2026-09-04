package boilerplate

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/branding"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Mostra configuracao efetiva do template da CLI",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Mostra nome, escopos e organizacao configurados",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := branding.Load()

		cmd.Printf("Name: %s\n", config.Name)
		cmd.Printf("Binary: %s\n", config.Binary)
		cmd.Printf("NPM scope: %s\n", config.NpmScope)
		cmd.Printf("Maven repo: %s\n", config.MavenRepo)
		cmd.Printf("GitHub org: %s\n", config.GitHubOrg)
		cmd.Printf("Maven groupId: %s\n", config.MavenGroup)

		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}
