package kubeconfig

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pkg/Kubeconfig/rename", func() {
	It("Should return the same name if it doesn't already exists as a Context", func() {
		testName := "alksdjflaksjdflaskjdfj" // no way this exists
		config, _ := Read(MainConfigPath)
		result := config.rename(testName, "context")

		Expect(result).To(Equal(testName))
	})
})
