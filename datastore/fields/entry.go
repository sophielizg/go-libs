package fields

import "github.com/sophielizg/go-libs/datastore/mutator"

type Entry[D any, PD mutator.Mutatable[D]] struct {
	Data    PD
	factory mutator.MutatableFactory[D, PD]
	mutator *mutator.FieldMutator
}

func (e *Entry[D, PD]) Mutator() *mutator.FieldMutator {
	if e.Data == nil {
		e.Data = e.factory.Create()
	}

	if e.mutator == nil {
		e.mutator = e.Data.Mutator()
	}

	return e.mutator
}

type KeyedEntry[K any, PK mutator.Mutatable[K], D any, PD mutator.Mutatable[D]] struct {
	Key         PK
	Data        PD
	keyFactory  mutator.MutatableFactory[K, PK]
	dataFactory mutator.MutatableFactory[D, PD]
	mutator     *mutator.FieldMutator
}

func (e *KeyedEntry[K, PK, D, PD]) Mutator() *mutator.FieldMutator {
	if e.Key == nil {
		e.Key = e.keyFactory.Create()
	}

	if e.Data == nil {
		e.Data = e.dataFactory.Create()
	}

	if e.mutator == nil {
		e.mutator = mutator.MergeFieldMutators(
			e.Key.Mutator(),
			e.Data.Mutator(),
		)
	}

	return e.mutator
}
