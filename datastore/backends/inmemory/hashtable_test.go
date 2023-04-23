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

	datastoretest.TestHashTableWithBackend(t, mockTable)

	mockTableBackend.Drop()
}
