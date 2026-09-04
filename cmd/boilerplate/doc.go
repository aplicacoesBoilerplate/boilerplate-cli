package boilerplate

import (
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/dxdoc"
	"github.com/aplicacoesBoilerplate/boilerplate-cli/internal/install"
	"github.com/spf13/cobra"
)

var docLocal bool

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Resolve documentacao declarada pelo manifesto DX",
}

var docVueCmd = &cobra.Command{
	Use:   "vue <package>",
	Short: "Resolve documentacao de package Vue/npm",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return printDocResolution(cmd, dxdoc.Request{
			Language: install.LanguageVue,
			Package:  args[0],
			Local:    docLocal,
		})
	},
}

var docJavaCmd = &cobra.Command{
	Use:   "java <artifact>",
	Short: "Resolve documentacao de package Java/Maven",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return printDocResolution(cmd, dxdoc.Request{
			Language: install.LanguageJava,
			Package:  args[0],
			Local:    docLocal,
		})
	},
}

func init() {
	docCmd.PersistentFlags().BoolVar(&docLocal, "local", false, "Prefere documentacao local quando declarada no manifesto")
	docCmd.AddCommand(docVueCmd, docJavaCmd)
	rootCmd.AddCommand(docCmd)
}

func printDocResolution(cmd *cobra.Command, request dxdoc.Request) error {
	resolution, err := dxdoc.Resolve(request)
	if err != nil {
		return err
	}

	cmd.Printf("Language: %s\n", resolution.Language)
	cmd.Printf("Package: %s\n", resolution.Package)
	cmd.Printf("Manifest: %s\n", resolution.ManifestHint)
	cmd.Printf("Docs source: %s\n", resolution.Source)
	cmd.Println("Status: resolucao por manifesto preparada; abertura de navegador entra na proxima etapa")

	return nil
}
