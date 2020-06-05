package cmd

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

var namespaceName string

var useCmd = &cobra.Command{
	Use:     "use",
	Short:   "Set the current context",
	Long:    `Set the current context in the main kubeconfig`,
	Aliases: []string{"u"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("You must provide the name of a kubeconfig context")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		contextName := args[0]

		config, err := kubeconfig.GetConfig()
		if err != nil {
			log.Fatal().Msgf("Could not read main config")
		}

		// change the current context
		err = config.SetCurrentContext(contextName)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		fmt.Printf("Using context '%s'\n", contextName)

		// change the namespace
		err = config.SetNamespace(contextName, namespaceName)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		if namespaceName != "" {
			fmt.Printf("Setting preferred namespace '%s'\n", namespaceName)
		}

		err = config.Save()
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, err := getContextsFromConfig(toComplete)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	},
}

// flags for this subcommand
func init() {
	useCmd.Flags().StringVarP(&namespaceName, "namespace", "n", "", "set a namespace to use")
}
