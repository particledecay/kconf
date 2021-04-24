package cmd

import (
	"github.com/spf13/cobra"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
)

// ListCmd displays all of the stored contexts in the kubeconfig file
func ListCmd() *cobra.Command {
	command := &cobra.Command{
		Use:     "list",
		Short:   "List all saved contexts",
		Long:    `Print a list of all contexts previously saved in kubeconfig`,
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := kc.GetConfig()
			if err != nil {
				return err
			}
			contexts, currentContext := config.List()
			for _, ctx := range contexts {
				if currentContext == ctx {
					kc.Out.Log().Msgf("* %s", ctx)
				} else {
					kc.Out.Log().Msgf("  %s", ctx)
				}
			}

			return nil
		},
	}

	return command
}
