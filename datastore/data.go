package datastore

type DataRowFields map[string]interface{}

type DataRowFieldsValidationFunc = func(DataRowFields) error

func convertDataRowFieldsToInterface[T DataRow](dataRowFieldsList []DataRowFields,
	validationFunc func(DataRowFields) error,
	dataRowFactory DataRowFactory[T]) ([]T, error) {
	dataRowResults := make([]T, len(dataRowFieldsList))
	for i, dataRowFields := range dataRowFieldsList {
		if dataRowFields == nil {
			dataRowResults[i] = dataRowFactory.CreateDefault()
			continue
		}

		err := validationFunc(dataRowFields)
		if err != nil {
			return nil, err
		}

		dataRowResults[i], err = dataRowFactory.CreateFromFields(dataRowFields)
		if err != nil {
			return nil, err
		}
	}

	return dataRowResults, nil
}

type DataRowFieldTypes map[string]FieldType

type DataRow interface {
	GetFields() DataRowFields
}

type DataRowScanFields struct {
	DataRow DataRowFields
}

type DataRowScan[V DataRow] struct {
	DataRow V
}

func convertDataRowToInterface[T DataRow](dataRows ...T) []DataRow {
	generic := make([]DataRow, len(dataRows))
	for i := range dataRows {
		generic[i] = dataRows[i]
	}

	return generic
}

type HashKey interface {
	DataRow
}

type HashTableScanFields struct {
	DataRowScanFields
	HashKey DataRowFields
}

type HashTableScan[V DataRow, H HashKey] struct {
	DataRowScan[V]
	HashKey H
}

func convertHashKeyToInterface[T HashKey](hashKeys ...T) []HashKey {
	generic := make([]HashKey, len(hashKeys))
	for i := range hashKeys {
		generic[i] = hashKeys[i]
	}

	return generic
}

type SortKey interface {
	HashKey
	GetComparators() Options
}

type SortTableScanFields struct {
	HashTableScanFields
	SortKey DataRowFields
}

type SortTableScan[V DataRow, H HashKey, S SortKey] struct {
	HashTableScan[V, H]
	SortKey S
}

func convertSortKeyToInterface[T SortKey](sortKeys ...T) []SortKey {
	generic := make([]SortKey, len(sortKeys))
	for i := range sortKeys {
		generic[i] = sortKeys[i]
	}

	return generic
}
