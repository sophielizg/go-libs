package datastoretest

import (
	"strconv"
	"testing"

	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/testutils"
)

// HELPERS

func GenerateEntries(numEntries int, idPrefix string) []*MockEntry {
	entries := make([]*MockEntry, numEntries)

	for i := 0; i < numEntries; i += 1 {
		entries[i] = &MockEntry{
			Key: &MockKey{
				Id: idPrefix + strconv.Itoa(i),
			},
			Data: &MockData{
				Data: strconv.Itoa(i),
			},
		}
	}

	return entries
}

// MOCKS

const (
	DataKey = "Data"
	IdKey   = "Id"
)

type MockData struct {
	Data fields.String
}

func (d *MockData) Mutator() *mutator.FieldMutator {
	return mutator.NewFieldMutator(
		mutator.WithAddress(DataKey, &d.Data),
	)
}

var MockDataSettings = &fields.RowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(DataKey, 63),
	),
	FieldOrder: fields.OrderedFieldKeys{DataKey},
}

type MockKey struct {
	Id fields.String
}

func (k *MockKey) Mutator() *mutator.FieldMutator {
	return mutator.NewFieldMutator(
		mutator.WithAddress(IdKey, &k.Id),
	)
}

var MockKeySettings = &fields.RowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(IdKey, 63),
	),
	FieldOrder: fields.OrderedFieldKeys{IdKey},
}

type MockEntry = fields.KeyedEntry[MockKey, *MockKey, MockData, *MockData]

type MockTable = datastore.HashTable[MockKey, *MockKey, MockEntry, *MockEntry]

func NewMockTable() *MockTable {
	return &MockTable{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Test"),
			datastore.WithDataSettings(MockDataSettings),
			datastore.WithKeySettings(MockKeySettings),
		),
	}
}

// TESTS

func TestHashTableCount(t *testing.T, mockTable *MockTable) {
	t.Helper()

	entries := GenerateEntries(1, "testcount")

	_, err := mockTable.Add(entries...)
	testutils.AssertOk(t, err)

	count, err := mockTable.Count()
	testutils.AssertOk(t, err)
	testutils.AssertTrue(t, count > 0)

	err = mockTable.Delete(fields.KeysOfEntries(entries)...)
	testutils.AssertOk(t, err)
}

// TODO
func TestHashTableScan(t *testing.T, mockTable *MockTable) {
	t.Helper()

	entries := GenerateEntries(1, "testcount")

	_, err := mockTable.Add(entries...)
	testutils.AssertOk(t, err)

	count, err := mockTable.Count()
	testutils.AssertOk(t, err)
	testutils.AssertTrue(t, count > 0)

	err = mockTable.Delete(fields.KeysOfEntries(entries)...)
	testutils.AssertOk(t, err)
}

func TestHashTableGet(t *testing.T, mockTable *MockTable) {
	t.Helper()

	entries := GenerateEntries(1, "testget")
	entry := entries[0]

	_, err := mockTable.Add(entry)
	testutils.AssertOk(t, err)

	actualEntries, err := mockTable.Get(entry.Key)
	testutils.AssertOk(t, err)
	testutils.AssertEquals(t, 1, len(actualEntries))
	testutils.AssertEquals(t, entry.Key.Id, actualEntries[0].Key.Id)
	testutils.AssertEquals(t, entry.Data.Data, actualEntries[0].Data.Data)

	err = mockTable.Delete(entry.Key)
	testutils.AssertOk(t, err)
}

func TestHashTableAdd(t *testing.T, mockTable *MockTable) {
	t.Helper()

	entries := GenerateEntries(1, "testadd")
	entry := entries[0]

	actualAddEntries, err := mockTable.Add(entry)
	testutils.AssertOk(t, err)
	testutils.AssertEquals(t, 1, len(actualAddEntries))
	testutils.AssertEquals(t, entry.Key.Id, actualAddEntries[0].Key.Id)
	testutils.AssertEquals(t, entry.Data.Data, actualAddEntries[0].Data.Data)

	actualGetEntries, err := mockTable.Get(entry.Key)
	testutils.AssertOk(t, err)
	testutils.AssertEquals(t, 1, len(actualGetEntries))

	err = mockTable.Delete(entry.Key)
	testutils.AssertOk(t, err)
}

func TestHashTableUpdate(t *testing.T, mockTable *MockTable) {
	t.Helper()

	entries := GenerateEntries(1, "testupdate")
	entry := entries[0]

	_, err := mockTable.Add(entry)
	testutils.AssertOk(t, err)

	entry.Data.Data = "updated"

	err = mockTable.Update(entry)
	testutils.AssertOk(t, err)

	actualEntries, err := mockTable.Get(entry.Key)
	testutils.AssertOk(t, err)
	testutils.AssertEquals(t, 1, len(actualEntries))
	testutils.AssertEquals(t, entry.Data.Data, actualEntries[0].Data.Data)

	err = mockTable.Delete(entry.Key)
	testutils.AssertOk(t, err)
}

func TestHashTableDelete(t *testing.T, mockTable *MockTable) {
	t.Helper()

	entries := GenerateEntries(1, "testdelete")
	entry := entries[0]

	_, err := mockTable.Add(entry)
	testutils.AssertOk(t, err)

	actualEntriesBeforeDelete, err := mockTable.Get(entry.Key)
	testutils.AssertOk(t, err)
	testutils.AssertEquals(t, 1, len(actualEntriesBeforeDelete))

	err = mockTable.Delete(entry.Key)
	testutils.AssertOk(t, err)

	actualEntriesAfterDelete, err := mockTable.Get(entry.Key)
	testutils.AssertOk(t, err)
	testutils.AssertEquals(t, 0, len(actualEntriesAfterDelete))
}
