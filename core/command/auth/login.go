package auth

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/application/auth/services"
	"github.com/spf13/cobra"
)

/**
 * newLoginCommand cria o comando de configuração de credenciais locais.
 */
func newLoginCommand(pService services.ILoginService) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Configura credenciais locais para GitHub Packages",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := services.TLoginRequest{
				DryRun: false,
			}

			return pService.Login(cmd.Context(), request)
		},
	}
}
