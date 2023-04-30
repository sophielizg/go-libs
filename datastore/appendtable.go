package datastore

import (
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/datastore/queries"
)

type AppendTable[E any, PE mutator.Mutatable[E]] struct {
	Settings *TableSettings
	*queries.Scannable[E, PE]
	*queries.Addable[E, PE]
	*queries.Transferable[E, PE]
}

func (t *AppendTable[E, PE]) Init() {
	t.Settings.ApplyOption(WithEntry[E, PE]())
	t.Scannable = &queries.Scannable[E, PE]{}
	t.Addable = &queries.Addable[E, PE]{}
	t.Transferable = &queries.Transferable[E, PE]{}
}

func (t *AppendTable[E, PE]) GetSettings() *TableSettings {
	return t.Settings
}

func (t *AppendTable[E, PE]) SetBackend(tableBackend AppendTableBackendQueries) {
	t.Scannable.SetBackend(tableBackend)
	t.Addable.SetBackend(tableBackend)
	t.Transferable.SetBackend(tableBackend)
}
