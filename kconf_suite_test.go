package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
)

func TestKconf(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kconf Suite")
}
