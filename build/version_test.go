package build_test

import (
	"io"
	"os"
	"testing"

	"github.com/particledecay/kconf/build"
	. "github.com/particledecay/kconf/test"
)

func TestPrintVersion(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"print nothing if version not set": func(t *testing.T) {
			// redirect stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				w.Close()
				os.Stdout = oldStdout

				// read captured stdout
				out, _ := io.ReadAll(r)

				if len(out) != 0 {
					t.Errorf("expected: empty, got: %s", string(out))
				}
			}()

			// all this function does is print to stdout
			build.Version = ""
			build.PrintVersion()
		},
		"print a version if set": func(t *testing.T) {
			// redirect stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				w.Close()
				os.Stdout = oldStdout

				// read captured stdout
				out, _ := io.ReadAll(r)

				if string(out) != "v1.2.3\n" {
					t.Errorf("expected: v1.2.3, got: %s", string(out))
				}
			}()

			// all this function does is print to stdout
			build.Version = "1.2.3"
			build.PrintVersion()
		},
	}

	for name, test := range tests {
		t.Run(name, test)
	}
}

func TestPrintLongVersion(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"print nothing if version not set": func(t *testing.T) {
			// redirect stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				w.Close()
				os.Stdout = oldStdout

				// read captured stdout
				out, _ := io.ReadAll(r)

				if len(out) != 0 {
					t.Errorf("expected: empty, got: %s", string(out))
				}
			}()

			// all this function does is print to stdout
			build.Version = ""
			err := build.PrintLongVersion()
			if err != nil {
				t.Errorf("expected: nil, got: %v", err)
			}
		},
		"print a version if set": func(t *testing.T) {
			// redirect stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				w.Close()
				os.Stdout = oldStdout

				// read captured stdout
				out, _ := io.ReadAll(r)

				// expected substrings
				expected := []string{
					"Version",
					"v1.2.3",
					"SHA",
					"abcdef1234",
					"Built On",
					"20200101",
				}

				AssertSubstrings(t, out, expected)
			}()

			// all this function does is print to stdout
			build.Version = "1.2.3"
			build.Commit = "abcdef1234"
			build.Date = "20200101"
			err := build.PrintLongVersion()
			if err != nil {
				t.Errorf("expected: nil, got: %v", err)
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
	}
}
