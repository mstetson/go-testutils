/*
Package testlog provides control over logging during tests.

Log messages written to a test case's log using t.Logf
are only printed if that particular test fails.
However, if code under test logs to a log.Logger, slog.Logger,
or the defult Loggers via functions in packages log or slog,
any test failure in the package will cause all log messages
for the package to be printed – even messages from passing tests.
This is very confusing.

To avoid that, this package allows test cases
to get a Logger or override the standard logger
in a way that redirects logging output into the test log.
This way, only logs for failing tests appear in test output.

Based on code from here:
https://github.com/golang/go/issues/22513
*/
package testlog

import (
	"io"
	"log"
	"log/slog"
	"testing"
)

// Override redirects the output of the default loggers to call t.Log.
// When t and all of its subtests are complete, the original log output is restored.
func Override(t testing.TB) {
	override(t, Writer{t})
}

// Tee, like Override, redirects the output of the default loggers to call t.Log.
// A copy of all log messages is also written to w.
// When t and all of its subtests are complete, the original log output is restored.
func Tee(t testing.TB, w io.Writer) {
	override(t, teeWriter{t, w})
}

func override(t testing.TB, w io.Writer) {
	// Package slog handles the default logger for package log,
	// so we dont' have to do that here.
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(w, nil)))

	t.Cleanup(func() {
		slog.SetDefault(old)
	})
}

// Logger returns a log.Logger that writes to t.Log.
func Logger(t testing.TB) *log.Logger {
	return log.New(Writer{t}, "", log.LstdFlags)
}

// Slogger returns a slog.Logger that writes to t.Log.
func Slogger(t testing.TB) *slog.Logger {
	return slog.New(slog.NewTextHandler(Writer{t}, nil))
}

// Writer treats a test case as an io.Writer,
// passing any output to the case's Log method.
type Writer struct {
	TB testing.TB
}

func (w Writer) Write(p []byte) (n int, err error) {
	w.TB.Helper()
	w.TB.Log(string(p))
	return len(p), nil
}

// TeeWriter is like Writer but also writes to w.
type teeWriter struct {
	TB testing.TB
	w  io.Writer
}

func (w teeWriter) Write(p []byte) (n int, err error) {
	w.TB.Helper()
	w.TB.Log(string(p))
	return w.w.Write(p)
}
