package kubeconfig

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var _ = Describe("Pkg/Kubeconfig/rename", func() {
	It("Should return the same name if it doesn't already exist as a Context", func() {
		k := MockConfig(0)
		testName := "test"
		result, err := k.rename(testName, "context")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(result).To(Equal(testName))
	})

	It("Should return an incremented name if the given name is already taken", func() {
		k := MockConfig(1)
		testName := "test"
		result, err := k.rename(testName, "context")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(result).To(Equal(fmt.Sprintf("%s-1", testName)))

		k = MockConfig(2)
		var keys []string
		for key := range k.Contexts {
			keys = append(keys, key)
		}
		result, err = k.rename(testName, "context")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(result).To(Equal(fmt.Sprintf("%s-2", testName)))
	})

	It("Should return an error if the type is not recognized", func() {
		k := MockConfig(1)
		testName := "test"
		result, err := k.rename(testName, "contextual") // contextual is not a thing

		Expect(err).Should(HaveOccurred())
		Expect(result).To(BeEmpty())
	})
})

var _ = Describe("Pkg/Kubeconfig/hasContext", func() {
	It("Should return false with an empty context", func() {
		context := &clientcmdapi.Context{
			LocationOfOrigin: "/home/user/.kube/config",
			Cluster:          "test",
			AuthInfo:         "test",
			Namespace:        "default",
		}
		k := MockConfig(0)
		result := k.hasContext(context)

		Expect(result).To(BeFalse())
	})

	It("Should return false if the context does not already exist", func() {
		context := &clientcmdapi.Context{
			LocationOfOrigin: "/home/user/.kube/config",
			Cluster:          "testing",
			AuthInfo:         "testing",
			Namespace:        "default",
		}
		k := MockConfig(1)
		result := k.hasContext(context)

		Expect(result).To(BeFalse())
	})

	It("Should return true if the context exists", func() {
		context := &clientcmdapi.Context{
			LocationOfOrigin: "/home/user/.kube/config",
			Cluster:          "test-1",
			AuthInfo:         "test-1",
			Namespace:        "default",
		}
		k := MockConfig(2)
		result := k.hasContext(context)

		Expect(result).To(BeTrue())
	})

	It("Should not return true if the context name is the same but other properties are not", func() {
		context := &clientcmdapi.Context{
			LocationOfOrigin: "/home/user/.kube/config",
			Cluster:          "test",
			AuthInfo:         "test-1",
			Namespace:        "default",
		}
		k := MockConfig(1)
		result := k.hasContext(context)

		Expect(result).To(BeFalse())
	})
})

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
		result, err := k.AddContext(testName, context)

		Expect(err).ShouldNot(HaveOccurred())
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
		result, err := k.AddContext(testName, context)

		Expect(err).ShouldNot(HaveOccurred())
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
		result, err := k.AddContext(testName, context)

		Expect(err).ShouldNot(HaveOccurred())
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
		result, err := k.AddContext(testName, context)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(result).To(Equal(testName))
	})
})

var _ = Describe("Pkg/Kubeconfig/hasCluster", func() {
	It("Should return false with an empty config", func() {
		cluster := &clientcmdapi.Cluster{
			LocationOfOrigin:         "/home/user/.kube/config",
			Server:                   "https://example-test.com:6443",
			InsecureSkipTLSVerify:    true,
			CertificateAuthority:     "bbbbbbbbbbbb",
			CertificateAuthorityData: []byte("bbbbbbbbbbbb"),
		}
		k := MockConfig(0)
		result := k.hasCluster(cluster)

		Expect(result).To(BeFalse())
	})

	It("Should return false if the cluster does not already exist", func() {
		cluster := &clientcmdapi.Cluster{
			LocationOfOrigin:         "/home/user/.kube/config",
			Server:                   "https://example-test-1.com:6443",
			InsecureSkipTLSVerify:    true,
			CertificateAuthority:     "bbbbbbbbbbbb",
			CertificateAuthorityData: []byte("bbbbbbbbbbbb"),
		}
		k := MockConfig(1)
		result := k.hasCluster(cluster)

		Expect(result).To(BeFalse())
	})

	It("Should return true if the cluster exists", func() {
		cluster := &clientcmdapi.Cluster{
			LocationOfOrigin:         "/home/user/.kube/config",
			Server:                   "https://example-test-1.com:6443",
			InsecureSkipTLSVerify:    true,
			CertificateAuthority:     "bbbbbbbbbbbb",
			CertificateAuthorityData: []byte("bbbbbbbbbbbb"),
		}
		k := MockConfig(2)
		result := k.hasCluster(cluster)

		Expect(result).To(BeTrue())
	})

	It("Should not return true if the server is the same but other properties are not", func() {
		cluster := &clientcmdapi.Cluster{
			LocationOfOrigin:         "/home/user/.kube/config",
			Server:                   "https://example-test-1.com:6443",
			InsecureSkipTLSVerify:    false,
			CertificateAuthority:     "bbbbbbbbbbbb",
			CertificateAuthorityData: []byte("bbbbbbbbbbbb"),
		}
		k := MockConfig(1)
		result := k.hasCluster(cluster)

		Expect(result).To(BeFalse())
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
		result, err := k.AddCluster(testName, cluster)

		Expect(err).ShouldNot(HaveOccurred())
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
		result, err := k.AddCluster(testName, cluster)

		Expect(err).ShouldNot(HaveOccurred())
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
		result, err := k.AddCluster(testName, cluster)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(result).To(Equal("test-1"))
		Expect(k.Clusters).Should(HaveKey("test-1"))

		Expect(k.Clusters["test-1"].InsecureSkipTLSVerify).To(BeFalse())
	})
})

var _ = Describe("Pkg/Kubeconfig/hasUser", func() {
	It("Should return false with an empty config", func() {
		user := &clientcmdapi.AuthInfo{
			LocationOfOrigin: "/home/user/.kube/config",
			Token:            "bbbbbbbbbbbb-test",
		}
		k := MockConfig(0)
		result := k.hasUser(user)

		Expect(result).To(BeFalse())
	})

	It("Should return false if the user does not already exist", func() {
		user := &clientcmdapi.AuthInfo{
			LocationOfOrigin: "/home/user/.kube/config",
			Token:            "bbbbbbbbbbbb-test-1",
		}
		k := MockConfig(1)
		result := k.hasUser(user)

		Expect(result).To(BeFalse())
	})

	It("Should return true if the user exists", func() {
		user := &clientcmdapi.AuthInfo{
			LocationOfOrigin: "/home/user/.kube/config",
			Token:            "bbbbbbbbbbbb-test-1",
		}
		k := MockConfig(2)
		result := k.hasUser(user)

		Expect(result).To(BeTrue())
	})

	It("Should not return true if the one field is the same but others are not", func() {
		user := &clientcmdapi.AuthInfo{
			LocationOfOrigin: "/home/user/.kube/CONFIG",
			Token:            "bbbbbbbbbbbb-test",
		}
		k := MockConfig(1)
		result := k.hasUser(user)

		Expect(result).To(BeFalse())
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
		result, err := k.AddUser(testName, user)

		Expect(err).ShouldNot(HaveOccurred())
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
		result, err := k.AddUser(testName, user)

		Expect(err).ShouldNot(HaveOccurred())
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
		result, err := k.AddUser(testName, user)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(result).To(Equal("test-1"))
		Expect(k.AuthInfos).Should(HaveKey("test-1"))

		Expect(k.AuthInfos["test-1"].Token).To(Equal("bbbbbbbbbbbb-test-1"))
	})
})
