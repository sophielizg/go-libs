package datastore

import (
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/datastore/queries"
)

type HashTable[K any, PK mutator.Mutatable[K], E any, PE mutator.Mutatable[E]] struct {
	Settings *TableSettings
	*queries.Scannable[E, PE]
	*queries.Countable
	*queries.CRUDable[K, PK, E, PE]
	*queries.Transferable[E, PE]
}

func (t *HashTable[K, PK, E, PE]) Init() {
	t.Settings.ApplyOption(WithEntry[E, PE]())
	t.Scannable = &queries.Scannable[E, PE]{}
	t.Countable = &queries.Countable{}
	t.CRUDable = &queries.CRUDable[K, PK, E, PE]{}
	t.Transferable = &queries.Transferable[E, PE]{}
}

func (t *HashTable[K, PK, E, PE]) GetSettings() *TableSettings {
	return t.Settings
}

func (t *HashTable[K, PK, E, PE]) SetBackend(tableBackend HashTableBackendQueries) {
	t.Scannable.SetBackend(tableBackend)
	t.Countable.SetBackend(tableBackend)
	t.CRUDable.SetBackend(tableBackend)
	t.Transferable.SetBackend(tableBackend)
}
