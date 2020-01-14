package kubeconfig

import (
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// MainConfigPath is the file path to the main config
var MainConfigPath string

// Read returns a Config object representing an entire Kubernetes config
func Read(filepath string) (*clientcmdapi.Config, error) {
	_, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			// return empty config object if file does not exist
			return clientcmdapi.NewConfig(), nil
		}
		return nil, err
	}
	kubeconfig, err := clientcmd.LoadFromFile(filepath)
	if err != nil {
		log.Error().Msgf("Error while reading %s: %v", filepath, err)
		return nil, err
	}
	return kubeconfig, nil
}

func init() {
	MainConfigPath = path.Join(os.Getenv("HOME"), ".kube", "config")
}
