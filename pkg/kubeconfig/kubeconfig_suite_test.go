package kubeconfig_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"

	. "github.com/particledecay/kconf/test"
)

func TestKubeconfig(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kubeconfig Suite")
	AfterSuite(CleanupFiles)
}
