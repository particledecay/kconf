package build

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKubeconfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Build Suite")
}
