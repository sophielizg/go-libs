package testutils

import (
	"errors"
	"testing"
)

func Assert(t *testing.T, condition bool, msg string) {
	t.Helper()

	if !condition {
		t.Errorf("assert error: %s", msg)
	}
}

func AssertTrue(t *testing.T, condition bool) {
	t.Helper()

	if !condition {
		t.Errorf("want: %v; got: %v", true, condition)
	}
}

func AssertEquals[T comparable](t *testing.T, expected, actual T) {
	t.Helper()

	if expected != actual {
		t.Errorf("want: %v; got: %v", expected, actual)
	}
}

func AssertErrorEquals(t *testing.T, expected, actual error) {
	t.Helper()

	if expected == nil {
		AssertOk(t, actual)
	} else if actual == nil || !errors.Is(actual, expected) {
		t.Errorf("want: %v; got: %v", expected, actual)
	}
}

func AssertOk(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("want: %v; got: %v", nil, err)
	}
}

func AssertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf("want: error; got: %v", err)
	}
}

func AssertNull[T any](t *testing.T, nullable *T) {
	t.Helper()

	if nullable != nil {
		t.Errorf("want: %v; got: %v", nil, nullable)
	}
}
