package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	// Helper marks the calling function as a test helper function.
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}
