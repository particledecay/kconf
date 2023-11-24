package cmd

import (
	"errors"
	"os"

	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var verbose bool

func rootCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "kconf",
		Short: "kconf manages your kubeconfigs",
		Long: `kconf allows you to add and delete kubeconfigs by merging kubeconfig
				files together and breaking them apart appropriately.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

			// debug mode
			if verbose {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
			log.Debug().Msg("debug messaging turned on")
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("an action is required")
			}
			return nil
		},
	}
	command.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "display debug messages")

	return command
}

// Execute combines all of the available command functions
func Execute() {
	root := rootCmd()
	root.AddCommand(AddCmd())
	root.AddCommand(RemoveCmd())
	root.AddCommand(ListCmd())
	root.AddCommand(ViewCmd())
	root.AddCommand(UseCmd())
	root.AddCommand(VersionCmd())
	root.AddCommand(NamespaceCmd())
	root.AddCommand(RenameCmd())
	root.AddCommand(VerifyCmd())

	completion := CompletionCmd()
	completion.AddCommand(completionBashCmd(root))
	completion.AddCommand(completionFishCmd(root))
	completion.AddCommand(completionPowerShellCmd(root))
	completion.AddCommand(completionZshCmd(root))
	root.AddCommand(completion)

	err := root.Execute()
	failOnError("", err)
}
