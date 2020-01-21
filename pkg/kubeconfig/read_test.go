package kubeconfig

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pkg/Kubeconfig/Read", func() {
	It("Should fail if kubeconfig doesn't exist", func() {
		config, err := Read("/some/nonexistent/path")

		Expect(config).To(BeNil())
		Expect(err).Should(HaveOccurred())
	})
})
