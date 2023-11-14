package kubeconfig_test

import (
	"testing"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
)

func TestRemove(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"remove a context": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)

			k, err := kc.GetConfig()
			if err != nil {
				t.Fatal(err)
			}

			err = k.Remove("test")
			if err != nil {
				t.Fatal(err)
			}

			AssertNotContext(t, k, "test")
		},
		"do not remove user if another context is using it": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 2, 2)

			k, err := kc.GetConfig()
			if err != nil {
				t.Fatal(err)
			}

			// force the second context to use the first user
			k.Contexts["test-1"].AuthInfo = "test"
			err = k.Remove("test")
			if err != nil {
				t.Error(err)
			}

			AssertNotContext(t, k, "test")

			_, ok := k.AuthInfos["test"]
			if !ok {
				t.Error("expected: user to exist, got: nil")
			}
		},
		"fail if context does not exist": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)

			k, err := kc.GetConfig()
			if err != nil {
				t.Fatal(err)
			}

			err = k.Remove("test-1")
			if err == nil {
				t.Error("expected: error, got: nil")
			}

			AssertContext(t, k, "test")
		},
	}

	for name, test := range tests {
		t.Run(name, test)
		PostTestCleanup()
	}
}
