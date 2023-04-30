package datastore

import (
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/datastore/queries"
)

type SortTable[K any, PK mutator.Mutatable[K], E any, PE mutator.Mutatable[E], C any, PC mutator.Mutatable[C]] struct {
	Settings *TableSettings
	*queries.Scannable[E, PE]
	*queries.Countable
	*queries.CRUDable[K, PK, E, PE]
	*queries.Sortable[K, PK, E, PE, C, PC]
	*queries.Transferable[E, PE]
}

func (t *SortTable[K, PK, E, PE, C, PC]) Init() {
	t.Settings.ApplyOption(WithEntry[E, PE]())
	t.Scannable = &queries.Scannable[E, PE]{}
	t.Countable = &queries.Countable{}
	t.CRUDable = &queries.CRUDable[K, PK, E, PE]{}
	t.Sortable = &queries.Sortable[K, PK, E, PE, C, PC]{
		SortFieldNames: t.Settings.SortFieldNames,
	}
	t.Transferable = &queries.Transferable[E, PE]{}
}

func (t *SortTable[K, PK, E, PE, C, PC]) GetSettings() *TableSettings {
	return t.Settings
}

func (t *SortTable[K, PK, E, PE, C, PC]) SetBackend(tableBackend SortTableBackendQueries) {
	t.Scannable.SetBackend(tableBackend)
	t.Countable.SetBackend(tableBackend)
	t.CRUDable.SetBackend(tableBackend)
	t.Sortable.SetBackend(tableBackend)
	t.Transferable.SetBackend(tableBackend)
}
