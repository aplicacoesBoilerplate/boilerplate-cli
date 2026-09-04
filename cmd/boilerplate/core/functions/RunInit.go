package functions

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/aplicacoesBoilerplate/boilerplate-cli/cmd/boilerplate/core/structs"
	"github.com/aplicacoesBoilerplate/boilerplate-cli/cmd/boilerplate/internal"
)

func RunInit(cmd *cobra.Command, args []string, flags structs.SInitFlags) error {
	if flags.DryRun {
		cmd.Println("Modo Dry-Run ativado: nada será alterado no disco.")
	} else {
		cmd.Println("TODO: detectar pom.xml e/ou package.json, injetar <repositories> e registry, criar settings se faltar")
		if err := internal.AddMavenDependency("Teste", "Estrutura", "0.0.1"); err != nil {
			return fmt.Errorf("Falha: %w", err)
		}
	}
	return nil
}
