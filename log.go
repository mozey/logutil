package logutil

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

// SetupLogger sets up logging using zerolog.
//
// Wrap new errors with WithStack
//		errors.WithStack(fmt.Errorf("foo"))
//
// Errors returned by build-in or third party packages
// should be wrapped using `errors.WithStack`.
// Avoid excessive use of `errors.Wrap`,
// it's not as useful as a stack trace,
// and makes the error message harder to read.
//
// Call `.WithStack` on the boundaries of your project.
// Then don't call it again internally to the project.
// The stack trace must take you to the line where
// your project is interfacing with the vendor code
//
func SetupLogger(consoleWriter bool) {
	// Prod
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = "created"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorFieldName = "message"
	zerolog.ErrorStackMarshaler = MarshalStack
	log.Logger = log.With().Caller().Logger()

	if consoleWriter {
		// Dev
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(ConsoleWriter{
			Out:           os.Stderr,
			NoColor:       false,
			TimeFormat:    "2006-01-02 15:04:05",
			MarshalIndent: true,
		})
	}
}

func PanicHandler() {
	if r := recover(); r != nil {
		err := fmt.Errorf("%s", r)
		// Use zerolog to print stack trace
		// https://github.com/rs/zerolog/pull/35
		err = errors.Wrap(err, "recovered panic")
		log.Error().Stack().Err(err).Msg("")
	}
}
