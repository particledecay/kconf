package cmd_test

import (
	"testing"

	"github.com/particledecay/kconf/cmd"
	. "github.com/particledecay/kconf/test"
)

func TestUseCmd(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"select current context": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 2)

			// select "test-1"
			currentContext := "test-1"
			useCmd := cmd.UseCmd()
			useCmd.SilenceErrors = true
			useCmd.SetArgs([]string{currentContext})

			err := useCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			// should read new kubeconfig with new values
			k := GetGlobalKubeconfig(t)

			AssertContext(t, k, currentContext)
			if k.CurrentContext != currentContext {
				t.Errorf("expected: %s, got: %s", currentContext, k.CurrentContext)
			}
		},
		"set preferred namespace": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 5)

			// select "test-2"
			currentContext := "test-2"
			namespace := "kube-system"
			useCmd := cmd.UseCmd()
			useCmd.SilenceErrors = true
			useCmd.SetArgs([]string{currentContext, "-n", namespace})

			err := useCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			// should read new kubeconfig with new values
			k := GetGlobalKubeconfig(t)

			AssertContext(t, k, currentContext)
			if k.CurrentContext != currentContext {
				t.Errorf("expected: %s, got: %s", currentContext, k.CurrentContext)
			}

			if k.Config.Contexts[currentContext].Namespace != namespace {
				t.Errorf("expected: %s, got: %s", namespace, k.Config.Contexts[currentContext].Namespace)
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}
