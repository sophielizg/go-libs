package queries

import (
	"github.com/sophielizg/go-libs/datastore/compare"
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/utils"
)

type SortableBackend interface {
	GetWithSortComparator(key mutator.MappedFieldValues, comparator mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error)
	UpdateWithSortComparator(entry mutator.MappedFieldValues, comparator mutator.MappedFieldValues) error
	DeleteWithSortComparator(key mutator.MappedFieldValues, comparator mutator.MappedFieldValues) error
}

type Sortable[K any, PK mutator.Mutatable[K], E any, PE mutator.Mutatable[E], C any, PC mutator.Mutatable[C]] struct {
	backend        SortableBackend
	keyFactory     mutator.MutatableFactory[K, PK]
	entryFactory   mutator.MutatableFactory[E, PE]
	KeySettings    *fields.RowSettings
	SortFieldNames fields.SortFieldNames
}

func (s *Sortable[K, PK, E, PE, C, PC]) SetBackend(tableBackend SortableBackend) {
	s.backend = tableBackend
}

func (s *Sortable[K, PK, E, PE, C, PC]) validateComparator(comparator PC) error {
	foundEmpty := false

	for _, fieldName := range s.KeySettings.FieldOrder {
		if !utils.SliceContains(s.SortFieldNames, fieldName) {
			continue
		}

		if !compare.IsNilComparator(comparator.Mutator().GetField(fieldName)) {
			if foundEmpty {
				return ComparatorMissingFieldsError
			} else {
				continue
			}
		}

		foundEmpty = true
	}

	return nil
}

func (s *Sortable[K, PK, E, PE, C, PC]) GetWithSortComparator(key PK, comparator PC) ([]PE, error) {
	if err := s.validateComparator(comparator); err != nil {
		return nil, err
	}

	entryFieldsList, err := s.backend.GetWithSortComparator(
		s.keyFactory.CreateFieldValues(key),
		comparator.Mutator().GetFields(),
	)
	if err != nil {
		return nil, err
	}

	return s.entryFactory.CreateFromFieldsList(entryFieldsList)
}

func (s *Sortable[K, PK, E, PE, C, PC]) UpdateWithSortComparator(entry PE, comparator PC) error {
	if err := s.validateComparator(comparator); err != nil {
		return err
	}

	return s.backend.UpdateWithSortComparator(
		s.entryFactory.CreateFieldValues(entry),
		comparator.Mutator().GetFields(),
	)
}

func (s *Sortable[K, PK, E, PE, C, PC]) DeleteWithSortComparator(key PK, comparator PC) error {
	if err := s.validateComparator(comparator); err != nil {
		return err
	}

	return s.backend.DeleteWithSortComparator(
		s.keyFactory.CreateFieldValues(key),
		comparator.Mutator().GetFields(),
	)
}
