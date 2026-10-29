package contracts

import "context"

/**
 * IBaseActionService define o contrato de recursos que executam uma única ação.
 * Qualquer struct com um método Execute compatível satisfaz este contrato.
 */
type IBaseActionService[T any] interface {
	Execute(ctx context.Context, req T) error
}
