package services

import "context"

/**
 * TLogoutRequest representa os dados necessários para executar logout.
 */
type TLogoutRequest struct {
	DryRun bool
}

/**
 * ILogoutService define o contrato de encerramento da sessão local.
 */
type ILogoutService interface {
	Logout(pContext context.Context, pRequest TLogoutRequest) error
}

/**
 * TLogoutService implementa o encerramento da autenticação local para GitHub Packages.
 */
type TLogoutService struct{}

/**
 * NewLogoutService cria a implementação padrão do serviço de Logout.
 */
func NewLogoutService() ILogoutService {
	return &TLogoutService{}
}

/**
 * Logout configura as entradas necessárias para encerrar a sessão no GitHub Packages.
 */
func (pService *TLogoutService) Logout(
	pContext context.Context,
	pRequest TLogoutRequest,
) error {
	// Futuramente:
	// 1. Ler a sessão do gh.
	// 2. Validar acesso.
	// 3. Atualizar Maven.
	// 4. Atualizar npm.
	// 5. Respeitar pRequest.DryRun.

	return nil
}

/**
 * Garante em compilação que TLogoutService atende ILogoutService, é o mais próximo que go tem do satifies.
 */
var _ ILogoutService = (*TLogoutService)(nil)
