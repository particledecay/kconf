package cmd

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

// RemoveCmd deletes a context and all associated resources from the existing kubeconfig file
func RemoveCmd() *cobra.Command {
	command := &cobra.Command{
		Use:     "rm",
		Short:   "Remove a kubeconfig from main file",
		Long:    `Remove a named context and associated resources from main kubeconfig file`,
		Aliases: []string{"r"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must provide the name of a kubeconfig context")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]
			config, err := kubeconfig.GetConfig()
			if err != nil {
				return err
			}
			err = config.Remove(contextName)
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
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}

			return list, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return command
}
