package logutil_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mozey/logutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

// TODO Capture os.Stdout and verify log output

func TestMarshalStack(t *testing.T) {
	logutil.SetupLogger(true)
	// Creating errors this way includes a stack trace
	err := errors.Errorf("testing")
	log.Error().Stack().Err(err).Msg("")

	// Existing err can be wrapped to include a stack trace
	err = fmt.Errorf("testing 2")
	err = errors.WithStack(err)
	log.Error().Stack().Err(err).Msg("")
}

func TestPanicHandler(t *testing.T) {
	logutil.SetupLogger(true)
	defer logutil.PanicHandler()
	panic("testing")
}

func TestConsoleWriterFalse(t *testing.T) {
	logutil.SetupLogger(false)
	err := errors.Errorf("testing")
	log.Error().Stack().Err(err).Msg("")
	// Must not double encode log JSON inside message property
}

func TestSetupLogger(t *testing.T) {
	tmp, err := os.MkdirTemp("", "mozey-logutil")
	require.NoError(t, err)
	defer (func() {
		_ = os.RemoveAll(tmp)
	})()

	filePath := filepath.Join(tmp, "test.log")
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	require.NoError(t, err)
	logutil.SetupLogger(true, f)

	// Logs will be written to
	// stdOut using ConsoleWriter...
	err = errors.Errorf("testing error")
	log.Error().Stack().Err(err).Msg("")
	log.Info().Str("foo", "bar").Float64("pi", 3.14).Msg("testing info")

	// ...and to file as JSON
	b, err := os.ReadFile(filePath)
	require.NoError(t, err)

	fmt.Println("\n--- File content\n", string(b))
}

func TestLogToFile(t *testing.T) {
	tmp, err := os.MkdirTemp("", "mozey-logutil")
	require.NoError(t, err)
	defer (func() {
		_ = os.RemoveAll(tmp)
	})()

	filePath := filepath.Join(tmp, "test.log")
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	require.NoError(t, err)
	logutil.LogToFile(f)

	// Logs will be written to file as JSON
	err = errors.Errorf("testing error")
	log.Error().Stack().Err(err).Msg("")
	log.Info().Str("foo", "bar").Float64("pi", 3.14).Msg("testing info")

	b, err := os.ReadFile(filePath)
	require.NoError(t, err)

	fmt.Println("\n--- File content\n", string(b))
}

func TestConsoleWriterToFile(t *testing.T) {
	tmp, err := os.MkdirTemp("", "mozey-logutil")
	require.NoError(t, err)
	defer (func() {
		_ = os.RemoveAll(tmp)
	})()

	filePath := filepath.Join(tmp, "test.log")
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	require.NoError(t, err)
	writer := logutil.DefaultConsoleWriter(f)
	logutil.SetupLogger(true, writer)

	// Console writer logs will be written to file
	err = errors.Errorf("testing error")
	log.Error().Stack().Err(err).Msg("")
	log.Info().Str("foo", "bar").Float64("pi", 3.14).Msg("testing info")

	b, err := os.ReadFile(filePath)
	require.NoError(t, err)

	fmt.Println("\n--- File content\n", string(b))
}

func TestConsoleWriterNoColor(t *testing.T) {
	// Default is to use colors, unless runtime.GOOS == "windows".
	logutil.SetupLogger(true)
	// This log will have color codes
	err := errors.Errorf("testing")
	log.Error().Stack().Err(err).Msg("")

	// Override by calling SetNoColor.
	// Useful when ConsoleWriter is used to write logs to a file
	logutil.SetNoColor(true)
	logutil.SetupLogger(true)
	// This log won't have color codes
	log.Error().Stack().Err(err).Msg("")
}
