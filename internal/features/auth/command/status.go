package command

import (
	"fmt"

	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/auth/contracts"
	"github.com/spf13/cobra"
)

/**
 * newStatusCommand cria o comando de diagnóstico da sessão autenticada.
 */
func newStatusCommand(pService contracts.IAuthServices) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Consulta a disponibilidade da sessão do GitHub CLI",
		Args:  cobra.NoArgs,
		RunE: func(pCommand *cobra.Command, pArguments []string) error {
			response, err := pService.Status(pCommand.Context())
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(pCommand.OutOrStdout(), response.Authenticated)
			return err
		},
	}
}
