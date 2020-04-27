package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/build"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print version",
	Long:    `Print version information`,
	Aliases: []string{"ver"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose == true {
			if err := build.PrintLongVersion(); err != nil {
				log.Error().Msgf("%v", err)
			}
		} else {
			build.PrintVersion()
		}
	},
}
