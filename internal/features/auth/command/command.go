package command

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/auth/application"
	"github.com/spf13/cobra"
)

/**
 * NewCommand cria o grupo de comandos de autenticação.
 */
func NewCommand() *cobra.Command {
	authService := application.NewAuthServices()

	// Adiciona apenas um comando 'auth'
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Gerencia autenticação para GitHub Packages usando o GitHub CLI.",
	}

	// Adiciona os comandos do grupo auth.
	authCmd.AddCommand(
		newLoginCommand(authService),
		newLogoutCommand(authService),
		newStatusCommand(authService),
	)

	return authCmd
}
