package queries

import "github.com/sophielizg/go-libs/datastore/mutator"

type UpdateableBackend interface {
	Update(keys []mutator.MappedFieldValues) error
}

type Updateable[E any, PE mutator.Mutatable[E]] struct {
	backend      UpdateableBackend
	entryFactory mutator.MutatableFactory[E, PE]
}

func (a *Updateable[E, PE]) SetBackend(tableBackend UpdateableBackend) {
	a.backend = tableBackend
}

func (a *Updateable[E, PE]) Update(entries ...PE) error {
	return a.backend.Update(
		a.entryFactory.CreateFieldValuesList(entries),
	)
}
