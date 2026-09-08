package root

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/core/registry"
	"github.com/spf13/cobra"
)

/**
 * NewCommand monta a árvore principal da CLI.
 */
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "boilerplate",
		Short: "CLI para aplicações Java e Vue",
		SilenceUsage: true,
		SilenceErrors: true,
	}

	// Consome o registry para estruturar a árvore da CLI.
	rootCmd.AddCommand(registry.NewTopLevelCommands()...)

	return rootCmd
}
