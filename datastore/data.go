package datastore

import (
	"github.com/sophielizg/go-libs/datastore/mutator"
)

type DataRow[V any] interface {
	*V
	Mutator() *mutator.FieldMutator
}

type HashKey[H any] interface {
	DataRow[H]
}

type SortKey[S any] interface {
	HashKey[S]
}

type SortKeyComparator interface {
	Mutator() *mutator.FieldMutator
}

type DataRowFactory[V any, PV DataRow[V]] struct{}

func (f DataRowFactory[V, PV]) Create() PV {
	return PV(new(V))
}

func (f DataRowFactory[V, PV]) CreateFromFields(fields mutator.MappedFieldValues) (PV, error) {
	dataRow := f.Create()
	err := dataRow.Mutator().SetFields(fields)
	return dataRow, err
}

func (f DataRowFactory[V, PV]) CreateFieldValues(dataRow PV) mutator.MappedFieldValues {
	return dataRow.Mutator().GetFields()
}

func (f DataRowFactory[V, PV]) CreateFromFieldsList(fieldsList []mutator.MappedFieldValues) ([]PV, error) {
	dataRows := make([]PV, len(fieldsList))

	for i, fields := range fieldsList {
		var err error
		dataRows[i], err = f.CreateFromFields(fields)
		if err != nil {
			return nil, err
		}
	}

	return dataRows, nil
}

func (f DataRowFactory[V, PV]) CreateFieldValuesList(dataRows []PV) []mutator.MappedFieldValues {
	fieldsList := make([]mutator.MappedFieldValues, len(dataRows))

	for i, dataRow := range dataRows {
		fieldsList[i] = f.CreateFieldValues(dataRow)
	}

	return fieldsList
}
