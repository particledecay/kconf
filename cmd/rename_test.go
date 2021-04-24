package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/particledecay/kconf/cmd"
	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
)

var _ = Describe("Cmd/RenameCmd", func() {

	It("Should rename an existing context", func() {
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}

		// rename "test" -> "some-other-thing"
		oldName := "test"
		newName := "some-other-thing"
		renameCmd := cmd.RenameCmd()
		renameCmd.SilenceErrors = true
		renameCmd.SetArgs([]string{oldName, newName})
		err = renameCmd.Execute()

		Expect(err).NotTo(HaveOccurred())

		// should read new kubeconfig with new values
		k, err = kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k).NotTo(ContainContext(oldName))
		Expect(k).To(ContainContext(newName))
	})

	It("Should fail when context doesn't exist", func() {
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}

		// rename "test-jlaskdjf" -> "some-other-thing"
		oldName := "test-jlaskdjf"
		newName := "some-other-thing"
		renameCmd := cmd.RenameCmd()
		renameCmd.SilenceErrors = true
		renameCmd.SetArgs([]string{oldName, newName})
		err = renameCmd.Execute()

		Expect(err).To(HaveOccurred())

		// should read new kubeconfig with new values
		k, err = kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k).NotTo(ContainContext(newName))
	})

	It("Should fail if at least two arguments are provided", func() {
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}

		// rename "test-jlaskdjf" -> "some-other-thing"
		oldName := "test-jlaskdjf"
		newName := "some-other-thing"
		renameCmd := cmd.RenameCmd()
		renameCmd.SilenceErrors = true
		renameCmd.SilenceUsage = true
		renameCmd.SetArgs([]string{oldName})
		err = renameCmd.Execute()

		Expect(err).To(HaveOccurred())

		// should read new kubeconfig with new values
		k, err = kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k).NotTo(ContainContext(newName))
	})
})
