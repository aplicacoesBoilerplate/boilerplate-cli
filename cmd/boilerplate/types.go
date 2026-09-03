package boilerplate

import (
	"context"
	"errors"
	"fmt"
)

// ExitCode is part of the public CLI contract and must remain stable.
type ExitCode int

const (
	ExitSuccess        ExitCode = 0
	ExitUsage          ExitCode = 2
	ExitConfiguration  ExitCode = 3
	ExitAuthentication ExitCode = 4
	ExitExternal       ExitCode = 5
	ExitConflict       ExitCode = 6
	ExitInternal       ExitCode = 10
)

// CLIError keeps an internal cause while exposing only a sanitized message.
type CLIError struct {
	code    ExitCode
	message string
	cause   error
}

func NewCLIError(code ExitCode, message string, cause error) error {
	if code == ExitSuccess {
		code = ExitInternal
	}
	return &CLIError{code: code, message: message, cause: cause}
}

func (e *CLIError) Error() string {
	return e.message
}

func (e *CLIError) Unwrap() error {
	return e.cause
}

func (e *CLIError) ExitCode() ExitCode {
	return e.code
}

func ExitCodeFor(err error) ExitCode {
	if err == nil {
		return ExitSuccess
	}
	var coded interface{ ExitCode() ExitCode }
	if errors.As(err, &coded) {
		return coded.ExitCode()
	}
	return ExitInternal
}

type Platform string

const (
	PlatformJava Platform = "java"
	PlatformVue  Platform = "vue"
	PlatformAll  Platform = "all"
)

type Visibility string

const (
	VisibilityPrivate  Visibility = "private"
	VisibilityInternal Visibility = "internal"
	VisibilityPublic   Visibility = "public"
)

type AuthAction string

const (
	AuthLogin  AuthAction = "login"
	AuthLogout AuthAction = "logout"
	AuthStatus AuthAction = "status"
)

type AuditFormat string

const (
	AuditFormatText AuditFormat = "text"
	AuditFormatJSON AuditFormat = "json"
)

type Workspace struct {
	Root         string
	JavaProjects []string
	VueProjects  []string
}

type AuthRequest struct {
	Root   string
	DryRun bool
	Action AuthAction
}

type InitRequest struct {
	Workspace Workspace
	DryRun    bool
}

type NewRequest struct {
	Root       string
	DryRun     bool
	Platform   Platform
	Name       string
	Owner      string
	Directory  string
	Visibility Visibility
}

type AddRequest struct {
	Workspace Workspace
	DryRun    bool
	Platform  Platform
	Package   string
	Version   string
}

type UpdateRequest struct {
	Workspace Workspace
	DryRun    bool
	Platform  Platform
}

type DoctorRequest struct {
	Workspace Workspace
}

type AuditRequest struct {
	Workspace Workspace
	DryRun    bool
	Format    AuditFormat
	Output    string
}

// Service isolates command parsing from filesystem, process and network effects.
type Service interface {
	Auth(context.Context, AuthRequest) error
	Init(context.Context, InitRequest) error
	New(context.Context, NewRequest) error
	Add(context.Context, AddRequest) error
	Update(context.Context, UpdateRequest) error
	Doctor(context.Context, DoctorRequest) error
	Audit(context.Context, AuditRequest) error
}

func usageError(format string, args ...any) error {
	return NewCLIError(ExitUsage, fmt.Sprintf(format, args...), nil)
}
