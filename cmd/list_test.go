package cmd_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/particledecay/kconf/cmd"
	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var _ = Describe("Cmd/ListCmd", func() {
	BeforeEach(func() {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	})

	AfterEach(func() {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		kc.MainConfigPath = path.Join(os.Getenv("HOME"), ".kube", "config")
		CleanupFiles()
	})

	It("Should print a list of contexts", func() {
		// write a kubeconfig that we'll import later
		k := MockConfig(5)
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
		oldOut := kc.Out
		// we use zerolog to print instead of straight os.Stdout, so override that
		kc.Out = log.Output(zerolog.ConsoleWriter{Out: w, PartsExclude: []string{"time", "level"}})

		// list contexts
		listCmd := cmd.ListCmd()
		listCmd.SilenceErrors = true
		listCmd.SilenceUsage = true
		err = listCmd.Execute()

		Expect(err).NotTo(HaveOccurred())

		// read captured stdout
		w.Close()
		out, _ := ioutil.ReadAll(r)

		// restore stdout
		kc.Out = oldOut

		// restore args back to what they were
		os.Args = origArgs[:]

		Expect(out).To(ContainSubstring("test"))
		Expect(out).To(ContainSubstring("test-1"))
		Expect(out).To(ContainSubstring("test-2"))
		Expect(out).To(ContainSubstring("test-3"))
		Expect(out).To(ContainSubstring("test-4"))
	})

	It("Should mark the current context if set", func() {
		// write a kubeconfig that we'll import later
		k := MockConfig(5)
		k.SetCurrentContext("test-1")
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
		oldOut := kc.Out
		// we use zerolog to print instead of straight os.Stdout, so override that
		kc.Out = log.Output(zerolog.ConsoleWriter{Out: w, PartsExclude: []string{"time", "level"}})

		// list contexts
		listCmd := cmd.ListCmd()
		listCmd.SilenceErrors = true
		listCmd.SilenceUsage = true
		err = listCmd.Execute()

		Expect(err).NotTo(HaveOccurred())

		// read captured stdout
		w.Close()
		out, _ := ioutil.ReadAll(r)

		// restore stdout
		kc.Out = oldOut

		// restore args back to what they were
		os.Args = origArgs[:]

		Expect(out).To(ContainSubstring("test"))
		Expect(out).To(ContainSubstring("* test-1"))
		Expect(out).To(ContainSubstring("test-2"))
		Expect(out).To(ContainSubstring("test-3"))
		Expect(out).To(ContainSubstring("test-4"))
	})
})
