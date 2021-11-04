package kubeconfig

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// for printing to stdout using zerolog
var Out zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{"time", "level"}})

// HasPipeData returns true if kconf is receiving data from stdin pipe
func HasPipeData() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
