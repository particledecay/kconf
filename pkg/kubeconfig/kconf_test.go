package kubeconfig_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	. "github.com/particledecay/kconf/test"
)

var _ = Describe("Pkg/Kubeconfig/AddContext", func() {
	It("Should add a context if it does not already exist in the kubeconfig", func() {
		context := &clientcmdapi.Context{
			LocationOfOrigin: "/home/user/.kube/config",
			Cluster:          "test",
			AuthInfo:         "test-1",
			Namespace:        "default",
		}
		k := MockConfig(0)
		testName := "testContext"
		result := k.AddContext(testName, context)

		Expect(result).To(Equal(testName))
		Expect(k.Contexts).Should(HaveKey(testName))

		Expect(k.Contexts[testName].Cluster).To(Equal("test"))
	})

	It("Should not add a context if it already exists in the kubeconfig", func() {
		context := &clientcmdapi.Context{
			LocationOfOrigin: "/home/user/.kube/config",
			Cluster:          "test",
			AuthInfo:         "test",
			Namespace:        "default",
		}
		k := MockConfig(1)
		testName := "test"
		result := k.AddContext(testName, context)

		Expect(result).To(BeEmpty())
		Expect(k.Contexts).Should(HaveKey(testName))

		Expect(k.Contexts[testName].Cluster).To(Equal("test"))
	})

	It("Should add a context if even one field is different", func() {
		context := &clientcmdapi.Context{
			LocationOfOrigin: "/home/user/.kube/config",
			Cluster:          "test-1",
			AuthInfo:         "test",
			Namespace:        "default",
		}
		k := MockConfig(1)
		testName := "test"
		result := k.AddContext(testName, context)

		Expect(result).To(Equal("test-1"))
		Expect(k.Contexts).Should(HaveKey("test-1"))

		Expect(k.Contexts["test-1"].Cluster).To(Equal("test-1"))
	})

	It("Should add a context with context name", func() {
		context := &clientcmdapi.Context{
			LocationOfOrigin: "/home/user/.kube/config",
			Cluster:          "test-1",
			AuthInfo:         "test",
			Namespace:        "default",
		}
		k := MockConfig(0)
		testName := "test"
		result := k.AddContext(testName, context)
		Expect(result).To(Equal(testName))
	})
})

var _ = Describe("Pkg/Kubeconfig/AddCluster", func() {
	It("Should add a cluster if it does not already exist in the kubeconfig", func() {
		cluster := &clientcmdapi.Cluster{
			LocationOfOrigin:         "/home/user/.kube/config",
			Server:                   "https://example-test.com:6443",
			InsecureSkipTLSVerify:    false,
			CertificateAuthority:     "bbbbbbbbbbbb",
			CertificateAuthorityData: []byte("bbbbbbbbbbbb"),
		}
		k := MockConfig(0)
		testName := "testCluster"
		result := k.AddCluster(testName, cluster)

		Expect(result).To(Equal(testName))
		Expect(k.Clusters).Should(HaveKey(testName))

		Expect(k.Clusters[testName].Server).To(Equal("https://example-test.com:6443"))
	})

	It("Should not add a cluster if it already exists in the kubeconfig", func() {
		cluster := &clientcmdapi.Cluster{
			LocationOfOrigin:         "/home/user/.kube/config",
			Server:                   "https://example-test.com:6443",
			InsecureSkipTLSVerify:    true,
			CertificateAuthority:     "bbbbbbbbbbbb",
			CertificateAuthorityData: []byte("bbbbbbbbbbbb"),
		}
		k := MockConfig(1)
		testName := "test"
		result := k.AddCluster(testName, cluster)

		Expect(result).To(BeEmpty())
		Expect(k.Clusters).Should(HaveKey(testName))

		Expect(k.Clusters[testName].Server).To(Equal("https://example-test.com:6443"))
	})

	It("Should add a cluster if even one field is different", func() {
		cluster := &clientcmdapi.Cluster{
			LocationOfOrigin:         "/home/user/.kube/config",
			Server:                   "https://example-test.com:6443",
			InsecureSkipTLSVerify:    false,
			CertificateAuthority:     "bbbbbbbbbbbb",
			CertificateAuthorityData: []byte("bbbbbbbbbbbb"),
		}
		k := MockConfig(1)
		testName := "test"
		result := k.AddCluster(testName, cluster)

		Expect(result).To(Equal("test-1"))
		Expect(k.Clusters).Should(HaveKey("test-1"))

		Expect(k.Clusters["test-1"].InsecureSkipTLSVerify).To(BeFalse())
	})
})

var _ = Describe("Pkg/Kubeconfig/AddUser", func() {
	It("Should add a user if it does not already exist in the kubeconfig", func() {
		user := &clientcmdapi.AuthInfo{
			LocationOfOrigin: "/home/user/.kube/config",
			Token:            "bbbbbbbbbbbb-test",
		}
		k := MockConfig(0)
		testName := "testUser"
		result := k.AddUser(testName, user)

		Expect(result).To(Equal(testName))
		Expect(k.AuthInfos).Should(HaveKey(testName))

		Expect(k.AuthInfos[testName].Token).To(Equal("bbbbbbbbbbbb-test"))
	})

	It("Should not add a user if it already exists in the kubeconfig", func() {
		user := &clientcmdapi.AuthInfo{
			LocationOfOrigin: "/home/user/.kube/config",
			Token:            "bbbbbbbbbbbb-test",
		}
		k := MockConfig(1)
		testName := "test"
		result := k.AddUser(testName, user)

		Expect(result).To(BeEmpty())
		Expect(k.AuthInfos).Should(HaveKey(testName))

		Expect(k.AuthInfos[testName].Token).To(Equal("bbbbbbbbbbbb-test"))
	})

	It("Should add a user if even one field is different", func() {
		user := &clientcmdapi.AuthInfo{
			LocationOfOrigin: "/home/user/.kube/config",
			Token:            "bbbbbbbbbbbb-test-1",
		}
		k := MockConfig(1)
		testName := "test"
		result := k.AddUser(testName, user)

		Expect(result).To(Equal("test-1"))
		Expect(k.AuthInfos).Should(HaveKey("test-1"))

		Expect(k.AuthInfos["test-1"].Token).To(Equal("bbbbbbbbbbbb-test-1"))
	})
})

var _ = Describe("Pkg/Kubeconfig/MoveContext", func() {
	It("Should move an existing context to a new context", func() {
		k := MockConfig(1)
		err := k.MoveContext("test", "test-44")

		Expect(err).NotTo(HaveOccurred())
		Expect(k).To(ContainContext("test-44"))
	})

	It("Should fail if moving a context that doesn't exist", func() {
		k := MockConfig(1)
		err := k.MoveContext("test-3", "test-44")

		Expect(err).To(HaveOccurred())
		Expect(k).NotTo(ContainContext("test-44"))
	})

	It("Should fail if the new context name already exists", func() {
		k := MockConfig(2)
		err := k.MoveContext("test", "test-1")

		Expect(err).To(HaveOccurred())
	})
})
