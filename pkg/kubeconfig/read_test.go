package kubeconfig_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/particledecay/kconf/pkg/kubeconfig"
)

var _ = Describe("Pkg/Kubeconfig/Read", func() {
	Context("When the kubeconfig file doesn't exist", func() {
		It("Should create an empty kubeconfig", func() {
			config, _ := Read("/some/nonexistent/path")
			Expect(config.Contexts).Should(BeEmpty())
		})
	})
})
