package runtime

import (
	"bytes"
	"context"
	"testing"
)

/**
 * TestExecuteAuthCommands garante que as operações de Auth estão registradas na árvore Cobra.
 */
func TestExecuteAuthCommands(pTesting *testing.T) {
	pTesting.Parallel()

	testCases := []struct {
		name           string
		arguments      []string
		expectedStdout string
	}{
		{
			name:      "login",
			arguments: []string{"auth", "login"},
		},
		{
			name:      "logout",
			arguments: []string{"auth", "logout"},
		},
		{
			name:           "status",
			arguments:      []string{"auth", "status"},
			expectedStdout: "true\n",
		},
	}

	for _, testCase := range testCases {
		pTesting.Run(testCase.name, func(pTesting *testing.T) {
			pTesting.Parallel()

			var stdout bytes.Buffer
			var stderr bytes.Buffer

			exitCode := Execute(
				context.Background(),
				testCase.arguments,
				&stdout,
				&stderr,
			)

			if exitCode != 0 {
				pTesting.Fatalf("código de saída = %d, stderr = %q", exitCode, stderr.String())
			}

			if actualStdout := stdout.String(); actualStdout != testCase.expectedStdout {
				pTesting.Errorf("stdout = %q, esperado %q", actualStdout, testCase.expectedStdout)
			}
		})
	}
}
