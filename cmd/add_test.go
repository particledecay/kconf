package cmd_test

import (
	"fmt"
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/particledecay/kconf/cmd"
	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
)

var _ = Describe("Cmd/AddCmd", func() {

	// restore the original config path to avoid weirdness
	AfterEach(func() {
		kc.MainConfigPath = path.Join(os.Getenv("HOME"), ".kube", "config")
		CleanupFiles()
	})

	It("Should add a new kubeconfig", func() {
		// create a kubeconfig to add
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}
		newConfig := kc.MainConfigPath

		// create the base kubeconfig
		tmpfile, err := MakeTmpFile()
		if err != nil {
			panic(err)
		}
		kc.MainConfigPath = tmpfile.Name()

		// add the kubeconfig
		addCmd := cmd.AddCmd()
		addCmd.SilenceErrors = true
		addCmd.SetArgs([]string{newConfig})
		err = addCmd.Execute()

		Expect(err).NotTo(HaveOccurred())

		k2, err := kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k2).To(ContainContext("test"))
	})

	It("Should rename cluster, user, and context if one already exists", func() {
		// create a kubeconfig to add
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}
		newConfig := kc.MainConfigPath

		// create the base kubeconfig with a context
		k2 := MockConfig(1)
		err = k2.Save()
		if err != nil {
			panic(err)
		}

		// add the kubeconfig
		addCmd := cmd.AddCmd()
		addCmd.SilenceErrors = true
		addCmd.SetArgs([]string{newConfig})
		err = addCmd.Execute()

		Expect(err).NotTo(HaveOccurred())

		k3, err := kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k3).To(ContainContext("test"))
		Expect(k3).To(ContainContext("test-1"))
		Expect(k3.Clusters).To(HaveKey("test"))
		Expect(k3.Clusters).To(HaveKey("test-1"))
		Expect(k3.AuthInfos).To(HaveKey("test"))
		Expect(k3.AuthInfos).To(HaveKey("test-1"))
	})

	It("Should add a new kubeconfig with a custom context name", func() {
		// create a kubeconfig to add
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}
		newConfig := kc.MainConfigPath

		// create the base kubeconfig
		tmpfile, err := MakeTmpFile()
		if err != nil {
			panic(err)
		}
		kc.MainConfigPath = tmpfile.Name()

		// add the kubeconfig
		newContext := "booyah"
		addCmd := cmd.AddCmd()
		addCmd.SilenceErrors = true
		addCmd.SetArgs([]string{
			newConfig,
			fmt.Sprintf("--context-name=%s", newContext),
		})
		err = addCmd.Execute()

		Expect(err).NotTo(HaveOccurred())

		k2, err := kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k2).To(ContainContext(newContext))
	})

	It("Should trigger a rename if a custom context name already exists", func() {
		// create a kubeconfig to add
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}
		newConfig := kc.MainConfigPath

		// create the base kubeconfig with a context
		k2 := MockConfig(2)
		err = k2.Save()
		if err != nil {
			panic(err)
		}

		// add the kubeconfig
		addCmd := cmd.AddCmd()
		addCmd.SilenceErrors = true
		addCmd.SetArgs([]string{
			newConfig,
			"--context-name=test-1",
		})
		err = addCmd.Execute()

		// Expect(err).To(HaveOccurred())
		Expect(err).NotTo(HaveOccurred())

		k3, err := kc.GetConfig()

		Expect(k3).To(ContainContext("test-1-1"))
	})

	It("Should fail if the kubeconfig is invalid", func() {
		// create a kubeconfig to add
		k := MockConfig(1)
		err := k.Save()
		if err != nil {
			panic(err)
		}
		newConfig := kc.MainConfigPath

		// create the base kubeconfig
		tmpfile, err := MakeTmpFile()
		if err != nil {
			panic(err)
		}
		tmpfile.Write([]byte("ajsdlkfjasldfkjlasdkf"))
		kc.MainConfigPath = tmpfile.Name()

		// add the kubeconfig
		addCmd := cmd.AddCmd()
		addCmd.SilenceUsage = true
		addCmd.SilenceErrors = true
		addCmd.SetArgs([]string{newConfig})
		err = addCmd.Execute()

		Expect(err).To(HaveOccurred())
	})
})
