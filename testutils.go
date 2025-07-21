/* Package testutils holds test helpers we use across projects. */
package testutils

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

// CheckError checks whether err matches an error expectation.
// If the error expectation is an empty string, no error is expected.
// Otherwise, the error's message is expected to contain the string.
// CheckError returns true only if err is nil, expected or not.
func CheckError(t testing.TB, err error, expected string) bool {
	t.Helper()
	if err == nil && expected != "" {
		t.Error("unexpected success; want error", expected)
	}
	if err != nil && expected == "" {
		t.Error("unexpected error; got", err, "want no error")
	} else if err != nil && !strings.Contains(err.Error(), expected) {
		t.Error("unexpected error; got", err, "want", expected)
	}
	return err == nil
}

// CheckDiff checks whether actual and expected are equal.
// If they are, it returns true.
// If not, it records a test error, logs a diff, and returns false.
func CheckDiff(t testing.TB, actual, expected []byte) bool {
	t.Helper()
	strActual := strings.ReplaceAll(string(actual), "\r", "")
	strExpected := strings.ReplaceAll(string(expected), "\r", "")
	edits := myers.ComputeEdits(span.URIFromPath("actual"), strActual, strExpected)
	if len(edits) == 0 {
		return true
	}
	t.Error("differences found")
	t.Logf("\n%v", gotextdiff.ToUnified("actual", "expected", strActual, edits))
	return false
}

// CheckDeepEqual checks whether actual and expected are equal.
// If they are, it returns true.
// If not, it records a test error, logs a diff, and returns false.
func CheckDeepEqual(t testing.TB, actual, expected any) bool {
	t.Helper()
	diff := deep.Equal(actual, expected)
	if diff == nil {
		return true
	}
	t.Error("actual != expected:", diff)
	return false
}
