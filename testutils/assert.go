package testutils

import (
	"testing"
)

func Assert(t *testing.T, condition bool, msg string) {
	t.Helper()

	if !condition {
		t.Errorf("assert error: %s", msg)
	}
}

func Equals[T comparable](t *testing.T, expected, actual T) {
	t.Helper()

	if expected != actual {
		t.Errorf("want: %v; got: %v", expected, actual)
	}
}

func ErrorEquals(t *testing.T, expected, actual error) {
	t.Helper()

	// TODO: update this to use errors.Is
	if expected != nil {
		Error(t, actual)
	} else {
		Ok(t, actual)
	}
}

func Ok(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("want: %v; got: %v", nil, err)
	}
}

func Error(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf("want: error; got: %v", err)
	}
}
