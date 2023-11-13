package cmd_test

import (
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/particledecay/kconf/build"
	"github.com/particledecay/kconf/cmd"
)

func TestVersionCmd(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"display version information": func(t *testing.T) {
			// redirect stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				w.Close()
				os.Stdout = oldStdout

				// read from r
				out, _ := io.ReadAll(r)

				want := "1.2.3"
				if ok, _ := regexp.Match(want, out); !ok {
					t.Errorf("expected '%s' in output '%s'", want, out)
				}

				// this shouldn't be in short version output
				dont_want := "asdfasdfasdf"
				if ok, _ := regexp.Match(dont_want, out); ok {
					t.Errorf("expected '%s' to not be in output '%s'", dont_want, out)
				}
			}()

			// set a fake version
			build.Version = "1.2.3"
			build.Commit = "asdfasdfasdf"
			build.Date = "2020-01-01"

			versionCmd := cmd.VersionCmd()
			versionCmd.SilenceErrors = true

			err := versionCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
	}
}
