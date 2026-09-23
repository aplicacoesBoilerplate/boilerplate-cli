package services

import "context"

/**
 * TLoginRequest representa os dados necessários para executar a autenticação.
 */
type TLoginRequest struct {
	DryRun bool
}

/**
 * ILoginService define o contrato da autenticação local.
 */
type ILoginService interface {
	Login(pContext context.Context, pRequest TLoginRequest) error
}

/**
 * TLoginService implementa a autenticação local para GitHub Packages.
 */
type TLoginService struct{}

/**
 * NewLoginService cria a implementação padrão do serviço de login.
 */
func NewLoginService() ILoginService {
	return &TLoginService{}
}

/**
 * Login configura as credenciais necessárias para GitHub Packages.
 */
func (pService *TLoginService) Login(
	pContext context.Context,
	pRequest TLoginRequest,
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
 * Garante em compilação que TLoginService atende ILoginService, é o mais próximo que go tem do satifies.
 */
var _ ILoginService = (*TLoginService)(nil)
