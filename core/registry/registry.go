package registry

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/core/command/auth"
	// "github.com/aplicacoesBoilerplate/boilerplate-cli/core/command/dependency"
	// "github.com/aplicacoesBoilerplate/boilerplate-cli/core/command/diagnostics"
	// "github.com/aplicacoesBoilerplate/boilerplate-cli/core/command/manifest"
	// "github.com/aplicacoesBoilerplate/boilerplate-cli/core/command/packages"
	// "github.com/aplicacoesBoilerplate/boilerplate-cli/core/command/project"
	"github.com/spf13/cobra"
)

/**
 * NewTopLevelCommands retorna os comandos registrados diretamente na raiz.
 */
func NewTopLevelCommands() []*cobra.Command {
	return []*cobra.Command{
		auth.NewCommand(),
		// project.NewCommand(),
		// packages.NewCommand(),
		// dependency.NewCommand(),
		// manifest.NewCommand(),
		// diagnostics.NewDoctorCommand(),
		// diagnostics.NewAuditCommand(),
	}
}
