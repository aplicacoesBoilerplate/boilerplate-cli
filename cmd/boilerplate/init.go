package boilerplate

import (
	"github.com/spf13/cobra"
	"github.com/aplicacoesBoilerplate/boilerplate-cli/cmd/boilerplate/core/functions"
	"github.com/aplicacoesBoilerplate/boilerplate-cli/cmd/boilerplate/core/structs"
)

var flags structs.SInitFlags

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstrap em projeto existente (detecta pom.xml/package.json)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return functions.RunInit(cmd, args, flags)
	},
}

func init() {
	initCmd.Flags().BoolVarP(&flags.DryRun, "dry-run", "d", false, "Faz uma varredura para garantir sucesso antes da execução do comando")
	rootCmd.AddCommand(initCmd)
}
