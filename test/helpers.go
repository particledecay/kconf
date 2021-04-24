package test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
)

// stores the filepaths for cleanup
var tmpfiles []string

// MockConfig generates a mock KConf with `num` number of resources
func MockConfig(num int) *kc.KConf {
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

	tmpfile, _ := MakeTmpFile()
	kc.MainConfigPath = tmpfile.Name()

	return &kc.KConf{Config: *config}
}

// MakeTmpConfig creates a new, empty file to be used as a kubeconfig
func MakeTmpFile() (*os.File, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	tmpfile, err := ioutil.TempFile("", fmt.Sprintf("config-%d", r1.Intn(1000)))
	if err != nil {
		return nil, err
	}

	// for later removal
	tmpfiles = append(tmpfiles, tmpfile.Name())

	return tmpfile, nil
}

// CleanupFiles removes any temporary files created during tests
func CleanupFiles() {
	for _, filepath := range tmpfiles {
		if err := os.Remove(filepath); err != nil {
			kc.Out.Log().Err(err).Msgf("could not remove '%s'", filepath)
		}
	}
}
