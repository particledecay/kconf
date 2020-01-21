package kubeconfig

import (
	"os"
	"sort"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Read loads a kubeconfig file and returns an api.Config (client-go) type
func Read(filepath string) (*clientcmdapi.Config, error) {
	config, err := clientcmd.LoadFromFile(filepath)
	if err != nil {
		log.Debug().Msgf("Error while reading %s: %v", filepath, err)
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
func (k *KConf) List() []string {
	contexts := []string{}
	for context := range k.Contexts {
		contexts = append(contexts, context)
	}

	sort.Strings(contexts)
	return contexts
}
