package test

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
)

// stores the filepaths for cleanup
var tmpfiles []string

// DummyCert is a fake certificate we can use for testing
var DummyCert = &x509.Certificate{
	SerialNumber: big.NewInt(2019),
	Subject: pkix.Name{
		Organization:  []string{"Company, INC."},
		Country:       []string{"US"},
		Province:      []string{"FL"},
		Locality:      []string{"Miami"},
		StreetAddress: []string{"100 SE 2nd St"},
		PostalCode:    []string{"33131"},
	},
	NotBefore:             time.Now(),
	NotAfter:              time.Now().AddDate(10, 0, 0),
	IsCA:                  true,
	ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	BasicConstraintsValid: true,
}

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
			CertificateAuthority:     "/etc/ssl/certs/dummy.crt",
			CertificateAuthorityData: DummyCert.Raw,
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

// MakeTmpFile creates a new, empty file to be used as a kubeconfig
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

func PostTestCleanup() {
	// reset the kubeconfig
	kc.MainConfigPath = path.Join(os.Getenv("HOME"), ".kube", "config")
	CleanupFiles()
}

func AssertContext(t *testing.T, k *kc.KConf, contextName string) {
	contexts, _ := k.List()
	for _, ctx := range contexts {
		if ctx == contextName {
			return
		}
	}
	t.Errorf("expected context '%s' in '%v'", contextName, contexts)
}

func AssertCluster(t *testing.T, k *kc.KConf, clusterName string) {
	var clusterNames = []string{}
	for cluster := range k.Clusters {
		if cluster == clusterName {
			return
		}
		clusterNames = append(clusterNames, cluster)
	}
	t.Errorf("expected cluster '%s' in '%v'", clusterName, clusterNames)
}

func AssertUser(t *testing.T, k *kc.KConf, userName string) {
	var userNames = []string{}
	for user := range k.AuthInfos {
		if user == userName {
			return
		}
		userNames = append(userNames, user)
	}
	t.Errorf("expected user '%s' in '%v'", userName, userNames)
}

func GenerateAndReplaceGlobalKubeconfig(t *testing.T, newContexts, baseContexts int) string {
	// save brand new kubeconfig and replace the global one
	k := MockConfig(newContexts)
	err := k.Save()
	if err != nil {
		t.Fatal(err)
	}
	newConfig := kc.MainConfigPath

	// replace the global kubeconfig
	if baseContexts == 0 {
		tmpfile, err := MakeTmpFile()
		if err != nil {
			t.Fatal(err)
		}
		kc.MainConfigPath = tmpfile.Name()
	} else {
		k2 := MockConfig(baseContexts)
		err := k2.Save()
		if err != nil {
			t.Fatal(err)
		}
	}

	return newConfig
}

func GetGlobalKubeconfig(t *testing.T) *kc.KConf {
	k, err := kc.GetConfig()
	if err != nil {
		t.Fatal(err)
	}

	return k
}
