package cmd_test

import (
	"testing"

	"github.com/particledecay/kconf/cmd"
	. "github.com/particledecay/kconf/test"
)

func TestRenameCmd(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"rename an existing context": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 1)

			// rename "test" -> "some-other-thing"
			oldName := "test"
			newName := "some-other-thing"
			renameCmd := cmd.RenameCmd()
			renameCmd.SilenceErrors = true
			renameCmd.SetArgs([]string{oldName, newName})

			err := renameCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			// should read new kubeconfig with new values
			k := GetGlobalKubeconfig(t)

			AssertNotContext(t, k, oldName)
			AssertContext(t, k, newName)
		},
		"fail when context doesn't exist": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 1)

			// rename "test-jlaskdjf" -> "some-other-thing"
			oldName := "test-jlaskdjf"
			newName := "some-other-thing"
			renameCmd := cmd.RenameCmd()
			renameCmd.SilenceErrors = true
			renameCmd.SetArgs([]string{oldName, newName})

			err := renameCmd.Execute()
			if err == nil {
				t.Error("expected error to occur")
			}

			// should have not modified kubeconfig
			k := GetGlobalKubeconfig(t)

			AssertNotContext(t, k, newName)
		},
		"fail if insufficient arguments are provided": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 1)

			// rename "test-jlaskdjf" -> "some-other-thing"
			oldName := "test-jlaskdjf"
			newName := "some-other-thing"
			renameCmd := cmd.RenameCmd()
			renameCmd.SilenceErrors = true
			renameCmd.SilenceUsage = true
			renameCmd.SetArgs([]string{oldName})

			err := renameCmd.Execute()
			if err == nil {
				t.Error("expected error to occur")
			}

			// should have not modified kubeconfig
			k := GetGlobalKubeconfig(t)

			AssertNotContext(t, k, newName)
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}
