# Go Test Utils

These are small utilities I use when testing Go code. There are three packages.

- `mock` has a simple structure for setting up mocks during tests.
- `testlog` helps redirect log output during tests.
- `testutils` (this directory) has helper functions for checking errors and comparing output.

## Mockable Pattern

When I test code, I like to mock out the public API of other packages I'm using *if* those packages may be slow, make network calls, or need errors simulated. I do this by adding a file to the package I'm writing called `mockable.go` that looks like this:

```go
package whatever

// Import the packages I want to mock out.
// For this example, I'll mock out a couple of parts of the http package.
import (
	"http"
)

// Here we set local function variables
// for functions and methods we want
// to be able to mock out in tests.
// The convention is that the name is unchanged,
// except leading punctuation is removed,
// and internal punctuation is replaced with _.
var (
	http_Get        = http.Get
	http_Client_Get = (*http.Client).Get
)
```

With that, wherever I would call `http.Get(url)` in the package, I replace that with `http_Get(url)`. And wherever I would call `c.Get(url)` with http.Client c, I replace that with `http_Client_Get(c, url)`. Then during test I can override the function variables as needed to simulate behaviors or synthesize return values.

```go
package whatever

import (
	"errors"
	"http"
	"testing"

	"github.com/mstetson/go-testutils/mock"
)

func TestSomething(t *testing.T) {
	testErr := errors.New("test error")
	mock.UntilCleanup(t,
		mock.Set(&http_Get, func(url string) (*http.Response, error) {
			return nil, testErr
		}))
	err := somethingThatCallsGet()
	if !errors.Is(err, testErr) {
		t.Error("Unexpected error, got", err, "want", testErr)
	}
}
```

## License

This software is released into the public domain. See LICENSE for details.

Thanks to SchoolsPLP, LLC for funding part of the work and allowing this code to be released freely.
