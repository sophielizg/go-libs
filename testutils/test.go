package testutils

import "testing"

type TestCase[I any, O any] struct {
	Name     string
	Input    I
	Expected O
}

type Tests[I any, O any] struct {
	Cases []TestCase[I, O]
	Func  func(I, O)
}

func (ts *Tests[I, O]) Run(t *testing.T) {
	t.Helper()
	for _, testCase := range ts.Cases {
		t.Log("starting test case: " + testCase.Name)
		ts.Func(testCase.Input, testCase.Expected)
	}
}
