package cmd_test

import (
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/particledecay/kconf/cmd"
	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
)

var _ = Describe("Cmd/RemoveCmd", func() {

	// restore the original config path to avoid weirdness
	AfterEach(func() {
		kc.MainConfigPath = path.Join(os.Getenv("HOME"), ".kube", "config")
		CleanupFiles()
	})

	It("Should remove a context from the kubeconfig", func() {
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}

		// remove "test"
		ctxName := "test"
		removeCmd := cmd.RemoveCmd()
		removeCmd.SilenceErrors = true
		removeCmd.SetArgs([]string{ctxName})
		err = removeCmd.Execute()

		Expect(err).NotTo(HaveOccurred())

		// should not have the "test" context
		k, err = kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k).NotTo(ContainContext(ctxName))
	})

	It("Should fail when context doesn't exist", func() {
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}

		// remove "test-1" which doesn't exist
		ctxName := "test-1"
		removeCmd := cmd.RemoveCmd()
		removeCmd.SilenceErrors = true
		removeCmd.SetArgs([]string{ctxName})
		err = removeCmd.Execute()

		Expect(err).To(HaveOccurred())

		// should have not modified kubeconfig
		k, err = kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k).To(ContainContext("test"))
	})
})
