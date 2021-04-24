package cmd

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
)

// UseCmd selects a new current context and optional default namespace
func UseCmd() *cobra.Command {
	var namespaceName string

	command := &cobra.Command{
		Use:     "use",
		Short:   "Set the current context",
		Long:    `Set the current context in the main kubeconfig`,
		Aliases: []string{"u"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must provide the name of a kubeconfig context")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			config, err := kc.GetConfig()
			if err != nil {
				return err
			}

			// change the current context
			err = config.SetCurrentContext(contextName)
			if err != nil {
				return err
			}
			kc.Out.Log().Msgf("using context '%s'", contextName)

			if namespaceName != "" {
				// change the namespace
				err = config.SetNamespace(contextName, namespaceName)
				if err != nil {
					return err
				}
				kc.Out.Log().Msgf("setting preferred namespace '%s'", namespaceName)
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
	command.Flags().StringVarP(&namespaceName, "namespace", "n", "", "set a namespace to use")

	return command
}
