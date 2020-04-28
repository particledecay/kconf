package kubeconfig

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pkg/Kubeconfig/SetNamespace", func() {
	It("Should set the namespace in the given context", func() {
		k := mockConfig(1)
		testCtx := "test"
		testNamespace := "superman"

		err := k.SetNamespace(testCtx, testNamespace)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(k.Config.Contexts[testCtx].Namespace).To(Equal(testNamespace))
	})

	It("Should fail if the context does not exist", func() {
		k := mockConfig(0)
		testCtx := "test-1"
		testNamespace := "superman"

		// cannot set namespace on a non-existent context
		err := k.SetNamespace(testCtx, testNamespace)
		Expect(err).Should(HaveOccurred())

		// also make sure the context doesn't actually exist
		_, ok := k.Config.Contexts[testCtx]
		Expect(ok).To(BeFalse())
	})
})

var _ = Describe("Pkg/Kubeconfig/SetCurrentContext", func() {
	It("Should set the current context if the context exists", func() {
		k := mockConfig(5)
		contextName := "test-2"

		err := k.SetCurrentContext(contextName)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(k.CurrentContext).To(Equal(contextName))
	})

	It("Should fail if the context does not exist", func() {
		k := mockConfig(0)
		testCtx := "test-1"

		err := k.SetCurrentContext(testCtx)
		Expect(err).Should(HaveOccurred())
	})
})
