package command

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/auth/contracts"
	"github.com/spf13/cobra"
)

/**
 * newLogoutCommand cria o comando para encerrar sessão autenticada.
 */
func newLogoutCommand(pService contracts.IAuthServices) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Encerra sessão autenticada no GitHub Packages",
		Args:  cobra.NoArgs,
		RunE: func(pCommand *cobra.Command, pArguments []string) error {
			lRequest := contracts.TLogoutRequest{
				DryRun: false,
			}

			return pService.Logout(pCommand.Context(), lRequest)
		},
	}
}
