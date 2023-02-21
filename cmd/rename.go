package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

// RenameCmd moves an existing context to a new context name
func RenameCmd() *cobra.Command {
	command := &cobra.Command{
		Use:     "rename",
		Short:   "Rename a kubeconfig context",
		Long:    `Rename a kubeconfig context`,
		Aliases: []string{"rn"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("you must provide the name of an existing context and a new context name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]
			renamedName := args[1]

			config, err := kubeconfig.GetConfig()
			if err != nil {
				return err
			}
			cmd.SilenceUsage = true

			err = config.MoveContext(contextName, renamedName)
			if err != nil {
				return err
			}

			err = config.Save()
			if err != nil {
				return err
			}

			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			list, err := getContextsFromConfig(toComplete)
			failOnError("", err)

			return list, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return command
}
