package cmd

import (
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/build"
)

// VersionCmd prints version information for this build
func VersionCmd() *cobra.Command {
	command := &cobra.Command{
		Use:     "version",
		Short:   "Print version",
		Long:    `Print version information`,
		Aliases: []string{"ver"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				if err := build.PrintLongVersion(); err != nil {
					return err
				}
			} else {
				build.PrintVersion()
			}

			return nil
		},
	}

	return command
}
