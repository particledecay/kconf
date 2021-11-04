package kubeconfig_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	kc "github.com/particledecay/kconf/pkg/kubeconfig"
	. "github.com/particledecay/kconf/test"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var _ = Describe("Pkg/Kubeconfig/Read", func() {
	BeforeEach(func() {
		kc.MainConfigPath = path.Join(os.Getenv("HOME"), ".kube", "config")
	})

	It("Should fail if kubeconfig doesn't exist", func() {
		config, err := kc.Read("/some/nonexistent/path")

		Expect(config).To(BeNil())
		Expect(err).Should(HaveOccurred())
	})

	It("Should read and return a valid kubeconfig", func() {
		_ = MockConfig(1)
		config, err := kc.Read(kc.MainConfigPath)

		Expect(config).NotTo(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("Should read a kubeconfig from stdin", func() {
		// redirect stdin
		oldStdin := os.Stdin
		r, w, _ := os.Pipe()
		os.Stdin = r

		// write some data to stdin
		_ = MockConfig(1)
		confdata, _ := os.ReadFile(kc.MainConfigPath)
		w.Write(confdata)
		w.Close()

		// run function to read from stdin
		config, err := kc.Read("")

		// restore stdin
		os.Stdin = oldStdin

		Expect(config).NotTo(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("Pkg/Kubeconfig/List", func() {
	It("Should return all context names", func() {
		contexts := []string{}
		k := MockConfig(3)
		for context := range k.Contexts {
			contexts = append(contexts, context)
		}
		sort.Strings(contexts)
		Expect(contexts).To(Equal([]string{"test", "test-1", "test-2"}))
	})
})

var _ = Describe("Pkg/Kubeconfig/Export", func() {
	It("Should return a single usable config when given a context name", func() {
		contextName := "test-3"
		config := clientcmdapi.NewConfig()
		config.Clusters[contextName] = &clientcmdapi.Cluster{
			LocationOfOrigin:         "/home/user/.kube/config",
			Server:                   fmt.Sprintf("https://example-%s.com:6443", contextName),
			InsecureSkipTLSVerify:    true,
			CertificateAuthority:     "bbbbbbbbbbbb",
			CertificateAuthorityData: []byte("bbbbbbbbbbbb"),
		}
		config.AuthInfos[contextName] = &clientcmdapi.AuthInfo{
			LocationOfOrigin: "/home/user/.kube/config",
			Token:            fmt.Sprintf("bbbbbbbbbbbb-%s", contextName),
		}
		config.Contexts[contextName] = &clientcmdapi.Context{
			LocationOfOrigin: "/home/user/.kube/config",
			Cluster:          contextName,
			AuthInfo:         contextName,
			Namespace:        "default",
		}
		config.CurrentContext = contextName

		k := MockConfig(5)

		// extract the one config out of the mocked configs
		result, err := k.Export(contextName)

		Expect(result).To(Equal(config))
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("Should fail if context doesn't exist", func() {
		contextName := "test-7"
		k := MockConfig(5)

		result, err := k.Export(contextName)

		Expect(result).To(BeNil())
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("Pkg/Kubeconfig/GetContent", func() {
	It("Should properly convert a config into bytes that can be written and used as a separate kubeconfig", func() {
		contextName := "test-3"
		k := MockConfig(5)

		// extract the bytes content
		content, err := k.GetContent(contextName)

		Expect(content).ToNot(BeEmpty()) // this helps with the validity check below so we don't get a valid empty config
		Expect(err).ShouldNot(HaveOccurred())

		// convert back to a config as a validity check
		config, err := clientcmd.Load(content)

		Expect(config).ToNot(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("Pkg/Kubeconfig/GetConfig", func() {
	It("Should return a kubeconfig", func() {
		k := MockConfig(1)

		// create a random file to ensure a config exists
		tmpfile, err := ioutil.TempFile("", "test-config-5")
		if err != nil {
			panic(err)
		}
		kc.MainConfigPath = tmpfile.Name()
		err = k.Save()
		if err != nil {
			panic(err)
		}

		k2, err := kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k2).To(ContainContext("test"))
	})

	It("Should create a new kubeconfig if a file doesn't already exist", func() {
		// create a random file to ensure a config exists
		tmpfile, err := ioutil.TempFile("", "test-config-6")
		if err != nil {
			panic(err)
		}
		filename := tmpfile.Name()
		err = os.Remove(filename)
		if err != nil {
			panic(err)
		}
		kc.MainConfigPath = filename

		k, err := kc.GetConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(k.Contexts).To(BeEmpty())
	})
})
