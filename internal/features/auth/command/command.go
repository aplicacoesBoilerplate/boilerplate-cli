package command

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/auth/application"
	"github.com/spf13/cobra"
)

/**
 * NewCommand cria o grupo de comandos  'auth' (recursos de autenticação).
 */
func NewCommand() *cobra.Command {
	lAuthService := application.NewAuthServices()

	// Adiciona apenas um comando 'auth'
	lAuthCmd := &cobra.Command{
		Use:   "auth",
		Short: "Gerencia autenticação para GitHub Packages usando o GitHub CLI.",
	}

	// Adiciona os comandos do grupo auth.
	lAuthCmd.AddCommand(
		newLoginCommand(lAuthService),
		newLogoutCommand(lAuthService),
		newStatusCommand(lAuthService),
	)

	return lAuthCmd
}
