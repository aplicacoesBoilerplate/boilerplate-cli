package boilerplate

import "github.com/spf13/cobra"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstrap em projeto existente (detecta pom.xml/package.json)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("TODO: detectar pom.xml e/ou package.json, injetar <repositories> e registry, criar settings se faltar")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
