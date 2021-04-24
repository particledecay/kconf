package cmd

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
)

// NamespaceCmd allows you to switch default namespaces within a context
func NamespaceCmd() *cobra.Command {
	command := &cobra.Command{
		Use:     "namespace",
		Short:   "Set preferred namespace",
		Long:    `Set the preferred namespace within the current context`,
		Aliases: []string{"ns"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must provide the name of a namespace within the current context")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			namespace := args[0]

			config, err := kc.GetConfig()
			if err != nil {
				return errors.New("could not read main config")
			}

			// fail if we have no current context
			if config.CurrentContext == "" {
				return errors.New("you must first set a current context before setting a preferred namespace")
			}

			if namespace == "" {
				return errors.New("namespace cannot be blank")
			}

			kc.Out.Log().Msgf("setting preferred namespace '%s'", namespace)
			_ = config.SetNamespace(config.CurrentContext, namespace)

			err = config.Save()
			if err != nil {
				return err
			}

			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			list, err := getNamespacesFromConfig(toComplete)
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}

			return list, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return command
}
