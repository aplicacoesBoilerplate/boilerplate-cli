package boilerplate

import (
	"fmt"
	"os"
	"strings"

	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/manifest"
	"github.com/spf13/cobra"
)

var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Inspeciona manifestos DX publicados pelos packages",
}

var manifestInspectCmd = &cobra.Command{
	Use:   "inspect <arquivo>",
	Short: "Valida e resume um package.manifest.json",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("nao foi possivel abrir o manifesto DX: %w", err)
		}
		defer file.Close()

		packageManifest, err := manifest.Decode(file)
		if err != nil {
			return err
		}

		cmd.Printf("Package: %s\n", packageManifest.Package)
		cmd.Printf("Language: %s\n", packageManifest.Language)
		cmd.Printf("Minimum CLI: %s\n", packageManifest.MinimumCLIVersion)
		cmd.Printf("Aliases: %s\n", joinOrDash(packageManifest.Aliases))
		cmd.Printf("Presets: %s\n", joinOrDash(packageManifest.PresetNames()))
		cmd.Printf("IDE: %s\n", joinOrDash(packageManifest.IDENames()))
		cmd.Printf("Docs: %s\n", joinOrDash(docNames(packageManifest.Docs)))

		return nil
	},
}

func init() {
	manifestCmd.AddCommand(manifestInspectCmd)
	rootCmd.AddCommand(manifestCmd)
}

func joinOrDash(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, ", ")
}

func docNames(docs manifest.Docs) []string {
	var names []string
	if docs.Remote != "" {
		names = append(names, "remote")
	}
	if docs.Local != "" {
		names = append(names, "local")
	}
	if docs.Javadoc != "" {
		names = append(names, "javadoc")
	}
	return names
}
