package testutils

import "testing"

func Case(t *testing.T, name string, f func(*testing.T)) {
	t.Run(name, f)
}

type TestCase[I any, O any] struct {
	Name     string
	Input    I
	Expected O
}

type Tests[I any, O any] struct {
	Cases []TestCase[I, O]
	Func  func(*testing.T, I, O)
}

func (ts *Tests[I, O]) Run(t *testing.T) {
	t.Helper()
	for _, testCase := range ts.Cases {
		Case(t, testCase.Name, func(t *testing.T) {
			ts.Func(t, testCase.Input, testCase.Expected)
		})
	}
}
