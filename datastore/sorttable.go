package datastore

import "errors"

type SortTable[V DataRow, H HashKey, S SortKey] struct {
	Backend        SortTableBackend
	DataRowFactory DataRowFactory[V]
	HashKeyFactory KeyFactory[H]
	SortKeyFactory KeyFactory[S]
	Name           string
	schema         *SortTableSchema
}

func (t *SortTable[V, H, S]) getSchema() *SortTableSchema {
	if t.schema == nil {
		t.schema = &SortTableSchema{
			HashTableSchema: HashTableSchema{
				BaseTableSchema: BaseTableSchema{
					Name:                 t.Name,
					DataRowSchemaFactory: t.DataRowFactory,
				},
				HashKeySchemaFactory:  t.HashKeyFactory,
				SupportedFieldOptions: t.getSupportedFieldOptions(),
			},
			SortKeySchemaFactory: t.SortKeyFactory,
		}
	}

	return t.schema
}

func (t *SortTable[V, H, S]) getSupportedFieldOptions() SupportedOptions {
	supported := t.Backend.SupportedFieldOptions()
	if supported == nil {
		supported = DefaultSupportedFieldOptions
	}

	return supported
}

func (t *SortTable[V, H, S]) ValidateSchema() error {
	err := t.getSchema().Validate()
	if err != nil {
		return err
	}

	sortKeyFieldTypes := t.SortKeyFactory.GetFieldTypes()
	sortOrder := t.SortKeyFactory.GetSortOrder()
	if len(sortKeyFieldTypes) != len(sortOrder) {
		return errors.New("SortKey field types and sort order must have the same number of fields")
	}

	for _, field := range sortOrder {
		fieldType, ok := sortKeyFieldTypes[field]
		if !ok {
			return errors.New("SortKey.GetFieldTypes() must return all fields in sort order")
		} else if fieldType == nil {
			return errors.New("SortKey must define types for all fields in sort order")
		} else if !fieldType.IsComparable() {
			return errors.New("All SortKey fields must be comparable")
		}
	}

	return t.Backend.ValidateSchema(t.getSchema())
}

func (t *SortTable[V, H, S]) ValidateSortKey(sortKey S) error {
	foundEmpty := false
	fieldComparators := sortKey.GetComparators()
	sortOrder := t.SortKeyFactory.GetSortOrder()

	for _, fieldName := range sortOrder {
		comparators, ok := fieldComparators[fieldName]
		if !ok || comparators == nil {
			foundEmpty = true
			continue
		} else if foundEmpty {
			return errors.New("All SortKey fields on the left side must be included in compare")
		}

		for _, comparator := range comparators {
			if !ComparatorTypes[comparator] {
				return errors.New("Found invalid type in comparators")
			}
		}
	}

	return nil
}

func (t *SortTable[V, H, S]) Scan(batchSize int) (chan SortTableScan[V, H, S], chan error) {
	scanDataRowChan, scanErrorChan := t.Backend.Scan(t.getSchema(), batchSize)
	return scan(
		batchSize,
		scanDataRowChan,
		scanErrorChan,
		func(scanDataRow *SortTableScanFields) (SortTableScan[V, H, S], error) {
			var err error
			res := SortTableScan[V, H, S]{}

			res.DataRow, err = t.DataRowFactory.CreateFromFields(scanDataRow.DataRow)

			if err == nil {
				res.HashKey, err = t.HashKeyFactory.CreateFromFields(scanDataRow.HashKey)
			}

			if err == nil {
				res.SortKey, err = t.SortKeyFactory.CreateFromFields(scanDataRow.SortKey)
			}

			return res, err
		},
	)
}

func (t *SortTable[V, H, S]) Get(hashKey H, sortKey S) (V, error) {
	vals, err := t.GetMultiple([]H{hashKey}, []S{sortKey})
	if err != nil {
		return t.DataRowFactory.CreateDefault(), err
	}

	if len(vals) == 0 {
		return t.DataRowFactory.CreateDefault(), errors.New("No DataRow returned from Get")
	} else {
		return vals[0], nil
	}
}

func (t *SortTable[V, H, S]) GetMultiple(hashKeys []H, sortKeys []S) ([]V, error) {
	if len(hashKeys) != len(sortKeys) {
		return nil, errors.New("The number of HashKeys and SortKeys must match")
	}

	genericHashKeys := convertHashKeyToInterface(hashKeys...)
	genericSortKeys := convertSortKeyToInterface(sortKeys...)

	dataRowFieldsList, err := t.Backend.GetMultiple(t.getSchema(), genericHashKeys, genericSortKeys)
	if err != nil {
		return nil, err
	} else if len(dataRowFieldsList) > len(hashKeys) {
		return nil, errors.New("Datastore constraint not satisfied, more values than keys returned")
	}

	return convertDataRowFieldsToInterface(
		dataRowFieldsList,
		t.getSchema().validateDataRowFields,
		t.DataRowFactory,
	)
}

func (t *SortTable[V, H, S]) Add(hashKey H, sortKey S, data V) (H, S, error) {
	hashKeys, sortKeys, err := t.AddMultiple([]H{hashKey}, []S{sortKey}, []V{data})
	if err != nil {
		return t.HashKeyFactory.CreateDefault(), t.SortKeyFactory.CreateDefault(), err
	}

	if len(hashKeys) != 1 || len(sortKeys) != 1 {
		return t.HashKeyFactory.CreateDefault(), t.SortKeyFactory.CreateDefault(), errors.New("Must return exactly one HashKey and SortKey from Add")
	} else {
		return hashKeys[0], sortKeys[0], nil
	}
}

func (t *SortTable[V, H, S]) AddMultiple(hashKeys []H, sortKeys []S, data []V) ([]H, []S, error) {
	if len(hashKeys) != len(sortKeys) || len(hashKeys) != len(data) {
		return nil, nil, errors.New("The numbers of HashKeys, SortKeys, and DataRows must match")
	}

	genericData := convertDataRowToInterface(data...)
	genericHashKeys := convertHashKeyToInterface(hashKeys...)
	genericSortKeys := convertSortKeyToInterface(sortKeys...)

	hashKeyFieldsList, sortKeyFieldsList, err := t.Backend.AddMultiple(
		t.getSchema(), genericHashKeys, genericSortKeys, genericData)
	if err != nil {
		return nil, nil, err
	} else if len(hashKeyFieldsList) != len(hashKeys) || len(sortKeyFieldsList) != len(sortKeys) {
		return nil, nil, errors.New("Datastore constraint not satisfied, must return exactly the same number of HashKeys and SortKeys as were input")
	}

	hashKeyResults, err := convertDataRowFieldsToInterface[H](
		hashKeyFieldsList,
		t.getSchema().validateHashKeyFields,
		t.HashKeyFactory,
	)

	if err != nil {
		return nil, nil, err
	}

	sortKeyResults, err := convertDataRowFieldsToInterface[S](
		sortKeyFieldsList,
		t.getSchema().validateSortKeyFields,
		t.SortKeyFactory,
	)
	return hashKeyResults, sortKeyResults, err
}

func (t *SortTable[V, H, S]) Update(hashKey H, sortKey S, data V) error {
	return t.UpdateMultiple([]H{hashKey}, []S{sortKey}, []V{data})
}

func (t *SortTable[V, H, S]) UpdateMultiple(hashKeys []H, sortKeys []S, data []V) error {
	if len(hashKeys) != len(data) {
		return errors.New("The number of HashKeys must match the number of data values")
	}

	genericData := convertDataRowToInterface(data...)
	genericHashKeys := convertHashKeyToInterface(hashKeys...)
	genericSortKeys := convertSortKeyToInterface(sortKeys...)

	return t.Backend.UpdateMultiple(t.getSchema(), genericHashKeys, genericSortKeys, genericData)
}

func (t *SortTable[V, H, S]) Delete(hashKey H, sortKey S) error {
	return t.DeleteMultiple([]H{hashKey}, []S{sortKey})
}

func (t *SortTable[V, H, S]) DeleteMultiple(hashKeys []H, sortKeys []S) error {
	if len(hashKeys) != len(sortKeys) {
		return errors.New("The number of HashKeys and SortKeys must match")
	}

	genericHashKeys := convertHashKeyToInterface(hashKeys...)
	genericSortKeys := convertSortKeyToInterface(sortKeys...)
	return t.Backend.DeleteMultiple(t.getSchema(), genericHashKeys, genericSortKeys)
}

func (t *SortTable[V, H, S]) GetWithSortKey(hashKey H, sortKey S) ([]V, []H, error) {
	err := t.ValidateSortKey(sortKey)
	if err != nil {
		return nil, nil, err
	}

	dataRowFieldsList, hashKeyFieldsList, err := t.Backend.GetWithSortKey(t.getSchema(), hashKey, sortKey)
	if err != nil {
		return nil, nil, err
	} else if len(dataRowFieldsList) != len(hashKeyFieldsList) {
		return nil, nil, errors.New("Datastore constraint not satisfied, number of DataRows and HashKeys must be equal")
	}

	dataRowResults := make([]V, len(dataRowFieldsList))
	hashKeyResults := make([]H, len(hashKeyFieldsList))
	for i := range dataRowFieldsList {
		err = t.getSchema().validateDataRowFields(dataRowFieldsList[i])
		if err != nil {
			return nil, nil, err
		}

		dataRowResults[i], err = t.DataRowFactory.CreateFromFields(dataRowFieldsList[i])
		if err != nil {
			return nil, nil, err
		}

		err = t.getSchema().validateHashKeyFields(hashKeyFieldsList[i])
		if err != nil {
			return nil, nil, err
		}

		hashKeyResults[i], err = t.HashKeyFactory.CreateFromFields(hashKeyFieldsList[i])
		if err != nil {
			return nil, nil, err
		}
	}

	return dataRowResults, hashKeyResults, nil
}

func (t *SortTable[V, H, S]) UpdateWithSortKey(hashKey H, sortKey S, data V) error {
	err := t.ValidateSortKey(sortKey)
	if err != nil {
		return err
	}

	return t.Backend.UpdateWithSortKey(t.getSchema(), hashKey, sortKey, data)
}

func (t *SortTable[V, H, S]) DeleteWithSortKey(hashKey H, sortKey S) error {
	err := t.ValidateSortKey(sortKey)
	if err != nil {
		return err
	}

	return t.Backend.DeleteWithSortKey(t.getSchema(), hashKey, sortKey)
}
