package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
)

// AddCmd merges a new kubeconfig into the existing kubeconfig file
func AddCmd() *cobra.Command {
	var contextName string

	command := &cobra.Command{
		Use:     "add",
		Short:   "Add in a new kubeconfig file and optional context name",
		Long:    `Add a new kubeconfig file to the existing merged config file and optional context name`,
		Aliases: []string{"a"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 && !kc.HasPipeData() {
				return errors.New("you must supply the path to a kubeconfig file")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := kc.GetConfig()
			if err != nil {
				return err
			}
			if config == nil {
				return fmt.Errorf("could not find kubeconfig at '%s'", kc.MainConfigPath)
			}

			filepath := ""
			if !kc.HasPipeData() {
				filepath = args[0]
			}
			newConfig, err := kc.Read(filepath)
			if err != nil {
				return err
			}

			config.Merge(newConfig, contextName)
			err = config.Save()
			if err != nil {
				return err
			}

			return nil

		},
	}
	command.Flags().StringVarP(&contextName, "context-name", "n", "", "override context name")

	return command
}
