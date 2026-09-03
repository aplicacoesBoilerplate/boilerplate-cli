package boilerplate

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type Dependencies struct {
	Service  Service
	Stdout   io.Writer
	Stderr   io.Writer
	Discover func(string) (Workspace, error)
	Runner   ProcessRunner
	HomeDir  func() (string, error)
}

type commandContext struct {
	service  Service
	discover func(string) (Workspace, error)
	root     string
	dryRun   bool
}

func NewRootCommand(dependencies Dependencies) *cobra.Command {
	dependencies = withDefaults(dependencies)
	state := &commandContext{service: dependencies.Service, discover: dependencies.Discover}
	root := &cobra.Command{
		Use:           "boilerplate",
		Short:         "CLI para aplicacoes Java e Vue do ecossistema aplicacoesBoilerplate",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) > 0 {
				return usageError("comando desconhecido: %s", args[0])
			}
			return command.Help()
		},
	}
	root.SetOut(dependencies.Stdout)
	root.SetErr(dependencies.Stderr)
	root.PersistentFlags().StringVar(&state.root, "root", ".", "diretorio raiz do projeto ou monorepo")
	root.PersistentFlags().BoolVar(&state.dryRun, "dry-run", false, "planeja a operacao sem gravar ou executar alteracoes")
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return NewCLIError(ExitUsage, err.Error(), nil)
	})

	root.AddCommand(
		newAuthCommand(state),
		newInitCommand(state),
		newNewCommand(state),
		newAddCommand(state),
		newUpdateCommand(state),
		newDoctorCommand(state),
		newAuditCommand(state),
	)
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	return root
}

func RunCLI(ctx context.Context, args []string, dependencies Dependencies) int {
	dependencies = withDefaults(dependencies)
	command := NewRootCommand(dependencies)
	command.SetArgs(args)
	if err := command.ExecuteContext(ctx); err != nil {
		code := ExitCodeFor(err)
		if code == ExitInternal && isCobraUsageError(err) {
			code = ExitUsage
		}
		_, _ = fmt.Fprintln(dependencies.Stderr, "erro:", err.Error())
		return int(code)
	}
	return int(ExitSuccess)
}

// Execute runs the process-facing command. Tests should call RunCLI instead.
func Execute() {
	os.Exit(RunCLI(context.Background(), os.Args[1:], Dependencies{}))
}

func withDefaults(dependencies Dependencies) Dependencies {
	if dependencies.Stdout == nil {
		dependencies.Stdout = os.Stdout
	}
	if dependencies.Stderr == nil {
		dependencies.Stderr = os.Stderr
	}
	if dependencies.Discover == nil {
		dependencies.Discover = DiscoverWorkspace
	}
	if dependencies.Runner == nil {
		dependencies.Runner = execProcessRunner{}
	}
	if dependencies.HomeDir == nil {
		dependencies.HomeDir = os.UserHomeDir
	}
	if dependencies.Service == nil {
		dependencies.Service = newDefaultService(
			newAuthService(dependencies.Runner, dependencies.HomeDir, dependencies.Stdout),
		)
	}
	return dependencies
}

func (c *commandContext) absoluteRoot() (string, error) {
	root, err := filepath.Abs(c.root)
	if err != nil {
		return "", NewCLIError(ExitConfiguration, "nao foi possivel resolver o diretorio raiz", err)
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return "", NewCLIError(ExitConfiguration, "o diretorio raiz nao existe ou nao e um diretorio", err)
	}
	return filepath.Clean(root), nil
}

func (c *commandContext) workspace() (Workspace, error) {
	root, err := c.absoluteRoot()
	if err != nil {
		return Workspace{}, err
	}
	return c.discover(root)
}

func isCobraUsageError(err error) bool {
	message := err.Error()
	return strings.HasPrefix(message, "unknown command") ||
		strings.HasPrefix(message, "unknown flag") ||
		strings.Contains(message, "requires") ||
		strings.Contains(message, "accepts")
}

func noArgs(_ *cobra.Command, args []string) error {
	if len(args) != 0 {
		return usageError("nenhum argumento posicional e aceito")
	}
	return nil
}

func exactArgs(count int) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if len(args) != count {
			return usageError("esperados %d argumentos posicionais, recebidos %d", count, len(args))
		}
		return nil
	}
}

func maximumArgs(count int) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if len(args) > count {
			return usageError("esperados no maximo %d argumentos posicionais, recebidos %d", count, len(args))
		}
		return nil
	}
}
