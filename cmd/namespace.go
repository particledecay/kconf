package cmd

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

var namespaceCmd = &cobra.Command{
	Use:     "namespace",
	Short:   "Set preferred namespace",
	Long:    `Set the preferred namespace within the current context`,
	Aliases: []string{"ns"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("You must provide the name of a namespace within the current context")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		namespace := args[0]

		config, err := kubeconfig.GetConfig()
		if err != nil {
			log.Fatal().Msgf("Could not read main config")
		}

		// fail if we have no current context
		if config.CurrentContext == "" {
			log.Fatal().Msgf("No current context detected. You must set one first with the `use` command.")
		}

		err = config.SetNamespace(config.CurrentContext, namespace)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		if namespace != "" {
			fmt.Printf("Setting preferred namespace '%s'\n", namespace)
		}

		err = config.Save()
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
	},
}
