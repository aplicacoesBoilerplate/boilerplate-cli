package boilerplate

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/install"
	"github.com/spf13/cobra"
)

var installPreset string
var installDryRun bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Instala packages Java ou Vue e aplica presets DX",
}

var installVueCmd = &cobra.Command{
	Use:   "vue <package>",
	Short: "Planeja instalacao de package Vue/npm",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return printInstallPlan(cmd, install.Request{
			Language: install.LanguageVue,
			Package:  args[0],
			Preset:   installPreset,
			DryRun:   installDryRun,
		})
	},
}

var installJavaCmd = &cobra.Command{
	Use:   "java <artifact>",
	Short: "Planeja instalacao de package Java/Maven",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return printInstallPlan(cmd, install.Request{
			Language: install.LanguageJava,
			Package:  args[0],
			Preset:   installPreset,
			DryRun:   installDryRun,
		})
	},
}

func init() {
	installCmd.PersistentFlags().StringVar(&installPreset, "preset", "default", "Preset de configuracao declarado pelo manifesto DX")
	installCmd.PersistentFlags().BoolVar(&installDryRun, "dry-run", false, "Mostra o plano sem alterar arquivos")
	installCmd.AddCommand(installVueCmd, installJavaCmd)
	rootCmd.AddCommand(installCmd)
}

func printInstallPlan(cmd *cobra.Command, request install.Request) error {
	plan, err := install.BuildPlan(request)
	if err != nil {
		return err
	}

	cmd.Printf("Language: %s\n", plan.Language)
	cmd.Printf("Package: %s\n", plan.Package)
	cmd.Printf("Preset: %s\n", plan.Preset)
	cmd.Printf("Workdir: %s\n", plan.WorkDir)
	cmd.Printf("Manifest: %s\n", plan.ManifestHint)
	if request.DryRun {
		cmd.Println("Mode: dry-run")
	} else {
		cmd.Println("Mode: plan-only")
	}
	cmd.Println("Actions:")
	for _, action := range plan.Actions {
		cmd.Printf("- %s\n", action)
	}

	return nil
}
