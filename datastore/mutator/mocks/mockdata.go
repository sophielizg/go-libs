package mocks

import (
	"testing"

	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/testutils"
)

type MockData struct {
	MockField1   string
	MockField2   int
	MockField3   *float32
	fieldMutator *mutator.FieldMutator
}

func (d *MockData) Mutator() *mutator.FieldMutator {
	if d.fieldMutator == nil {
		d.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress("1", &d.MockField1),
			mutator.WithAddress("2", &d.MockField2),
			mutator.WithAddress("3", &d.MockField3),
		)
	}

	return d.fieldMutator
}

func AssertMockDataEquals(t *testing.T, expected, actual *MockData) {
	t.Helper()
	testutils.AssertEquals(t, expected.MockField1, actual.MockField1)
	testutils.AssertEquals(t, expected.MockField2, actual.MockField2)
	testutils.AssertEquals(t, expected.MockField3, actual.MockField3)
}
