package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/particledecay/kconf/pkg/kubeconfig"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Get the completion script for a shell",
	Long:  `Generate the completion script for a particular shell`,
}

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Get kconf completions for bash",
	Long:  `Generate the bash script for kconf completions`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

var completionFishCmd = &cobra.Command{
	Use:   "fish",
	Short: "Get kconf completions for fish shell",
	Long:  `Generate the fish script for kconf completions`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenFishCompletion(os.Stdout, true)
	},
}

var completionPowerShellCmd = &cobra.Command{
	Use:     "powershell",
	Short:   "Get kconf completions for fish shell",
	Long:    `Generate the fish script for kconf completions`,
	Aliases: []string{"ps"},
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenPowerShellCompletion(os.Stdout)
	},
}

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Get kconf completions for zsh",
	Long:  `Generate the zsh script for kconf completions`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenZshCompletion(os.Stdout)
	},
}

func getContextsFromConfig(partial string) (out []string, err error) {
	config, err := kubeconfig.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("could not read main config: %w", err)
	}

	for name := range config.Contexts {
		out = append(out, name)
	}

	return out, nil
}

func getNamespacesFromConfig(partial string) (out []string, err error) {
	config, err := kubeconfig.GetConfig()
	restConfig, err := kubeconfig.GetRestConfig(config)

	// fail if we have no current context
	if config.CurrentContext == "" {
		return []string{""}, errors.New("No current context detected. You must set one first with the `use` command.")
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
