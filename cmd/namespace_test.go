package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/particledecay/kconf/cmd"
	. "github.com/particledecay/kconf/test"
)

var _ = Describe("Cmd/NamespaceCmd", func() {
	It("Should fail if a namespace is not provided", func() {
		namespaceCmd := cmd.NamespaceCmd()
		namespaceCmd.SilenceErrors = true
		namespaceCmd.SilenceUsage = true
		err := namespaceCmd.Execute()

		Expect(err).To(HaveOccurred())
	})

	It("Should fail if a current context has not been set", func() {
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}

		namespaceCmd := cmd.NamespaceCmd()
		namespaceCmd.SilenceErrors = true
		namespaceCmd.SilenceUsage = true
		namespaceCmd.SetArgs([]string{"default"})
		err = namespaceCmd.Execute()

		Expect(err).To(HaveOccurred())
	})

	It("Should fail if the namespace is blank", func() {
		k := MockConfig(2)
		k.SetCurrentContext("test-1")
		err := k.Save()
		if err != nil {
			panic(err)
		}

		namespaceCmd := cmd.NamespaceCmd()
		namespaceCmd.SilenceErrors = true
		namespaceCmd.SilenceUsage = true
		namespaceCmd.SetArgs([]string{""})
		err = namespaceCmd.Execute()

		Expect(err).To(HaveOccurred())
	})

	It("Should set the desired namespace", func() {
		k := MockConfig(2)
		k.SetCurrentContext("test-1")
		err := k.Save()
		if err != nil {
			panic(err)
		}

		namespaceCmd := cmd.NamespaceCmd()
		namespaceCmd.SilenceErrors = true
		namespaceCmd.SilenceUsage = true
		namespaceCmd.SetArgs([]string{"kube-system"})
		err = namespaceCmd.Execute()

		Expect(err).NotTo(HaveOccurred())
	})
})
