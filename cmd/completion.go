package cmd

import (
	"errors"
	"fmt"
	"os"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CompletionCmd is the base autocompletion command
func CompletionCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "completion",
		Short: "Get the completion script for a shell",
		Long:  `Generate the completion script for a particular shell`,
	}

	return command
}

func completionBashCmd(command *cobra.Command) *cobra.Command {
	completion := &cobra.Command{
		Use:   "bash",
		Short: "Get kconf completions for bash",
		Long:  `Generate the bash script for kconf completions`,
		Run: func(cmd *cobra.Command, args []string) {
			command.GenBashCompletion(os.Stdout)
		},
	}

	return completion
}

func completionFishCmd(command *cobra.Command) *cobra.Command {
	completion := &cobra.Command{
		Use:   "fish",
		Short: "Get kconf completions for fish shell",
		Long:  `Generate the fish script for kconf completions`,
		Run: func(cmd *cobra.Command, args []string) {
			command.GenFishCompletion(os.Stdout, true)
		},
	}

	return completion
}

func completionPowerShellCmd(command *cobra.Command) *cobra.Command {
	completion := &cobra.Command{
		Use:     "powershell",
		Short:   "Get kconf completions for fish shell",
		Long:    `Generate the fish script for kconf completions`,
		Aliases: []string{"ps"},
		Run: func(cmd *cobra.Command, args []string) {
			command.GenPowerShellCompletion(os.Stdout)
		},
	}

	return completion
}

func completionZshCmd(command *cobra.Command) *cobra.Command {
	completion := &cobra.Command{
		Use:   "zsh",
		Short: "Get kconf completions for zsh",
		Long:  `Generate the zsh script for kconf completions`,
		Run: func(cmd *cobra.Command, args []string) {
			command.GenZshCompletion(os.Stdout)
		},
	}

	return completion
}

func getContextsFromConfig(partial string) (out []string, err error) {
	config, err := kc.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("could not read main config: %w", err)
	}

	for name := range config.Contexts {
		out = append(out, name)
	}

	return out, nil
}

func getNamespacesFromConfig(partial string) (out []string, err error) {
	config, err := kc.GetConfig()

	// fail if we have no current context
	if config.CurrentContext == "" {
		return []string{""}, errors.New("No current context detected. You must set one first with the `use` command.")
	}

	restConfig, err := kc.GetRestConfig(config)
	if err != nil {
		return []string{""}, err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return []string{""}, err
	}

	list, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return []string{""}, err
	}

	var namespaces []string
	for _, namespace := range list.Items {
		namespaces = append(namespaces, namespace.Name)
	}

	return namespaces, nil
}
