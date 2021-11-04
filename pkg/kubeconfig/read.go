package kubeconfig

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Read loads a kubeconfig file and returns an api.Config (client-go) type
func Read(filepath string) (*clientcmdapi.Config, error) {
	var config *clientcmdapi.Config
	var err error

	if filepath == "" && HasPipeData() {
		pipedData, _ := io.ReadAll(os.Stdin)
		config, err = clientcmd.Load(pipedData)
	} else {
		config, err = clientcmd.LoadFromFile(filepath)
	}

	if err != nil {
		log.Debug().
			Err(err).
			Msgf("error while reading '%s'", filepath)
		return nil, err
	}
	return config, nil
}

// GetConfig reads the main kubeconfig and returns a KConf type
func GetConfig() (*KConf, error) {
	k := &KConf{}
	_, err := os.Stat(MainConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug().Msgf("Main config does not yet exist")
			// return empty config object if file does not exist
			k.Config = *clientcmdapi.NewConfig()
			return k, nil
		}
		return nil, err
	}
	kubeconfig, err := Read(MainConfigPath)
	if err != nil {
		log.Debug().Msgf("Error while reading main config: %v", err)
		return nil, err
	}
	k.Config = *kubeconfig
	return k, nil
}

// List returns an array of contexts
func (k *KConf) List() ([]string, string) {
	currentContext := k.Config.CurrentContext
	contexts := []string{}
	for context := range k.Contexts {
		contexts = append(contexts, context)
	}
	sort.Strings(contexts)
	return contexts, currentContext
}

// Export returns a single context's config from a kubeconfig file
func (k *KConf) Export(name string) (*clientcmdapi.Config, error) {
	context, ok := k.Contexts[name]
	if !ok { // the context never existed
		return nil, fmt.Errorf("could not find context '%s'", name)
	}

	// create new config for only this context
	config := *clientcmdapi.NewConfig()
	config.Contexts[name] = context
	config.AuthInfos[context.AuthInfo] = k.AuthInfos[context.AuthInfo]
	config.Clusters[context.Cluster] = k.Clusters[context.Cluster]
	config.CurrentContext = name

	return &config, nil
}

// GetContent converts a single config into writeable content
func (k *KConf) GetContent(name string) ([]byte, error) {
	config, err := k.Export(name)
	if err != nil {
		log.Debug().Msgf("could not export context '%s'", name)
		return []byte{}, err
	}

	content, err := clientcmd.Write(*config)
	if err != nil {
		log.Debug().Msgf("error in clientcmd.Write command for context '%s'", name)
		return []byte{}, err
	}

	return content, nil
}
