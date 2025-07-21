package mock

import "testing"

func TestMock(t *testing.T) {
	thing := "foo"
	t.Run("sub", func(t *testing.T) {
		UntilCleanup(t, Set(&thing, "bar"))
		if thing != "bar" {
			t.Error("wrong thing, got", thing, "want bar")
		}
	})
	if thing != "foo" {
		t.Error("wrong thing, got", thing, "want foo")
	}
}
