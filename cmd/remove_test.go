package cmd_test

import (
	"testing"

	"github.com/particledecay/kconf/cmd"
	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
	"github.com/rs/zerolog"
)

func TestRemoveCmd(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"remove a context": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 1)

			// remove "test"
			ctxName := "test"
			removeCmd := cmd.RemoveCmd()
			removeCmd.SilenceErrors = true
			removeCmd.SetArgs([]string{ctxName})
			err := removeCmd.Execute()

			if err != nil {
				t.Error(err)
			}

			// should not have the "test" context
			k, err := kc.GetConfig()
			if err != nil {
				t.Error(err)
			}

			contexts, _ := k.List()
			for _, ctx := range contexts {
				if ctx == ctxName {
					t.Errorf("expected context '%s' not to be in '%v'", ctxName, contexts)
				}
			}
		},
		"fail when context doesn't exist": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 1)

			// remove "test-1" which doesn't exist
			ctxName := "test-1"
			removeCmd := cmd.RemoveCmd()
			removeCmd.SilenceErrors = true
			removeCmd.SetArgs([]string{ctxName})
			err := removeCmd.Execute()

			if err == nil {
				t.Error("expected error to occur")
			}

			// should have not modified kubeconfig
			k, err := kc.GetConfig()
			if err != nil {
				t.Error(err)
			}

			AssertContext(t, k, "test")
		},
	}

	for name, test := range tests {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		t.Run(name, test)
		PostTestCleanup()
	}
}
