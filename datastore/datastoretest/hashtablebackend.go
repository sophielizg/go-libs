package datastoretest

import (
	"testing"

	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/testutils"
)

// HELPERS

func AssertMockDataRowEquals(t *testing.T, expected, actual *MockDataRow) {
	t.Helper()

	testutils.AssertEquals(t, expected.Data, actual.Data)
}

func AssertMockHashKeyEquals(t *testing.T, expected, actual *MockHashKey) {
	t.Helper()

	testutils.AssertEquals(t, expected.Id, actual.Id)
}

// MOCKS

const (
	DataKey = "Data"
	IdKey   = "Id"
)

type MockDataRow struct {
	Data         fields.String
	fieldMutator *mutator.FieldMutator
}

func (v *MockDataRow) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(DataKey, &v.Data),
		)
	}

	return v.fieldMutator
}

var MockDataRowSettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(DataKey, 63),
	),
	FieldOrder: fields.OrderedFieldKeys{DataKey},
}

type MockHashKey struct {
	Id           fields.String
	fieldMutator *mutator.FieldMutator
}

func (v *MockHashKey) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(IdKey, &v.Id),
		)
	}

	return v.fieldMutator
}

var MockHashKeySettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(IdKey, 63),
	),
	FieldOrder: fields.OrderedFieldKeys{IdKey},
}

type MockTable = datastore.HashTable[MockDataRow, *MockDataRow, MockHashKey, *MockHashKey]

type MockTableScan = datastore.HashTableScan[MockDataRow, *MockDataRow, MockHashKey, *MockHashKey]

func NewMockTable() *MockTable {
	return &MockTable{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Test"),
			datastore.WithDataRowSettings(&MockDataRowSettings),
			datastore.WithHashKeySettings(&MockHashKeySettings),
		),
	}
}

// TESTS

func TestHashTableWithBackend(t *testing.T, mockTable *MockTable) {
	t.Helper()

	testKeys := []*MockHashKey{
		{Id: "test1"},
		{Id: "test2"},
		{Id: "test3"},
	}
	testData := []*MockDataRow{
		{Data: "test1"},
		{Data: "test2"},
		{Data: "test3"},
	}
	updateKey := &MockHashKey{Id: "test1"}
	updateData := &MockDataRow{Data: "test updated"}

	_, err := mockTable.AddMultiple(testKeys, testData)
	testutils.AssertOk(t, err)

	getActualData, err := mockTable.Get(testKeys...)
	testutils.AssertOk(t, err)
	testutils.AssertEquals(t, len(testData), len(getActualData))
	for i := range testData {
		AssertMockDataRowEquals(t, testData[i], getActualData[i])
	}

	err = mockTable.Update(updateKey, updateData)
	testutils.AssertOk(t, err)

	actualUpdatedData, err := mockTable.Get(updateKey)
	testutils.AssertOk(t, err)
	testutils.AssertEquals(t, 1, len(actualUpdatedData))
	AssertMockDataRowEquals(t, updateData, actualUpdatedData[0])

	err = mockTable.Delete(testKeys...)
	testutils.AssertOk(t, err)

	actualDataAfterDelete, err := mockTable.Get(testKeys...)
	testutils.AssertOk(t, err)
	testutils.AssertEquals(t, len(testData), len(actualDataAfterDelete))
	for i := range testData {
		testutils.AssertNull(t, actualDataAfterDelete[i])
	}
}
