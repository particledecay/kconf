package cmd

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

var (
	contextName string
)

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add in a new kubeconfig file and optional context name",
	Long:    `Add a new kubeconfig file to the existing merged config file and optional context name`,
	Aliases: []string{"a"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("You must supply the path to a kubeconfig file")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		config, err := kubeconfig.GetConfig()
		if err != nil {
			log.Fatal().Msgf("Error while reading main config: %v", err)
		}

		newConfig, err := kubeconfig.Read(filepath)
		if err != nil {
			log.Fatal().Msgf("Error while reading %s: %v", filepath, err)
		}
		if config == nil {
			log.Fatal().Msgf("Could not find kubeconfig at %s", filepath)
		}

		err = config.Merge(newConfig, contextName)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		err = config.Save()
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
	},
}

// flags for this subcommand
func init() {
	addCmd.Flags().StringVarP(&contextName, "context-name", "n", "", "override context name")
}
