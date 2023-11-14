package kubeconfig_test

import (
	"fmt"
	"testing"

	runtimeapi "k8s.io/apimachinery/pkg/runtime"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
	"github.com/rs/zerolog"
)

func TestSetNamespace(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"set namespace in the given context": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k := GetGlobalKubeconfig(t)

			testCtx := "test"
			testNamespace := "superman"

			err := k.SetNamespace(testCtx, testNamespace)
			if err != nil {
				t.Errorf(fmt.Sprintf("expected: nil, got: %v", err))
			}
			if k.Config.Contexts[testCtx].Namespace != testNamespace {
				t.Errorf(fmt.Sprintf("expected: %s, got: %s", testNamespace, k.Config.Contexts[testCtx].Namespace))
			}
		},
		"fail if context does not exist": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 0)
			k := GetGlobalKubeconfig(t)

			testCtx := "test-1"
			testNamespace := "superman"

			// cannot set namespace on a non-existent context
			err := k.SetNamespace(testCtx, testNamespace)
			if err == nil {
				t.Errorf("expected: error, got: nil")
			}

			// also make sure the context doesn't actually exist
			if _, ok := k.Config.Contexts[testCtx]; ok {
				t.Errorf("expected: nil, got: %v", k.Config.Contexts[testCtx])
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}

func TestSetCurrentContext(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"set current context": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 5, 5)
			k := GetGlobalKubeconfig(t)

			testCtx := "test-2"

			err := k.SetCurrentContext(testCtx)
			if err != nil {
				t.Errorf(fmt.Sprintf("expected: nil, got: %v", err))
			}

			if k.CurrentContext != testCtx {
				t.Errorf(fmt.Sprintf("expected: %s, got: %s", testCtx, k.CurrentContext))
			}
		},
		"fail if context does not exist": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 0)
			k := GetGlobalKubeconfig(t)

			testCtx := "test-1"

			err := k.SetCurrentContext(testCtx)
			if err == nil {
				t.Errorf("expected: error, got: nil")
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}

func TestSave(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"save valid kubeconfig": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k := GetGlobalKubeconfig(t)

			err := k.Save()
			if err != nil {
				t.Errorf(fmt.Sprintf("expected: nil, got: %v", err))
			}

			// read saved kubeconfig
			k2 := GetGlobalKubeconfig(t)

			AssertContext(t, k2, "test")
		},
		"fail if kubeconfig is unwriteable": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k := GetGlobalKubeconfig(t)

			// this path should not be writeable
			kc.MainConfigPath = "/sbin/kconf.config"

			err := k.Save()
			if err == nil {
				t.Errorf("expected: error, got: nil")
			}
		},
	}

	for name, test := range tests {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		t.Run(name, test)
		PostTestCleanup()
	}
}

func TestMerge(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"merge kubeconfig into existing kubeconfig": func(t *testing.T) {
			contextName := "booyah"
			config := clientcmdapi.NewConfig()
			config.Clusters[contextName] = &clientcmdapi.Cluster{
				LocationOfOrigin:         "/home/user/.kube/config",
				Server:                   fmt.Sprintf("https://example-%s.com:6443", contextName),
				InsecureSkipTLSVerify:    true,
				CertificateAuthority:     "/etc/ssl/certs/dummy.crt",
				CertificateAuthorityData: DummyCert.Raw,
				Extensions:               map[string]runtimeapi.Object{},
			}
			config.AuthInfos[contextName] = &clientcmdapi.AuthInfo{
				LocationOfOrigin: "/home/user/.kube/config",
				Token:            fmt.Sprintf("bbbbbbbbbbbb-%s", contextName),
			}
			config.Contexts[contextName] = &clientcmdapi.Context{
				LocationOfOrigin: "/home/user/.kube/config",
				Cluster:          contextName,
				AuthInfo:         contextName,
				Namespace:        "default",
				Extensions:       map[string]runtimeapi.Object{},
			}

			k := MockConfig(2)

			k.Merge(config, "booyah")

			AssertContext(t, k, "test")
			AssertContext(t, k, "test-1")
			AssertContext(t, k, "booyah") // merged context
		},
		"rename context if it already exists": func(t *testing.T) {
			k := MockConfig(1)
			k2 := MockConfig(2)

			// we need k and k2 to be two unique configs, so delete the first config in k2
			err := k2.Remove("test")
			if err != nil {
				t.Fatal(err)
			}

			// force context names to be the same but with different configs
			k2.Contexts["test"] = k2.Contexts["test-1"]
			delete(k2.Contexts, "test-1")

			k.Merge(&k2.Config, "")

			AssertContext(t, k, "test")
			AssertContext(t, k, "test-1") // renamed context
		},
		"rename cluster if it already exists": func(t *testing.T) {
			k := MockConfig(1)
			k2 := MockConfig(2)

			// we need k and k2 to be two unique configs, so delete the first config in k2
			err := k2.Remove("test")
			if err != nil {
				t.Fatal(err)
			}

			// force cluster names to be the same but with different configs
			k2.Clusters["test"] = k2.Clusters["test-1"]
			delete(k2.Clusters, "test-1")

			k.Merge(&k2.Config, "")

			AssertCluster(t, k, "test")
			AssertCluster(t, k, "test-1") // renamed cluster
		},
		"rename user if it already exists": func(t *testing.T) {
			k := MockConfig(1)
			k2 := MockConfig(2)

			// we need k and k2 to be two unique configs, so delete the first config in k2
			err := k2.Remove("test")
			if err != nil {
				t.Fatal(err)
			}

			// force user names to be the same but with different configs
			k2.AuthInfos["test"] = k2.AuthInfos["test-1"]
			delete(k2.AuthInfos, "test-1")

			k.Merge(&k2.Config, "")

			AssertUser(t, k, "test")
			AssertUser(t, k, "test-1") // renamed user
		},
		"do not merge identical contexts": func(t *testing.T) {
			k := MockConfig(1)
			k2 := MockConfig(1)

			k.Merge(&k2.Config, "")

			AssertContext(t, k, "test")
			AssertNotContext(t, k, "test-1")

			if len(k.Contexts) != len(k2.Contexts) {
				t.Errorf("expected: %d, got: %d", len(k2.Contexts), len(k.Contexts))
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}
