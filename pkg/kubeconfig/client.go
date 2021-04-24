package kubeconfig

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetRestConfig returns a rest.Config from a KConf type
func GetRestConfig(config *KConf) (*rest.Config, error) {
	content, err := config.GetContent(config.CurrentContext)
	if err != nil {
		return nil, err
	}

	// first get the DirectClientConfig
	clientConfig, err := clientcmd.NewClientConfigFromBytes(content)
	if err != nil {
		return nil, err
	}

	// get the rest.Config
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return restConfig, nil
}
