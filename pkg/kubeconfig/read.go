package kubeconfig

import (
	"fmt"
	"os"
	"sort"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Read returns a Config object representing an entire Kubernetes config
func Read(filepath string) (*KConf, error) {
	k := &KConf{}
	_, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug().Msgf("File not found: %s", filepath)
			// return empty config object if file does not exist
			k.Config = *clientcmdapi.NewConfig()
			return k, nil
		}
		return nil, err
	}
	kubeconfig, err := clientcmd.LoadFromFile(filepath)
	if err != nil {
		log.Error().Msgf("Error while reading %s: %v", filepath, err)
		return nil, err
	}
	k.Config = *kubeconfig
	return k, nil
}

// List reads the kubeconfig and returns all of the available contexts
func List() error {
	contexts := []string{}
	mainConfig, err := Read(MainConfigPath)
	if err != nil {
		log.Debug().Msg("Could not read main config")
		return err
	}

	for context := range mainConfig.Contexts {
		contexts = append(contexts, context)
	}

	sort.Strings(contexts)
	for _, context := range contexts {
		fmt.Printf("%s\n", context)
	}
	return nil
}
