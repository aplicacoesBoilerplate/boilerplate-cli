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
 * NewStatusService cria o serviço responsável pela consulta da sessão autenticada.
 */
func NewStatusService(pValidateTokenService *TValidateTokenService) *TStatusService {
	return &TStatusService{
		validateTokenService: pValidateTokenService,
	}
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
