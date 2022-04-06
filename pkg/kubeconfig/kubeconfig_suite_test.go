package kubeconfig_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"

	. "github.com/particledecay/kconf/test"
)

func TestKubeconfig(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	RegisterFailHandler(Fail)
	AfterSuite(CleanupFiles)
	RunSpecs(t, "Kubeconfig Suite")
}
