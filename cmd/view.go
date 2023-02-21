package cmd

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/particledecay/kconf/pkg/kubeconfig"
)

// ViewCmd displays all the resources associated with a context
func ViewCmd() *cobra.Command {
	command := &cobra.Command{
		Use:     "view",
		Short:   "View a specific context's config",
		Long:    `Display all of the config resources associated with a specific context`,
		Aliases: []string{"v"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must provide the name of a kubeconfig context")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]
			config, err := kubeconfig.GetConfig()
			if err != nil {
				return err
			}
			cmd.SilenceUsage = true

			// convert config into bytes
			content, err := config.GetContent(contextName)
			if err != nil {
				return err
			}

			// print config content
			_, _ = os.Stdout.Write(content)

			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			list, err := getContextsFromConfig(toComplete)
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}

			return list, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return command
}
