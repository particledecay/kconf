package cmd_test

import (
	"testing"

	"github.com/particledecay/kconf/cmd"
	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
)

func TestNamespaceCmd(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"fail if namespace not provided": func(t *testing.T) {
			namespaceCmd := cmd.NamespaceCmd()
			namespaceCmd.SilenceErrors = true
			namespaceCmd.SilenceUsage = true

			err := namespaceCmd.Execute()
			if err == nil {
				t.Error("expected error to occur")
			}
		},
		"fail if current context not set": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 1)

			namespaceCmd := cmd.NamespaceCmd()
			namespaceCmd.SilenceErrors = true
			namespaceCmd.SilenceUsage = true
			namespaceCmd.SetArgs([]string{"default"})

			err := namespaceCmd.Execute()
			if err == nil {
				t.Error("expected error to occur")
			}
		},
		"fail if namespace is blank": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 2)

			k, _ := kc.GetConfig()
			err := k.SetCurrentContext("test-1")
			if err != nil {
				t.Fatal(err)
			}
			err = k.Save()
			if err != nil {
				t.Fatal(err)
			}

			namespaceCmd := cmd.NamespaceCmd()
			namespaceCmd.SilenceErrors = true
			namespaceCmd.SilenceUsage = true
			namespaceCmd.SetArgs([]string{""})

			err = namespaceCmd.Execute()
			if err == nil {
				t.Error("expected error to occur")
			}
		},
		"set desired namespace": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 2)

			k, _ := kc.GetConfig()
			err := k.SetCurrentContext("test-1")
			if err != nil {
				t.Fatal(err)
			}
			err = k.Save()
			if err != nil {
				t.Fatal(err)
			}

			namespaceCmd := cmd.NamespaceCmd()
			namespaceCmd.SilenceErrors = true
			namespaceCmd.SilenceUsage = true
			namespaceCmd.SetArgs([]string{"kube-system"})

			err = namespaceCmd.Execute()
			if err != nil {
				t.Error("error should not have occurred")
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}
