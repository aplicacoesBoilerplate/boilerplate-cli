package application

import "github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/auth/contracts"

/**
 * TAuthServices agrupa as operações de autenticação da aplicação.
 */
type TAuthServices struct {
	*TLoginService
	*TLogoutService
	*TStatusService
}

/**
 * NewAuthServices cria a implementação padrão das operações de autenticação.
 */
func NewAuthServices() contracts.IAuthServices {
	lValidateTokenService := NewValidateTokenService()

	return &TAuthServices{
		TLoginService:  NewLoginService(lValidateTokenService),
		TLogoutService: NewLogoutService(lValidateTokenService),
		TStatusService: NewStatusService(lValidateTokenService),
	}
}

var _ contracts.IAuthServices = (*TAuthServices)(nil)
