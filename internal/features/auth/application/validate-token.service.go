package application

import "context"

/**
 * TValidateTokenService implementa a validação de sessão do GitHub CLI.
 */
type TValidateTokenService struct{}

/**
 * NewValidateTokenService cria o serviço compartilhado de validação de sessão.
 */
func NewValidateTokenService() *TValidateTokenService {
	return &TValidateTokenService{}
}

/**
 * ValidateToken verifica se a sessão do GitHub CLI está apta para comandos autenticados.
 */
func (pService *TValidateTokenService) ValidateToken(pContext context.Context) error {
	// Futuramente:
	// 1. Executar "gh auth token".
	// 2. Validar o código de saída sem registrar o token.
	_ = pContext

	return nil
}
