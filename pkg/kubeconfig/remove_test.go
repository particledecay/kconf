package kubeconfig_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/particledecay/kconf/test"
)

var _ = Describe("Pkg/Kubeconfig/Remove", func() {
	It("Should remove a context and all unused resources", func() {
		k := MockConfig(1)

		err := k.Remove("test")

		Expect(err).NotTo(HaveOccurred())
		Expect(k).NotTo(ContainContext("test"))
	})

	It("Should not remove a user if another context is using it", func() {
		k := MockConfig(2)

		// force the second context to use the first user
		k.Contexts["test-1"].AuthInfo = "test"
		err := k.Remove("test")

		_, ok := k.AuthInfos["test"]

		Expect(err).NotTo(HaveOccurred())
		Expect(k).NotTo(ContainContext("test")) // test context not there
		Expect(ok).To(BeTrue())                 // but user should be
	})

	It("Should fail if the context doesn't exist", func() {
		k := MockConfig(1)

		err := k.Remove("test-1")

		Expect(err).To(HaveOccurred())
		Expect(k).To(ContainContext("test"))
	})
})
