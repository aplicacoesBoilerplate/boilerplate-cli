package application

import (
	"context"

	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/auth/contracts"
)

/**
 * TLoginService implementa a configuração de credenciais para GitHub Packages.
 */
type TLoginService struct {
	validateTokenService *TValidateTokenService
}

/**
 * Login configura as credenciais necessárias para GitHub Packages.
 */
func (pService *TLoginService) Login(
	pContext context.Context,
	pRequest contracts.TLoginRequest,
) error {
	if err := pService.validateTokenService.ValidateToken(pContext); err != nil {
		return err
	}

	// Futuramente:
	// 1. Atualizar Maven.
	// 2. Atualizar npm.
	// 3. Respeitar pRequest.DryRun.
	return nil
}
