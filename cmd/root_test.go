package cmd_test

import (
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/particledecay/kconf/cmd"
)

func TestRootCmd(t *testing.T) {
	var tests = map[string]func(*testing.T){
		"display help text": func(t *testing.T) {
			// redirect stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				w.Close()
				os.Stdout = oldStdout

				// read from r
				out, _ := io.ReadAll(r)

				var expected = []string{"kconf", "Usage:", "Flags:"}
				for _, want := range expected {
					if ok, _ := regexp.Match(want, out); !ok {
						t.Errorf("expected '%s' in output '%s'", want, out)
					}
				}
			}()

			cmd.Execute()
		},
	}

	for name, test := range tests {
		t.Run(name, test)
	}
}
