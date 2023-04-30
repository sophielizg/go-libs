package queries

import "github.com/sophielizg/go-libs/datastore/mutator"

type DeleteableBackend interface {
	Delete(keys []mutator.MappedFieldValues) error
}

type Deleteable[K any, PK mutator.Mutatable[K]] struct {
	backend    DeleteableBackend
	keyFactory mutator.MutatableFactory[K, PK]
}

func (a *Deleteable[K, PK]) SetBackend(tableBackend DeleteableBackend) {
	a.backend = tableBackend
}

func (a *Deleteable[K, PK]) Delete(keys ...PK) error {
	return a.backend.Delete(
		a.keyFactory.CreateFieldValuesList(keys),
	)
}
