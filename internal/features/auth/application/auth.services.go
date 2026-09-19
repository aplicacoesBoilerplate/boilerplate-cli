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
	validateTokenService := &TValidateTokenService{}

	return &TAuthServices{
		TLoginService: &TLoginService{
			validateTokenService: validateTokenService,
		},
		TLogoutService: &TLogoutService{
			validateTokenService: validateTokenService,
		},
		TStatusService: &TStatusService{
			validateTokenService: validateTokenService,
		},
	}
}

var _ contracts.IAuthServices = (*TAuthServices)(nil)
