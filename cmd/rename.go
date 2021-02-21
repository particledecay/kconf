package cmd

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

var renameCmd = &cobra.Command{
	Use:     "rename",
	Short:   "Rename a kubeconfig context",
	Long:    `Rename a kubeconfig context`,
	Aliases: []string{"rn"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("You must provide the name of an existing context and a new context name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		contextName := args[0]
		renamedName := args[1]
		config, err := kubeconfig.GetConfig()
		if err != nil {
			log.Fatal().Msgf("Error while reading main config: %v", err)
		}

		// check that old name exists, new name does not
		if _, ok := config.Contexts[contextName]; !ok {
			log.Fatal().Msgf("Could not find a context named '%s'", contextName)
		}
		if _, ok := config.Contexts[renamedName]; ok {
			log.Fatal().Msgf("There is already a context named '%s'", renamedName)
		}

		// rename it
		config.Contexts[renamedName] = config.Contexts[contextName]
		delete(config.Contexts, contextName)

		err = config.Save()
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
	},
}
