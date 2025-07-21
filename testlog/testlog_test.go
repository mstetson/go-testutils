package testlog

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

// Test doesn't do anything but provide coverage.
func Test(t *testing.T) {
	Override(t)
	l := Logger(t)
	l.Printf("to custom logger")
	log.Printf("to standard logger")

	var buf bytes.Buffer
	Tee(t, &buf)
	log.Printf("tee to standard logger")
	if !strings.Contains(buf.String(), "tee to standard logger") {
		t.Fail()
	}
}
