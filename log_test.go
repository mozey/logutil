package logutil_test

import (
	"fmt"
	"github.com/mozey/logutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestMarshalStack(t *testing.T) {
	// TODO Capture os.Stdout and verify log output
	logutil.SetupLogger(true)
	err := errors.WithStack(fmt.Errorf("testing"))
	log.Error().Stack().Err(err).Msg("")
}

func TestPanicHandler(t *testing.T) {
	logutil.SetupLogger(true)
	defer logutil.PanicHandler()
	panic("testing")
}
