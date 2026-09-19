package command

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/auth/contracts"
	"github.com/spf13/cobra"
)

/**
 * newLoginCommand cria o comando de configuração de credenciais locais.
 */
func newLoginCommand(pService contracts.IAuthServices) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Configura credenciais locais para GitHub Packages",
		Args:  cobra.NoArgs,
		RunE: func(pCommand *cobra.Command, pArguments []string) error {
			request := contracts.TLoginRequest{
				DryRun: false,
			}

			return pService.Login(pCommand.Context(), request)
		},
	}
}
