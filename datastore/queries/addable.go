package queries

import "github.com/sophielizg/go-libs/datastore/mutator"

type AddableBackend interface {
	Add(entries []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error)
}

type Addable[E any, PE mutator.Mutatable[E]] struct {
	backend      AddableBackend
	entryFactory mutator.MutatableFactory[E, PE]
}

func (a *Addable[E, PE]) SetBackend(tableBackend AddableBackend) {
	a.backend = tableBackend
}

func (a *Addable[E, PE]) Add(entries ...PE) ([]PE, error) {
	entryFieldsList, err := a.backend.Add(
		a.entryFactory.CreateFieldValuesList(entries),
	)

	if err != nil {
		return nil, err
	}

	return a.entryFactory.CreateFromFieldsList(entryFieldsList)
}
