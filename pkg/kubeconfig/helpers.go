package kubeconfig

import (
	"fmt"
	"io/ioutil"
	"os"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// stores the filepaths for cleanup
var tmpfiles []string

// MockConfig generates a mock KConf with `num` number of resources
func MockConfig(num int) *KConf {
	config := clientcmdapi.NewConfig()
	for i := 0; i < num; i++ {
		var name string
		if i == 0 {
			name = "test"
		} else {
			name = fmt.Sprintf("test-%d", i)
		}
		config.Clusters[name] = &clientcmdapi.Cluster{
			LocationOfOrigin:         "/home/user/.kube/config",
			Server:                   fmt.Sprintf("https://example-%s.com:6443", name),
			InsecureSkipTLSVerify:    true,
			CertificateAuthority:     "bbbbbbbbbbbb",
			CertificateAuthorityData: []byte("bbbbbbbbbbbb"),
		}
		config.AuthInfos[name] = &clientcmdapi.AuthInfo{
			LocationOfOrigin: "/home/user/.kube/config",
			Token:            fmt.Sprintf("bbbbbbbbbbbb-%s", name),
		}
		config.Contexts[name] = &clientcmdapi.Context{
			LocationOfOrigin: "/home/user/.kube/config",
			Cluster:          name,
			AuthInfo:         name,
			Namespace:        "default",
		}
	}
	return &KConf{Config: *config}
}

// WriteTmpConfig writes a kubeconfig to a temporary file instead of the default location
func WriteTmpConfig(config *KConf) error {
	tmpfile, err := ioutil.TempFile("", "config")
	if err != nil {
		return err
	}

	// for later removal
	tmpfiles = append(tmpfiles, tmpfile.Name())

	// override the default kubeconfig path
	MainConfigPath = tmpfile.Name()

	err = config.Save()
	if err != nil {
		return err
	}

	return nil
}

// CleanupFiles removes any temporary files created during tests
func CleanupFiles() {
	for _, filepath := range tmpfiles {
		if err := os.Remove(filepath); err != nil {
			fmt.Print(err)
		}
	}
}
