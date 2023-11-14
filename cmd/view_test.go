package cmd_test

import (
	"io"
	"os"
	"testing"

	"github.com/particledecay/kconf/cmd"
	. "github.com/particledecay/kconf/test"
)

func TestViewCmd(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"print context to stdout": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 3)

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				w.Close()
				os.Stdout = oldStdout

				// read from r
				out, _ := io.ReadAll(r)

				// check that the output contains the contexts we expect
				AssertSubstrings(t, out, []string{"test-2"})
			}()

			// view context
			viewCmd := cmd.ViewCmd()
			viewCmd.SilenceErrors = true
			viewCmd.SilenceUsage = true
			viewCmd.SetArgs([]string{"test-2"})

			err := viewCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}
		},
		"fail if no arguments provided": func(t *testing.T) {
			viewCmd := cmd.ViewCmd()
			viewCmd.SilenceErrors = true
			viewCmd.SilenceUsage = true

			err := viewCmd.Execute()
			if err == nil {
				t.Error("expected error to occur")
			}
		},
		"fail if context does not exist": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 3)

			viewCmd := cmd.ViewCmd()
			viewCmd.SilenceErrors = true
			viewCmd.SilenceUsage = true
			viewCmd.SetArgs([]string{"test-5"})

			err := viewCmd.Execute()
			if err == nil {
				t.Error("expected error to occur")
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}
