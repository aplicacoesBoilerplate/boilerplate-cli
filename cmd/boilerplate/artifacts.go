package boilerplate

import (
	"fmt"
	"os"

	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/dxartifacts"
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/manifest"
	"github.com/spf13/cobra"
)

var artifactsPreset string
var artifactsPackageRoot string
var artifactsProjectRoot string

var artifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "Planeja aplicacao de artefatos DX declarados por um manifesto",
}

var artifactsPlanCmd = &cobra.Command{
	Use:   "plan <manifesto>",
	Short: "Mostra snippets, templates e metadados de IDE que seriam aplicados",
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

		plan, err := dxartifacts.BuildPlan(dxartifacts.Request{
			Manifest:    packageManifest,
			Preset:      artifactsPreset,
			PackageRoot: artifactsPackageRoot,
			ProjectRoot: artifactsProjectRoot,
		})
		if err != nil {
			return err
		}

		cmd.Printf("Package: %s\n", plan.Package)
		cmd.Printf("Language: %s\n", plan.Language)
		cmd.Printf("Preset: %s\n", plan.Preset)
		cmd.Println("Artifacts:")
		if len(plan.Artifacts) == 0 {
			cmd.Println("- none")
			return nil
		}
		for _, artifact := range plan.Artifacts {
			cmd.Printf("- %s: %s -> %s\n", artifact.Kind, artifact.Source, artifact.Destination)
		}

		return nil
	},
}

func init() {
	artifactsPlanCmd.Flags().StringVar(&artifactsPreset, "preset", "default", "Preset declarado no manifesto DX")
	artifactsPlanCmd.Flags().StringVar(&artifactsPackageRoot, "package-root", ".", "Raiz do package que contem os artefatos")
	artifactsPlanCmd.Flags().StringVar(&artifactsProjectRoot, "project-root", ".", "Raiz do projeto consumidor")
	artifactsCmd.AddCommand(artifactsPlanCmd)
	rootCmd.AddCommand(artifactsCmd)
}
