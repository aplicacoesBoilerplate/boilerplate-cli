package registry

import (
	authCommand "github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/auth/command"
	// "github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/dependency/command"
	// "github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/diagnostics/command"
	// "github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/manifest/command"
	// "github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/packages/command"
	// "github.com/aplicacoesBoilerplate/boilerplate-cli/internal/features/project/command"
	"github.com/spf13/cobra"
)

/**
 * NewTopLevelCommands retorna os comandos registrados diretamente na raiz.
 */
func NewTopLevelCommands() []*cobra.Command {
	return []*cobra.Command{
		authCommand.NewCommand(),
		// project.NewCommand(),
		// packages.NewCommand(),
		// dependency.NewCommand(),
		// manifest.NewCommand(),
		// diagnostics.NewDoctorCommand(),
		// diagnostics.NewAuditCommand(),
	}
}
