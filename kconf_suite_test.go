package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKconf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kconf Suite")
}
