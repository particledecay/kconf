package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

var (
	verbose     bool
	contextName string
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
	Short: "Add in a new kubeconfig file and option context name",
	Long:  `Add a new kubeconfig file to the existing merged config file and option context name`,
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
		contextName := args[0]
		config, err := kubeconfig.GetConfig()
		if err != nil {
			log.Fatal().Msgf("Error while reading main config: %v", err)
		}
		err = config.Remove(contextName)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		err = config.Save()
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved contexts",
	Long:  `Print a list of all contexts previously saved in kubeconfig`,
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

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a specific context's config",
	Long:  `Display all of the config resources associated with a specific context`,
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

func init() {
	// flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "display debug messages")
	addCmd.Flags().StringVarP(&contextName, "contextName", "n", "", "set context name")
}

// Execute combines all of the available command functions
func Execute() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(viewCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msgf("Error during execution: %v", err)
	}
}
