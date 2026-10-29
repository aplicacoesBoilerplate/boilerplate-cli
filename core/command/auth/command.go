package auth

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/application/auth/services"
	"github.com/spf13/cobra"
)


/**
 * NewCommand cria o grupo de comandos de autenticação.
 */
func NewCommand() *cobra.Command {
	// 'Instância' do service.
	loginService := services.NewLoginService()

	// Adiciona apenas um comando 'auth'
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Gerencia autenticação para GitHub Packages usando o GitHub CLI.",
	}

	// Adiciona os comandos do grupo auth.
	authCmd.AddCommand(
		newLoginCommand(loginService),
		// newLogoutCommand(),
		// newStatusCommand(),
	)

	return authCmd
}
