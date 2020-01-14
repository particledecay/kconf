package cmd

import (
	"errors"

	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

var (
	verbose bool
)

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

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add in a new kubeconfig file",
	Long:  `Add a new kubeconfig file to the existing merged config file`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("You must supply the path to a kubeconfig file")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		config, err := kubeconfig.Read(args[0])
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		if config == nil {
			log.Fatal().Msgf("Could not find file at %s", args[0])
		}
		kubeconfig.Merge(config)
	},
}

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a kubeconfig from main file",
	Long:  `Remove a named context and associated resources from main kubeconfig file`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("You must provide the name of a kubeconfig context")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig.Remove(args[0])
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
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msgf("Error during execution: %v", err)
	}
}
