package application

import (
	"context"

	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/auth/contracts"
)

/**
 * TStatusService implementa a consulta de disponibilidade da sessão do GitHub CLI.
 */
type TStatusService struct {
	validateTokenService *TValidateTokenService
}

/**
 * Status consulta a disponibilidade da sessão do GitHub CLI.
 */
func (pService *TStatusService) Status(
	pContext context.Context,
) (contracts.TStatusResponse, error) {
	if err := pService.validateTokenService.ValidateToken(pContext); err != nil {
		return contracts.TStatusResponse{}, err
	}

	return contracts.TStatusResponse{Authenticated: true}, nil
}

var _ contracts.IAuthServices = (*TAuthServices)(nil)
