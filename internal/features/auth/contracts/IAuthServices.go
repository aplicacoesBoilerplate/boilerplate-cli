package contracts

import "context"

/**
 * TLoginRequest representa os dados necessários para configurar a autenticação.
 */
type TLoginRequest struct {
	DryRun bool
}

/**
 * TLogoutRequest representa os dados necessários para encerrar a autenticação.
 */
type TLogoutRequest struct {
	DryRun bool
}

/**
 * TStatusResponse representa a disponibilidade da sessão do GitHub CLI.
 */
type TStatusResponse struct {
	Authenticated bool
}

/**
 * IAuthServices define as operações de autenticação expostas pela CLI.
 */
type IAuthServices interface {
	Login(pContext context.Context, pRequest TLoginRequest) error
	Logout(pContext context.Context, pRequest TLogoutRequest) error
	Status(pContext context.Context) (TStatusResponse, error)
}
