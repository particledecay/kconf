package kubeconfig_test

import (
	"testing"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
	"github.com/rs/zerolog"
)

func TestRestConfig(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"fail if current context not set": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k, _ := kc.GetConfig()
			k.CurrentContext = ""

			config, err := kc.GetRestConfig(k)
			if err == nil {
				t.Error("expected error to occur")
			}

			if config != nil {
				t.Errorf("expected: nil, got: %v", config)
			}
		},
		"fail with an unusable kubeconfig": func(t *testing.T) {
			_ = GenerateAndReplaceGlobalKubeconfig(t, 1, 1)
			k, _ := kc.GetConfig()
			k.CurrentContext = "test"

			config, err := kc.GetRestConfig(k)
			if err == nil {
				t.Error("expected error to occur")
			}

			if config != nil {
				t.Errorf("expected: nil, got: %v", config)
			}
		},
	}

	for name, test := range tests {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		t.Run(name, test)
		PostTestCleanup()
	}
}
