package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/particledecay/kconf/cmd"
	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
)

var _ = Describe("Cmd/UseCmd", func() {
	It("Should select a current context", func() {
		k := MockConfig(2)
		err := k.Save()
		if err != nil {
			panic(err)
		}

		currentContext := "test-1"
		useCmd := cmd.UseCmd()
		useCmd.SilenceErrors = true
		useCmd.SetArgs([]string{currentContext})
		err = useCmd.Execute()

		Expect(err).NotTo(HaveOccurred())

		k2, err := kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k2).To(ContainContext(currentContext))
		Expect(k2.CurrentContext).To(Equal(currentContext))
	})

	It("Should set a preferred namespace if provided", func() {
		k := MockConfig(5)
		err := k.Save()
		if err != nil {
			panic(err)
		}

		currentContext := "test-2"
		namespace := "kube-system"
		useCmd := cmd.UseCmd()
		useCmd.SilenceErrors = true
		useCmd.SetArgs([]string{currentContext, "-n", namespace})
		err = useCmd.Execute()

		Expect(err).NotTo(HaveOccurred())

		k2, err := kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k2).To(ContainContext(currentContext))
		Expect(k2.CurrentContext).To(Equal(currentContext))
		Expect(k2.Config.Contexts[currentContext].Namespace).To(Equal(namespace))
	})
})
