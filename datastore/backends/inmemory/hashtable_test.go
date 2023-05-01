package inmemory_test

import (
	"testing"

	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/datastore/backends/inmemory"
	"github.com/sophielizg/go-libs/datastore/datastoretest"
	"github.com/sophielizg/go-libs/testutils"
)

func TestHashTableBackend(t *testing.T) {
	conn := inmemory.NewConnection()
	mockTable := datastoretest.NewMockTable()
	mockTableBackend := &inmemory.HashTableBackend{}

	group := datastore.NewConnectionGroup(
		datastore.WithConnection(conn),
	)
	err := group.RegisterTables(
		datastore.RegisterHashTable[*inmemory.Connection](mockTable, mockTableBackend),
	)
	testutils.AssertOk(t, err)

	testutils.Case(t, "count", func(t *testing.T) {
		datastoretest.TestHashTableCount(t, mockTable)
	})
	testutils.Case(t, "get", func(t *testing.T) {
		datastoretest.TestHashTableGet(t, mockTable)
	})
	testutils.Case(t, "add", func(t *testing.T) {
		datastoretest.TestHashTableAdd(t, mockTable)
	})
	testutils.Case(t, "update", func(t *testing.T) {
		datastoretest.TestHashTableUpdate(t, mockTable)
	})
	testutils.Case(t, "delete", func(t *testing.T) {
		datastoretest.TestHashTableDelete(t, mockTable)
	})

	mockTableBackend.Drop()
}
