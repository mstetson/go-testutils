package testutils

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCheckError(t *testing.T) {
	tcs := []struct {
		errIn    error
		expected string
		ret      bool
		err      string
	}{
		{
			errIn:    nil,
			expected: "",
			ret:      true,
		},
		{
			errIn:    nil,
			expected: "something",
			ret:      true,
			err:      "unexpected success; want error something\n",
		},
		{
			errIn:    fmt.Errorf("whatever"),
			expected: "what",
			ret:      false,
		},
		{
			errIn:    fmt.Errorf("whatever"),
			expected: "",
			ret:      false,
			err:      "unexpected error; got whatever want no error\n",
		},
		{
			errIn:    fmt.Errorf("whatever"),
			expected: "something",
			ret:      false,
			err:      "unexpected error; got whatever want something\n",
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			tt := new(mockTB)
			ret := CheckError(tt, tc.errIn, tc.expected)
			if ret != tc.ret {
				t.Error("return got", ret, "want", tc.ret)
			}
			if tt.err.String() != tc.err {
				t.Error("err got", tt.err.String(), "want", tc.err)
			}
		})
	}
}

func TestCheckDiff(t *testing.T) {
	tcs := []struct {
		actual   string
		expected string
		ret      bool
		err      string
		log      string
	}{
		{
			/* empty */
			ret: true,
		},
		{
			actual:   test1Actual,
			expected: test1Expected,
			err:      "differences found\n",
			log:      test1Diff,
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			tt := new(mockTB)
			ret := CheckDiff(tt, []byte(tc.actual), []byte(tc.expected))
			if ret != tc.ret {
				t.Error("return got", ret, "want", tc.ret)
			}
			if tt.err.String() != tc.err {
				t.Error("err got", tt.err.String(), "want", tc.err)
			}
			if tt.log.String() != tc.log {
				t.Error("log got", tt.log.String(), "want", tc.log)
			}
		})
	}
}

var test1Actual = `a
B
c
D
e
`

var test1Expected = `a
b
c
d
e
`

var test1Diff = `
--- actual
+++ expected
@@ -1,5 +1,5 @@
 a
-B
+b
 c
-D
+d
 e

`

type mockTB struct {
	testing.TB
	err bytes.Buffer
	log bytes.Buffer
}

func (t *mockTB) Helper() {
}

func (t *mockTB) Error(args ...any) {
	fmt.Fprintln(&t.err, args...)
}

func (t *mockTB) Logf(format string, args ...any) {
	fmt.Fprintf(&t.log, format+"\n", args...)
}
