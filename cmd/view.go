package cmd

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View a specific context's config",
	Long:    `Display all of the config resources associated with a specific context`,
	Aliases: []string{"v"},
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
			log.Fatal().Msg("Could not read main config")
		}

		// convert config into bytes
		content, err := config.GetContent(contextName)
		if err != nil {
			log.Fatal().Msgf("Error while converting context '%s': %v", contextName, err)
		}

		// print config content
		os.Stdout.Write(content)
	},
}
