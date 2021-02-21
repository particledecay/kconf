package cmd

import (
	"errors"

	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "kconf",
	Short: "kconf manages your kubeconfigs",
	Long: `kconf allows you to add and delete kubeconfigs by merging kubeconfig
			files together and breaking them apart appropriately.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("An action is required")
		}
		return nil
	},
}

func init() {
	// flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "display debug messages")
}

// Execute combines all of the available command functions
func Execute() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(namespaceCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(renameCmd)

	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionFishCmd)
	completionCmd.AddCommand(completionPowerShellCmd)
	completionCmd.AddCommand(completionZshCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msgf("Error during execution: %v", err)
	}
}
