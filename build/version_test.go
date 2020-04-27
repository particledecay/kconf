package build

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Build/PrintVersion", func() {
	It("Should print nothing if a version is not set", func() {
		// redirect stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// all this function does is print to stdout
		Version = ""
		PrintVersion()

		// read captured stdout
		w.Close()
		out, _ := ioutil.ReadAll(r)

		// restore stdout
		os.Stdout = oldStdout

		Expect(out).To(BeEmpty())
	})

	It("Should print a version if it was set", func() {
		// redirect stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// all this function does is print to stdout
		Version = "1.2.3"
		PrintVersion()

		// read captured stdout
		w.Close()
		out, _ := ioutil.ReadAll(r)

		// restore stdout
		os.Stdout = oldStdout

		Expect(string(out)).To(Equal("v1.2.3\n"))
	})
})

var _ = Describe("Build/PrintLongVersion", func() {
	It("Should print nothing if a version is not set", func() {
		// redirect stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// all this function does is print to stdout
		Version = ""
		PrintLongVersion()

		// read captured stdout
		w.Close()
		out, _ := ioutil.ReadAll(r)

		// restore stdout
		os.Stdout = oldStdout

		Expect(out).To(BeEmpty())
	})

	It("Should print a version if it was set", func() {
		// redirect stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// all this function does is print to stdout
		Version = "1.2.3"
		Commit = "abcdef1234"
		Date = "20200101"
		PrintLongVersion()

		// read captured stdout
		w.Close()
		out, _ := ioutil.ReadAll(r)

		// restore stdout
		os.Stdout = oldStdout

		Expect(string(out)).To(ContainSubstring("Version:"))
		Expect(string(out)).To(ContainSubstring("v1.2.3"))
		Expect(string(out)).To(ContainSubstring("SHA:"))
		Expect(string(out)).To(ContainSubstring("abcdef1234"))
		Expect(string(out)).To(ContainSubstring("Built On:"))
		Expect(string(out)).To(ContainSubstring("20200101"))
	})
})
