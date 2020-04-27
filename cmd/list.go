package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all saved contexts",
	Long:    `Print a list of all contexts previously saved in kubeconfig`,
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := kubeconfig.GetConfig()
		if err != nil {
			log.Fatal().Msgf("Could not read main config")
		}
		contexts, currentContext := config.List()
		for _, ctx := range contexts {
			if currentContext == ctx {
				fmt.Println("*", ctx)
			} else {
				fmt.Println(" ", ctx)
			}
		}
	},
}
