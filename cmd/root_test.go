package cmd_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/particledecay/kconf/cmd"
)

var _ = Describe("Cmd/Execute", func() {
	It("Should display help text", func() {
		// FIXME: Ginkgo supplies its own CLI flags which screws up stdout redirection
		// see https://github.com/onsi/ginkgo/issues/285#issuecomment-290575636
		// here we are truncating all args that aren't the command itself
		origArgs := os.Args[:]
		os.Args = os.Args[:1]

		// redirect stdout
		r, w, _ := os.Pipe()
		oldOut := os.Stdout
		os.Stdout = w

		cmd.Execute()

		// read captured stdout
		_ = w.Close()
		out, _ := ioutil.ReadAll(r)

		// restore stdout
		os.Stdout = oldOut

		// restore args back to what they were
		os.Args = origArgs[:]

		Expect(out).To(ContainSubstring("kconf"))
		Expect(out).To(ContainSubstring("Usage:"))
		Expect(out).To(ContainSubstring("Flags:"))
	})
})
