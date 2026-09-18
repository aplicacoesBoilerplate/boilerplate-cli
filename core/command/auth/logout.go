package auth

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/application/auth/services"
	"github.com/spf13/cobra"
)

/**
 * newLogoutCommand cria o comando para encerrar sessão autenticada.
 */
func newLogoutCommand(pService services.ILoginService) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Encerra sessão autenticada no GitHub Packages",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := services.TLoginRequest{
				DryRun: false,
			}

			return pService.Login(cmd.Context(), request)
		},
	}
}
