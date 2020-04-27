package cmd

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

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

		contexts, currentContext := config.List()
		if currentContext == contextName {
			fmt.Printf("Current context is already '%s'\n", currentContext)
			return
		}

		hasContext := false
		for _, ctx := range contexts {
			if contextName == ctx {
				hasContext = true
				break
			}
		}
		if hasContext {
			config.WriteCurrentContext(contextName)
			config.Save()
			if err == nil {
				fmt.Printf("Using context '%s'\n", contextName)
			} else {
				log.Fatal().Msgf("Could not save current context to main config")
			}
		} else {
			fmt.Printf("Context '%s' not found\n", contextName)
		}
	},
}
