package kubeconfig_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
)

var _ = Describe("Pkg/Kubeconfig/GetRestConfig", func() {
	It("Should fail if no current context is set", func() {
		k := MockConfig(1)
		k.CurrentContext = ""

		config, err := kc.GetRestConfig(k)

		Expect(err).To(HaveOccurred())
		Expect(config).To(BeNil())
	})

	It("Should fail with an unusable kubeconfig", func() {
		k := MockConfig(1)
		k.SetCurrentContext("test")

		config, err := kc.GetRestConfig(k)

		Expect(err).To(HaveOccurred())
		Expect(config).To(BeNil())
	})
})
