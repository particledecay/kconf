package kubeconfig_test

import (
	"fmt"
	"os"
	"sort"
	"testing"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
	"github.com/rs/zerolog"
	runtimeapi "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestRead(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"fail if kubeconfig doesn't exist": func(t *testing.T) {
			config, err := kc.Read("/some/nonexistent/path")

			if config != nil {
				t.Errorf("expected: nil, got: %v", config)
			}

			if err == nil {
				t.Error("expected: error, got: nil")
			}
		},
		"read and return a valid kubeconfig": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			config, err := kc.Read(kc.MainConfigPath)

			if config == nil {
				t.Error("expected: config, got: nil")
			}

			if err != nil {
				t.Errorf("expected: nil, got: %v", err)
			}
		},
		"read a kubeconfig from stdin": func(t *testing.T) {
			// redirect stdin
			oldStdin := os.Stdin
			r, w, _ := os.Pipe()
			os.Stdin = r

			defer func() {
				os.Stdin = oldStdin
			}()

			// write some data to stdin
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			confdata, _ := os.ReadFile(kc.MainConfigPath)
			w.Write(confdata)
			w.Close()

			// run function to read from stdin
			config, err := kc.Read("")

			if config == nil {
				t.Error("expected: config, got: nil")
			}

			if err != nil {
				t.Errorf("expected: nil, got: %v", err)
			}
		},
	}

	for name, test := range tests {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		t.Run(name, test)
		PostTestCleanup()
	}
}

func TestList(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"return all context names": func(t *testing.T) {
			contexts := []string{}
			_ = GenerateAndReplaceGlobalKubeconfig(t, 3, 3)
			k, _ := kc.GetConfig()
			for context := range k.Contexts {
				contexts = append(contexts, context)
			}
			sort.Strings(contexts)
			if len(contexts) != 3 {
				t.Errorf("expected: 3, got: %d", len(contexts))
			}

			// compare expected list to returned list
			expected := []string{"test", "test-1", "test-2"}
			for i := range contexts {
				if contexts[i] != expected[i] {
					t.Errorf("expected: %v, got: %v", expected, contexts)
				}
			}
		},
	}

	for name, test := range tests {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		t.Run(name, test)
		PostTestCleanup()
	}
}

func TestExport(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"return a single usable config when given a context name": func(t *testing.T) {
			contextName := "test-3"
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
			config.CurrentContext = contextName

			_ = GenerateAndReplaceGlobalKubeconfig(t, 5, 5)
			k, _ := kc.GetConfig()
			result, err := k.Export(contextName)

			if err != nil {
				t.Errorf("expected: nil, got: %v", err)
			}

			if kc.IsEqualCluster(result.Clusters[contextName], config.Clusters[contextName]) == false {
				t.Errorf("expected: %v, got: %v", config.Clusters[contextName], result.Clusters[contextName])
			}
			if kc.IsEqualUser(result.AuthInfos[contextName], config.AuthInfos[contextName]) == false {
				t.Errorf("expected: %v, got: %v", config.AuthInfos[contextName], result.AuthInfos[contextName])
			}
			if kc.IsEqualContext(result.Contexts[contextName], config.Contexts[contextName]) == false {
				t.Errorf("expected: %v, got: %v", config.Contexts[contextName], result.Contexts[contextName])
			}
		},
		"fail if context does not exist": func(t *testing.T) {
			contextName := "test-7"
			_ = GenerateAndReplaceGlobalKubeconfig(t, 5, 5)
			k, _ := kc.GetConfig()
			result, err := k.Export(contextName)

			if result != nil {
				t.Errorf("expected: nil, got: %v", result)
			}

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

func TestGetContent(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"convert a config into a separate kubeconfig": func(t *testing.T) {
			contextName := "test-3"
			_ = GenerateAndReplaceGlobalKubeconfig(t, 5, 5)
			k, _ := kc.GetConfig()

			// extract the bytes content
			content, err := k.GetContent(contextName)

			if len(content) == 0 {
				t.Errorf("expected: non-empty, got: empty")
			}

			if err != nil {
				t.Errorf("expected: nil, got: %v", err)
			}

			// convert back to a config as a validity check
			config, err := clientcmd.Load(content)

			if config == nil {
				t.Errorf("expected: non-nil, got: nil")
			}

			if err != nil {
				t.Errorf("expected: nil, got: %v", err)
			}
		},
	}

	for name, test := range tests {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		t.Run(name, test)
		PostTestCleanup()
	}
}

func TestGetConfig(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"return a kubeconfig": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)

			k, err := kc.GetConfig()

			if k == nil {
				t.Errorf("expected: non-nil, got: nil")
			}

			if err != nil {
				t.Errorf("expected: nil, got: %v", err)
			}

			AssertContext(t, k, "test")
		},
		"create new kubeconfig if file does not exist": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 0)

			// remove it
			err := os.Remove(kc.MainConfigPath)
			if err != nil {
				t.Fatal(err)
			}

			k, err := kc.GetConfig()

			if err != nil {
				t.Errorf("expected: nil, got: %v", err)
			}

			if len(k.Contexts) != 0 {
				t.Errorf("expected: 0, got: %d", len(k.Contexts))
			}
		},
	}

	for name, test := range tests {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		t.Run(name, test)
		PostTestCleanup()
	}
}
