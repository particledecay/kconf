package kubeconfig

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetRestConfig returns a rest.Config from a KConf type
func GetRestConfig(config *KConf) (*rest.Config, error) {
	apiConfig := config.Config

	// first get the DirectClientConfig
	clientConfig := clientcmd.NewDefaultClientConfig(apiConfig, &clientcmd.ConfigOverrides{})

	// get the rest.Config
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return restConfig, nil
}
