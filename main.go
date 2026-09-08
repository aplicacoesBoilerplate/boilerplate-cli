package main

import (
	"context"
	"os"

	"github.com/aplicacoesBoilerplate/boilerplate-cli/core/runtime"
)

func main() {
	exitCode := runtime.Execute(
		context.Background(),
		os.Args[1:],
		os.Stdout,
		os.Stderr,
	)

	os.Exit(exitCode)
}
