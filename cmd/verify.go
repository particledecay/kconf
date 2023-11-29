package cmd

import (
	"errors"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
)

// VerifyCmd checks the provided kubeconfig file for errors
func VerifyCmd() *cobra.Command {
	command := &cobra.Command{
		Use:     "verify",
		Short:   "Check for errors in the provided kubeconfig file",
		Long:    `Check the provided kubeconfig file for errors (contexts, clusters, users) and return any issues found`,
		Aliases: []string{"vfy"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 && !kc.HasPipeData() {
				return errors.New("you must supply the path to a kubeconfig file")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			filepath := ""
			if !kc.HasPipeData() {
				filepath = args[0]
			}

			config, err := kc.Read(filepath)
			if err != nil {
				return err
			}

			err = kc.ValidateConfig(config)
			if err != nil {
				kc.Out.Log().Msg("the following error(s) were found in the provided kubeconfig file:")
				rgx := regexp.MustCompile(`(?m)\[(.*?)\]`)
				for i, e := range strings.Split(rgx.FindStringSubmatch(err.Error())[1], ", ") {
					kc.Out.Log().Msgf("%d. %s", i+1, string(e))
				}
			} else {
				kc.Out.Log().Msg("no errors were found in the provided kubeconfig file")
			}

			return nil

		},
	}

	return command
}
