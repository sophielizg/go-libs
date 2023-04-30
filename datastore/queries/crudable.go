package queries

import "github.com/sophielizg/go-libs/datastore/mutator"

type CRUDableBackend interface {
	GetableBackend
	AddableBackend
	UpdateableBackend
	DeleteableBackend
}

type CRUDable[K any, PK mutator.Mutatable[K], E any, PE mutator.Mutatable[E]] struct {
	Getable[K, PK, E, PE]
	Addable[E, PE]
	Updateable[E, PE]
	Deleteable[K, PK]
}

func (t *CRUDable[K, PK, E, PE]) SetBackend(tableBackend CRUDableBackend) {
	t.Getable.SetBackend(tableBackend)
	t.Addable.SetBackend(tableBackend)
	t.Updateable.SetBackend(tableBackend)
	t.Deleteable.SetBackend(tableBackend)
}
