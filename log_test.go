package logutil_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mozey/logutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestMarshalStack(t *testing.T) {
	// TODO Capture os.Stdout and verify log output
	logutil.SetupLogger(true)
	err := errors.WithStack(fmt.Errorf("testing"))
	log.Error().Stack().Err(err).Msg("")
}

func TestPanicHandler(t *testing.T) {
	// TODO Capture os.Stdout and verify log output
	logutil.SetupLogger(true)
	defer logutil.PanicHandler()
	panic("testing")
}

func TestConsoleWriteFalse(t *testing.T) {
	// TODO Capture os.Stdout and verify log output
	logutil.SetupLogger(false)
	err := errors.WithStack(fmt.Errorf("testing"))
	log.Error().Stack().Err(err).Msg("")
	// Must not double encode log JSON inside message property
}

func TestSetupLogger(t *testing.T) {
	tmp, err := ioutil.TempDir("", "mozey-logutil")
	require.NoError(t, err)
	defer (func() {
		_ = os.RemoveAll(tmp)
	})()

	filePath := filepath.Join(tmp, "test.log")
	f, err := os.OpenFile(filePath,
		os.O_WRONLY|os.O_CREATE, 0644)
	logutil.SetupLogger(true, f)

	// Logs will be written to
	// stdOut using ConsoleWriter...
	err = errors.WithStack(fmt.Errorf("testing error"))
	log.Error().Stack().Err(err).Msg("")
	log.Info().Str("foo", "bar").Float64("pi", 3.14).Msg("testing info")

	// ...and to file as JSON
	b, err := ioutil.ReadFile(filePath)
	require.NoError(t, err)

	fmt.Println(string(b))
}

func TestLogToFile(t *testing.T) {
	tmp, err := ioutil.TempDir("", "mozey-logutil")
	require.NoError(t, err)
	defer (func() {
		_ = os.RemoveAll(tmp)
	})()

	filePath := filepath.Join(tmp, "test.log")
	f, err := os.OpenFile(filePath,
		os.O_WRONLY|os.O_CREATE, 0644)
	logutil.LogToFile(f)

	// Logs will be written to file as JSON
	err = errors.WithStack(fmt.Errorf("testing error"))
	log.Error().Stack().Err(err).Msg("")
	log.Info().Str("foo", "bar").Float64("pi", 3.14).Msg("testing info")

	b, err := ioutil.ReadFile(filePath)
	require.NoError(t, err)

	fmt.Println(string(b))
}
