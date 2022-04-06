package kubeconfig_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
)

var _ = Describe("Pkg/Kubeconfig/SetNamespace", func() {
	It("Should set the namespace in the given context", func() {
		k := MockConfig(1)
		testCtx := "test"
		testNamespace := "superman"

		err := k.SetNamespace(testCtx, testNamespace)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(k.Config.Contexts[testCtx].Namespace).To(Equal(testNamespace))
	})

	It("Should fail if the context does not exist", func() {
		k := MockConfig(0)
		testCtx := "test-1"
		testNamespace := "superman"

		// cannot set namespace on a non-existent context
		err := k.SetNamespace(testCtx, testNamespace)
		Expect(err).Should(HaveOccurred())

		// also make sure the context doesn't actually exist
		Expect(k).NotTo(ContainContext(testCtx))
	})
})

var _ = Describe("Pkg/Kubeconfig/SetCurrentContext", func() {
	It("Should set the current context if the context exists", func() {
		k := MockConfig(5)
		contextName := "test-2"

		err := k.SetCurrentContext(contextName)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(k.CurrentContext).To(Equal(contextName))
	})

	It("Should fail if the context does not exist", func() {
		k := MockConfig(0)
		testCtx := "test-1"

		err := k.SetCurrentContext(testCtx)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("Pkg/Kubeconfig/Save", func() {

	// restore the original config path to avoid weirdness
	AfterEach(func() {
		kc.MainConfigPath = path.Join(os.Getenv("HOME"), ".kube", "config")
	})

	It("Should save a valid kubeconfig", func() {
		k := MockConfig(1)

		// create temp location for config
		tmpfile, err := ioutil.TempFile("", "config")
		if err != nil {
			panic(err)
		}
		kc.MainConfigPath = tmpfile.Name()

		// save kubeconfig to the temp location
		err = k.Save()

		Expect(err).NotTo(HaveOccurred())

		// read saved kubeconfig
		k2, err := kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k2).To(ContainContext("test"))
	})

	It("Should fail if we can't write to the kubeconfig file", func() {
		// this should not be writeable
		k := MockConfig(1)
		kc.MainConfigPath = "/sbin/kconf.config"
		err := k.Save()

		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("Pkg/Kubeconfig/Merge", func() {
	It("Should merge a kubeconfig into an existing kubeconfig", func() {
		k := MockConfig(2)
		k2 := MockConfig(1)

		k.Merge(&k2.Config, "booyah")

		Expect(k).To(ContainContext("test"))
		Expect(k).To(ContainContext("test-1"))
	})

	It("Should rename a context when one exists with the same name", func() {
		k := MockConfig(1)
		k2 := MockConfig(2)

		// we need k and k2 to be two unique configs, so delete the first config in k2
		err := k2.Remove("test")

		Expect(err).NotTo(HaveOccurred())

		k.Merge(&k2.Config, "")

		Expect(k).To(ContainContext("test"))
		Expect(k).To(ContainContext("test-1"))
	})

	It("Should rename a cluster when one exists with the same name", func() {
		k := MockConfig(1)
		k2 := MockConfig(2)

		// we need k and k2 to be two unique configs, so delete the first config in k2
		err := k2.Remove("test")

		Expect(err).NotTo(HaveOccurred())

		k2.Clusters["test"] = k2.Clusters["test-1"]
		delete(k2.Clusters, "test-1")
		k.Merge(&k2.Config, "")

		Expect(k).To(ContainContext("test"))
		Expect(k).To(ContainContext("test-1"))
		Expect(k.Clusters).To(HaveKey("test"))
		Expect(k.Clusters).To(HaveKey("test-1"))
	})

	It("Should rename a user when one exists with the same name", func() {
		k := MockConfig(1)
		k2 := MockConfig(2)

		// we need k and k2 to be two unique configs, so delete the first config in k2
		err := k2.Remove("test")

		Expect(err).NotTo(HaveOccurred())

		k2.AuthInfos["test"] = k2.AuthInfos["test-1"]
		delete(k2.AuthInfos, "test-1")
		k.Merge(&k2.Config, "")

		Expect(k).To(ContainContext("test"))
		Expect(k).To(ContainContext("test-1"))
		Expect(k.AuthInfos).To(HaveKey("test"))
		Expect(k.AuthInfos).To(HaveKey("test-1"))
	})

	It("Should not merge a completely identical kubeconfig", func() {
		k := MockConfig(1)
		k2 := MockConfig(1)

		k.Merge(&k2.Config, "")

		Expect(k).To(ContainContext("test"))
		Expect(k).NotTo(ContainContext("test-1"))
		Expect(len(k.Contexts)).To(Equal(1))
	})
})
