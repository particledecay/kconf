package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func failOnError(msg string, err error) {
	if err != nil {
		if msg != "" {
			msg = fmt.Sprintf("%s: ", msg)
		}
		log.Fatal().Err(err).Msgf("%s%v", msg, err)
		os.Exit(1)
	}
}
