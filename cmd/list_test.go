package cmd_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/particledecay/kconf/cmd"
	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestListCmd(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"print a list of contexts": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 5)

			logBuffer := new(bytes.Buffer)
			oldOut := kc.Out
			kc.Out = log.Output(zerolog.ConsoleWriter{Out: logBuffer, PartsExclude: []string{"time", "level"}})

			defer func() {
				kc.Out = oldOut
			}()

			// list contexts
			listCmd := cmd.ListCmd()
			listCmd.SilenceErrors = true
			listCmd.SilenceUsage = true

			err := listCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			// read from logBuffer
			out, _ := io.ReadAll(logBuffer)

			// check that the output contains the contexts we expect
			var expectedContexts = []string{"test", "test-1", "test-2", "test-3", "test-4"}
			AssertSubstrings(t, out, expectedContexts)
		},
		"mark current context if set": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 5)

			logBuffer := new(bytes.Buffer)
			oldOut := kc.Out
			kc.Out = log.Output(zerolog.ConsoleWriter{Out: logBuffer, PartsExclude: []string{"time", "level"}})

			defer func() {
				kc.Out = oldOut
			}()

			// set current context
			k, _ := kc.GetConfig()
			err := k.SetCurrentContext("test-1")
			if err != nil {
				t.Fatal(err)
			}
			err = k.Save()
			if err != nil {
				t.Fatal(err)
			}

			// list contexts
			listCmd := cmd.ListCmd()
			listCmd.SilenceErrors = true
			listCmd.SilenceUsage = true

			err = listCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			// read from logBuffer
			out, _ := io.ReadAll(logBuffer)

			// check that test-1 is marked with an asterisk
			AssertSubstrings(t, out, []string{"* test-1"})
		},
	}

	for name, test := range tests {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		t.Run(name, test)
		PostTestCleanup()
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}
