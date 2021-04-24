package kubeconfig_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/particledecay/kconf/test"
)

var _ = Describe("Pkg/Kubeconfig/Remove", func() {
	It("Should remove a context and all its resources", func() {
		k := MockConfig(1)

		err := k.Remove("test")

		Expect(err).NotTo(HaveOccurred())
		Expect(k).NotTo(ContainContext("test"))
	})

	It("Should fail if the context doesn't exist", func() {
		k := MockConfig(1)

		err := k.Remove("test-1")

		Expect(err).To(HaveOccurred())
		Expect(k).To(ContainContext("test"))
	})
})
