package kubeconfig

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// for printing to stdout using zerolog
var Out zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{"time", "level"}})
