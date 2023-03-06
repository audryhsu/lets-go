package assert

import "testing"

// Equal is a generic helper functions for making test assertions
func Equal[T comparable](t *testing.T, actual T, expected T) {
	// designates Equal() is a test helper func. Now test runner knows to report the filename and line number of the code which *called* our Equal() function
	t.Helper()

	if actual != expected {
		t.Errorf("actual %q, expected %q", actual, expected)
	}
}
