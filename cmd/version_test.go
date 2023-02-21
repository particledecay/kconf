package cmd_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/particledecay/kconf/build"
	"github.com/particledecay/kconf/cmd"
)

var _ = Describe("Cmd/VersionCmd", func() {
	It("Should display short version information", func() {
		// FIXME: Ginkgo supplies its own CLI flags which screws up stdout redirection
		// see https://github.com/onsi/ginkgo/issues/285#issuecomment-290575636
		// here we are truncating all args that aren't the command itself
		origArgs := os.Args[:]
		os.Args = os.Args[:1]

		// redirect stdout
		r, w, _ := os.Pipe()
		oldOut := os.Stdout
		os.Stdout = w

		// set a fake version
		build.Version = "1.2.3"
		build.Commit = "asdfasdfasdf"
		build.Date = "2020-01-01"

		versionCmd := cmd.VersionCmd()
		versionCmd.SilenceErrors = true
		err := versionCmd.Execute()

		Expect(err).NotTo(HaveOccurred())

		// read captured stdout
		_ = w.Close()
		out, _ := ioutil.ReadAll(r)

		// restore stdout
		os.Stdout = oldOut

		// restore args back to what they were
		os.Args = origArgs[:]

		Expect(out).To(ContainSubstring("1.2.3"))
		Expect(out).NotTo(ContainSubstring("asdfasdfasdf")) // short version doesn't contain this
	})
})
