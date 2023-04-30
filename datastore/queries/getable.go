package queries

import "github.com/sophielizg/go-libs/datastore/mutator"

type GetableBackend interface {
	Get(keys []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error)
}

type Getable[K any, PK mutator.Mutatable[K], E any, PE mutator.Mutatable[E]] struct {
	backend      GetableBackend
	keyFactory   mutator.MutatableFactory[K, PK]
	entryFactory mutator.MutatableFactory[E, PE]
}

func (a *Getable[K, PK, E, PE]) SetBackend(tableBackend GetableBackend) {
	a.backend = tableBackend
}

func (a *Getable[K, PK, E, PE]) Get(keys ...PK) ([]PE, error) {
	entryFieldsList, err := a.backend.Get(
		a.keyFactory.CreateFieldValuesList(keys),
	)

	if err != nil {
		return nil, err
	}

	return a.entryFactory.CreateFromFieldsList(entryFieldsList)
}
