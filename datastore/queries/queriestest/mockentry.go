package queriestest

import (
	"testing"

	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/testutils"
)

var (
	DataKey = "Data"
	IdKey   = "Id"
)

type MockData struct {
	Data string
}

func (d *MockData) Mutator() *mutator.FieldMutator {
	return mutator.NewFieldMutator(
		mutator.WithAddress(DataKey, &d.Data),
	)
}

type MockKey struct {
	Id string
}

func (k *MockKey) Mutator() *mutator.FieldMutator {
	return mutator.NewFieldMutator(
		mutator.WithAddress(IdKey, &k.Id),
	)
}

type MockNonKeyedEntry = fields.Entry[MockData, *MockData]

type MockKeyedEntry = fields.KeyedEntry[MockKey, *MockKey, MockData, *MockData]

func AssertMockNonKeyedEntryEquals(t *testing.T, expected, actual *MockNonKeyedEntry) {
	t.Helper()

	testutils.AssertEquals(t, expected.Data.Data, actual.Data.Data)
}

func AssertMockKeyEquals(t *testing.T, expected, actual *MockKey) {
	t.Helper()

	testutils.AssertEquals(t, expected.Id, actual.Id)
}

func AssertMockKeyedEntryEquals(t *testing.T, expected, actual *MockKeyedEntry) {
	t.Helper()

	AssertMockKeyEquals(t, expected.Key, actual.Key)
	testutils.AssertEquals(t, expected.Data.Data, actual.Data.Data)
}

func AssertMockNonKeyedEntryFieldsEqual(t *testing.T, expected, actual mutator.MappedFieldValues) {
	t.Helper()

	expectedDataVal, ok := expected[DataKey].(fields.String)
	testutils.AssertTrue(t, ok)
	actualDataVal, ok := actual[DataKey].(fields.String)
	testutils.AssertTrue(t, ok)
	testutils.AssertEquals(t, expectedDataVal, actualDataVal)
}

func AssertMockKeyFieldsEqual(t *testing.T, expected, actual mutator.MappedFieldValues) {
	t.Helper()

	expectedIdVal, ok := expected[IdKey].(fields.String)
	testutils.AssertTrue(t, ok)
	actualIdVal, ok := actual[IdKey].(fields.String)
	testutils.AssertTrue(t, ok)
	testutils.AssertEquals(t, expectedIdVal, actualIdVal)
}

func AssertMockKeyedEntryFieldsEqual(t *testing.T, expected, actual mutator.MappedFieldValues) {
	t.Helper()

	// all key fields should be present and equal
	AssertMockKeyFieldsEqual(t, expected, actual)

	// the same fields as the non keyed entry should also be present and equal
	AssertMockNonKeyedEntryFieldsEqual(t, expected, actual)
}
