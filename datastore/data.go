package datastore

import "github.com/sophielizg/go-libs/datastore/fields"

type DataRow[V any] interface {
	*V
	Builder() *fields.DataRowBuilder
}

type HashKey[H any] interface {
	DataRow[H]
}

type SortKey[S any] interface {
	HashKey[S]
	// GetComparators() Options
}

type DataRowFactory[V any, PV DataRow[V]] struct{}

func (f DataRowFactory[V, PV]) Create() PV {
	return PV(new(V))
}

func (f DataRowFactory[V, PV]) CreateFromFields(fields fields.MappedFieldValues) (PV, error) {
	dataRow := f.Create()
	err := dataRow.Builder().SetFields(fields)
	return dataRow, err
}

func (f DataRowFactory[V, PV]) CreateFieldValues(dataRow PV) fields.MappedFieldValues {
	return dataRow.Builder().GetFields()
}

func (f DataRowFactory[V, PV]) CreateFromFieldsList(fieldsList []fields.MappedFieldValues) ([]PV, error) {
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

func (f DataRowFactory[V, PV]) CreateFieldValuesList(dataRows []PV) []fields.MappedFieldValues {
	fieldsList := make([]fields.MappedFieldValues, len(dataRows))

	for i, dataRow := range dataRows {
		fieldsList[i] = f.CreateFieldValues(dataRow)
	}

	return fieldsList
}
