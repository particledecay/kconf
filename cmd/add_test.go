package cmd_test

import (
	"fmt"
	"testing"

	"github.com/particledecay/kconf/cmd"
	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
	"github.com/rs/zerolog"
)

func TestAddCmd(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"add a new kubeconfig": func(t *testing.T) {
			newConfig := GenerateAndReplaceGlobalKubeconfig(t, 1, 0)

			// add the kubeconfig
			addCmd := cmd.AddCmd()
			addCmd.SilenceErrors = true
			addCmd.SetArgs([]string{newConfig})

			err := addCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			k2 := GetGlobalKubeconfig(t)
			AssertContext(t, k2, "test")
		},
		"rename context, cluster, and user": func(t *testing.T) {
			newConfig := GenerateAndReplaceGlobalKubeconfig(t, 1, 1)

			// add the kubeconfig
			addCmd := cmd.AddCmd()
			addCmd.SilenceErrors = true
			addCmd.SetArgs([]string{newConfig})

			err := addCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			k2 := GetGlobalKubeconfig(t)
			AssertContext(t, k2, "test")
			AssertContext(t, k2, "test-1")
			AssertCluster(t, k2, "test")
			AssertCluster(t, k2, "test-1")
			AssertUser(t, k2, "test")
			AssertUser(t, k2, "test-1")
		},
		"add new kubeconfig with custom context name": func(t *testing.T) {
			newConfig := GenerateAndReplaceGlobalKubeconfig(t, 1, 1)

			// add the kubeconfig
			newContext := "booyah"
			addCmd := cmd.AddCmd()
			addCmd.SilenceErrors = true
			addCmd.SetArgs([]string{
				newConfig,
				fmt.Sprintf("--context-name=%s", newContext),
			})

			err := addCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			k2 := GetGlobalKubeconfig(t)
			AssertContext(t, k2, newContext)
		},
		"rename if custom context already exists": func(t *testing.T) {
			newConfig := GenerateAndReplaceGlobalKubeconfig(t, 1, 1)

			// add the kubeconfig
			addCmd := cmd.AddCmd()
			addCmd.SilenceErrors = true
			addCmd.SetArgs([]string{
				newConfig,
				"--context-name=test",
			})

			err := addCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			k2 := GetGlobalKubeconfig(t)
			AssertContext(t, k2, "test-1")
		},
		"fail if kubeconfig is invalid": func(t *testing.T) {
			newConfig := GenerateAndReplaceGlobalKubeconfig(t, 1, 1)

			// create the base kubeconfig
			tmpfile, err := MakeTmpFile()
			if err != nil {
				t.Fatal(err)
			}
			_, _ = tmpfile.Write([]byte("ajsdlkfjasldfkjlasdkf"))
			kc.MainConfigPath = tmpfile.Name()

			// add the kubeconfig
			addCmd := cmd.AddCmd()
			addCmd.SilenceUsage = true
			addCmd.SilenceErrors = true
			addCmd.SetArgs([]string{newConfig})
			err = addCmd.Execute()

			if err == nil {
				t.Error("Expected error when kubeconfig is invalid")
			}
		},
	}

	for name, test := range tests {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		t.Run(name, test)
		PostTestCleanup()
	}
}
