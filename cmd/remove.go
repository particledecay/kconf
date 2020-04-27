package cmd

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

var removeCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Remove a kubeconfig from main file",
	Long:    `Remove a named context and associated resources from main kubeconfig file`,
	Aliases: []string{"r"},
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
			log.Fatal().Msgf("Error while reading main config: %v", err)
		}
		err = config.Remove(contextName)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		err = config.Save()
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
	},
}
