package runtime

import (
	"context"
	"fmt"
	"io"

	"github.com/aplicacoesBoilerplate/boilerplate-cli/core/root"
)

/**
 * Execute monta e executa a CLI, retornando o código de saída do processo.
 */
func Execute(
	pContext context.Context,
	pArguments []string,
	pStdout io.Writer,
	pStderr io.Writer,
) int {
	rootCmd := root.NewCommand()

	rootCmd.SetArgs(pArguments)
	rootCmd.SetOut(pStdout)
	rootCmd.SetErr(pStderr)

	if err := rootCmd.ExecuteContext(pContext); err != nil {
		_, _ = fmt.Fprintln(pStderr, "erro:", err.Error())
		return 1
	}

	return 0
}
