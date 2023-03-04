package datastore

type DataRowFields map[string]interface{}

type DataRowFieldsValidationFunc = func(DataRowFields) error

func convertDataRowFieldsToInterface[T DataRow](dataRowFieldsList []DataRowFields,
	validationFunc func(DataRowFields) error,
	dataRowFactory DataRowFactory[T]) ([]T, error) {
	dataRowResults := make([]T, len(dataRowFieldsList))
	for i, dataRowFields := range dataRowFieldsList {
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

func convertHashKeyToInterface[T HashKey](hashKeys ...T) []HashKey {
	generic := make([]HashKey, len(hashKeys))
	for i := range hashKeys {
		generic[i] = hashKeys[i]
	}

	return generic
}

func convertSortKeyToInterface[T SortKey](sortKeys ...T) []SortKey {
	generic := make([]SortKey, len(sortKeys))
	for i := range sortKeys {
		generic[i] = sortKeys[i]
	}

	return generic
}

type SortKey interface {
	HashKey
}
