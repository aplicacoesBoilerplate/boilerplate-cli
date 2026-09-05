package boilerplate

import "github.com/spf13/cobra"

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Gerencia autenticacao para GitHub Packages",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Le token de gh auth token e escreve ~/.m2/settings.xml e ~/.npmrc",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("TODO: gh auth token -> ~/.m2/settings.xml (github-boilerplate) e ~/.npmrc (@aplicacoesBoilerplate:registry)")
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Valida credenciais (gh api /user e teste de resolucao)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("TODO: validar via gh api /user e mvn dependency:get")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd, authStatusCmd)
	rootCmd.AddCommand(authCmd)
}
