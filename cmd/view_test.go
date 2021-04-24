package cmd_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/particledecay/kconf/cmd"
	. "github.com/particledecay/kconf/test"
)

var _ = Describe("Cmd/ViewCmd", func() {
	It("Should print the given context to stdout", func() {
		k := MockConfig(3)
		err := k.Save()
		if err != nil {
			panic(err)
		}

		// FIXME: Ginkgo supplies its own CLI flags which screws up stdout redirection
		// see https://github.com/onsi/ginkgo/issues/285#issuecomment-290575636
		// here we are truncating all args that aren't the command itself
		origArgs := os.Args[:]
		os.Args = os.Args[:1]

		// redirect stdout
		r, w, _ := os.Pipe()
		oldStdout := os.Stdout
		os.Stdout = w

		viewCmd := cmd.ViewCmd()
		viewCmd.SilenceErrors = true
		viewCmd.SilenceUsage = true
		viewCmd.SetArgs([]string{"test-2"})
		err = viewCmd.Execute()

		Expect(err).NotTo(HaveOccurred())

		// read captured stdout
		w.Close()
		out, _ := ioutil.ReadAll(r)

		// restore stdout
		os.Stdout = oldStdout

		// restore args back to what they were
		os.Args = origArgs[:]

		Expect(out).To(ContainSubstring("test-2"))
		Expect(out).To(ContainSubstring("certificate-authority: bbbb"))
	})

	It("Should fail if no arguments are provided", func() {
		viewCmd := cmd.ViewCmd()
		viewCmd.SilenceErrors = true
		viewCmd.SilenceUsage = true
		err := viewCmd.Execute()

		Expect(err).To(HaveOccurred())
	})

	It("Should fail if viewing a context that does not exist", func() {
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}

		viewCmd := cmd.ViewCmd()
		viewCmd.SilenceErrors = true
		viewCmd.SilenceUsage = true
		viewCmd.SetArgs([]string{"test-1"})
		err = viewCmd.Execute()

		Expect(err).To(HaveOccurred())
	})
})
