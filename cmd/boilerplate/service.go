package boilerplate

import "context"

type unavailableService struct{}

func unavailable(operation string) error {
	return NewCLIError(ExitConfiguration, operation+" ainda nao esta configurado nesta versao", nil)
}

func (unavailableService) Auth(context.Context, AuthRequest) error {
	return unavailable("auth")
}

func (unavailableService) Init(context.Context, InitRequest) error {
	return unavailable("init")
}

func (unavailableService) New(context.Context, NewRequest) error {
	return unavailable("new")
}

func (unavailableService) Add(context.Context, AddRequest) error {
	return unavailable("add")
}

func (unavailableService) Update(context.Context, UpdateRequest) error {
	return unavailable("update")
}

func (unavailableService) Doctor(context.Context, DoctorRequest) error {
	return unavailable("doctor")
}

func (unavailableService) Audit(context.Context, AuditRequest) error {
	return unavailable("audit")
}
