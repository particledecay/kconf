package kubeconfig_test

import (
	"testing"

	apiruntime "k8s.io/apimachinery/pkg/runtime"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
)

func TestAddContext(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"add a context if not already in kubeconfig": func(t *testing.T) {
			context := &clientcmdapi.Context{
				LocationOfOrigin: "/home/user/.kube/config",
				Cluster:          "test",
				AuthInfo:         "test-1",
				Namespace:        "default",
				Extensions:       map[string]apiruntime.Object{},
			}
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 0)
			k, _ := kc.GetConfig()
			testName := "testContext"
			result := k.AddContext(testName, context)

			if result != testName {
				t.Errorf("expected: %s, got: %s", testName, result)
			}

			AssertContext(t, k, testName)

			if k.Contexts[testName].Cluster != "test" {
				t.Errorf("expected: test, got: %s", k.Contexts[testName].Cluster)
			}
		},
		"do not add a context if already in kubeconfig": func(t *testing.T) {
			context := &clientcmdapi.Context{
				LocationOfOrigin: "/home/user/.kube/config",
				Cluster:          "test",
				AuthInfo:         "test",
				Namespace:        "default",
				Extensions:       map[string]apiruntime.Object{},
			}
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k, _ := kc.GetConfig()
			testName := "test"
			result := k.AddContext(testName, context)

			if result != "" {
				t.Errorf("expected: empty string, got: %s", result)
			}

			AssertContext(t, k, testName)

			if k.Contexts[testName].Cluster != "test" {
				t.Errorf("expected: test, got: %s", k.Contexts[testName].Cluster)
			}
		},
		"add a context if already in kubeconfig but with different fields": func(t *testing.T) {
			context := &clientcmdapi.Context{
				LocationOfOrigin: "/home/user/.kube/config",
				Cluster:          "test-1",
				AuthInfo:         "test",
				Namespace:        "default",
				Extensions:       map[string]apiruntime.Object{},
			}
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k, _ := kc.GetConfig()
			testName := "test"
			result := k.AddContext(testName, context)

			if result != "test-1" {
				t.Errorf("expected: test-1, got: %s", result)
			}

			AssertContext(t, k, "test-1")

			if k.Contexts["test-1"].Cluster != "test-1" {
				t.Errorf("expected: test-1, got: %s", k.Contexts["test-1"].Cluster)
			}
		},
		"add a context with context name": func(t *testing.T) {
			context := &clientcmdapi.Context{
				LocationOfOrigin: "/home/user/.kube/config",
				Cluster:          "test-1",
				AuthInfo:         "test",
				Namespace:        "default",
				Extensions:       map[string]apiruntime.Object{},
			}
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 0)
			k, _ := kc.GetConfig()
			testName := "test"
			result := k.AddContext(testName, context)

			if result != testName {
				t.Errorf("expected: %s, got: %s", testName, result)
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}

func TestAddCluster(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"add a cluster not already in kubeconfig": func(t *testing.T) {
			cluster := &clientcmdapi.Cluster{
				LocationOfOrigin:         "/home/user/.kube/config",
				Server:                   "https://example-test.com:6443",
				InsecureSkipTLSVerify:    false,
				CertificateAuthority:     "/etc/ssl/certs/dummy.crt",
				CertificateAuthorityData: DummyCert.Raw,
				Extensions:               map[string]apiruntime.Object{},
			}
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 0)
			k, _ := kc.GetConfig()
			testName := "testCluster"
			result := k.AddCluster(testName, cluster)

			if result != testName {
				t.Errorf("expected: %s, got: %s", testName, result)
			}

			AssertCluster(t, k, testName)

			if k.Clusters[testName].Server != "https://example-test.com:6443" {
				t.Errorf("expected: https://example-test.com:6443, got: %s", k.Clusters[testName].Server)
			}
		},
		"do not add a cluster if already in kubeconfig": func(t *testing.T) {
			cluster := &clientcmdapi.Cluster{
				LocationOfOrigin:         "/home/user/.kube/config",
				Server:                   "https://example-test.com:6443",
				InsecureSkipTLSVerify:    true,
				CertificateAuthority:     "/etc/ssl/certs/dummy.crt",
				CertificateAuthorityData: DummyCert.Raw,
				Extensions:               map[string]apiruntime.Object{},
			}
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k, _ := kc.GetConfig()
			testName := "test"
			result := k.AddCluster(testName, cluster)

			if result != "" {
				t.Errorf("expected: empty string, got: %s", result)
			}

			AssertCluster(t, k, testName)

			if k.Clusters[testName].Server != "https://example-test.com:6443" {
				t.Errorf("expected: https://example-test.com:6443, got: %s", k.Clusters[testName].Server)
			}
		},
		"add a cluster if already in kubeconfig but with different fields": func(t *testing.T) {
			cluster := &clientcmdapi.Cluster{
				LocationOfOrigin:         "/home/user/.kube/config",
				Server:                   "https://example-test.com:6443",
				InsecureSkipTLSVerify:    false,
				CertificateAuthority:     "/etc/ssl/certs/dummy.crt",
				CertificateAuthorityData: DummyCert.Raw,
				Extensions:               map[string]apiruntime.Object{},
			}
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k, _ := kc.GetConfig()
			testName := "test"
			result := k.AddCluster(testName, cluster)

			if result != "test-1" {
				t.Errorf("expected: test-1, got: %s", result)
			}

			AssertCluster(t, k, "test-1")

			if k.Clusters["test-1"].InsecureSkipTLSVerify != false {
				t.Errorf("expected: false, got: %v", k.Clusters["test-1"].InsecureSkipTLSVerify)
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}

func TestAddUser(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"add a user not already in kubeconfig": func(t *testing.T) {
			userToken := "bbbbbbbbbbbb-test"
			user := &clientcmdapi.AuthInfo{
				LocationOfOrigin: "/home/user/.kube/config",
				Token:            userToken,
			}
			_ = GenerateAndReplaceGlobalKubeconfig(t, 0, 0)
			k := GetGlobalKubeconfig(t)
			testName := "testUser"
			result := k.AddUser(testName, user)

			if result != testName {
				t.Errorf("expected: %s, got: %s", testName, result)
			}

			AssertUser(t, k, testName)

			if k.AuthInfos[testName].Token != userToken {
				t.Errorf("expected: %s, got: %s", userToken, k.AuthInfos[testName].Token)
			}
		},
		"do not add a user if already in kubeconfig": func(t *testing.T) {
			userToken := "bbbbbbbbbbbb-test"
			user := &clientcmdapi.AuthInfo{
				LocationOfOrigin: "/home/user/.kube/config",
				Token:            userToken,
			}
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k := GetGlobalKubeconfig(t)
			testName := "test"
			result := k.AddUser(testName, user)

			if result != "" {
				t.Errorf("expected: empty string, got: %s", result)
			}

			AssertUser(t, k, testName)

			if k.AuthInfos[testName].Token != userToken {
				t.Errorf("expected: %s, got: %s", userToken, k.AuthInfos[testName].Token)
			}
		},
		"add a user if already in kubeconfig but with different fields": func(t *testing.T) {
			userToken := "bbbbbbbbbbbb-test-1"
			user := &clientcmdapi.AuthInfo{
				LocationOfOrigin: "/home/user/.kube/config",
				Token:            userToken,
			}
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k := GetGlobalKubeconfig(t)
			testName := "test"
			result := k.AddUser(testName, user)

			if result != "test-1" {
				t.Errorf("expected: test-1, got: %s", result)
			}

			AssertUser(t, k, "test-1")

			if k.AuthInfos["test-1"].Token != userToken {
				t.Errorf("expected: %s, got: %s", userToken, k.AuthInfos["test-1"].Token)
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}

func TestMoveContext(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"move context to new name": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k := GetGlobalKubeconfig(t)
			err := k.MoveContext("test", "test-44")

			if err != nil {
				t.Errorf("expected: nil, got: %s", err)
			}

			AssertNotContext(t, k, "test")
			AssertContext(t, k, "test-44")
		},
		"fail if old context does not exist": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k := GetGlobalKubeconfig(t)
			err := k.MoveContext("test-3", "test-44")

			if err == nil {
				t.Errorf("expected: error, got: nil")
			}

			AssertNotContext(t, k, "test-3")
			AssertNotContext(t, k, "test-44")
		},
		"fail if the new context already exists": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 2, 2)
			k := GetGlobalKubeconfig(t)
			err := k.MoveContext("test", "test-1")

			if err == nil {
				t.Errorf("expected: error, got: nil")
			}

			AssertContext(t, k, "test")
			AssertContext(t, k, "test-1")
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}
