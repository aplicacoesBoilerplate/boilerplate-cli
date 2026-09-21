package application

import (
	"context"

	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/auth/contracts"
)

/**
 * TLogoutService implementa a remoção de credenciais para GitHub Packages.
 */
type TLogoutService struct {
	validateTokenService *TValidateTokenService
}

/**
 * NewLogoutService cria o serviço responsável pela remoção de credenciais.
 */
func NewLogoutService(pValidateTokenService *TValidateTokenService) *TLogoutService {
	return &TLogoutService{
		validateTokenService: pValidateTokenService,
	}
}

/**
 * Logout remove as credenciais locais configuradas para GitHub Packages.
 */
func (pService *TLogoutService) Logout(
	pContext context.Context,
	pRequest contracts.TLogoutRequest,
) error {
	if err := pService.validateTokenService.ValidateToken(pContext); err != nil {
		return err
	}

	// Futuramente:
	// 1. Remover as configurações gerenciadas no Maven.
	// 2. Remover as configurações gerenciadas no npm.
	// 3. Respeitar pRequest.DryRun.
	return nil
}
