package datastore

import "errors"

// A key-value table that supports a partition key to sort values
type SortTable[B SortTableBackendOps, V any, PV DataRow[V], H any, PH HashKey[H], S any, PS SortKey[S]] struct {
	Backend        B
	Settings       *TableSettings
	DataRowFactory DataRowFactory[V, PV]
	HashKeyFactory DataRowFactory[H, PH]
	SortKeyFactory DataRowFactory[S, PS]
}

func (t *SortTable[B, V, PV, H, PH, S, PS]) GetSettings() *TableSettings {
	return t.Settings
}

func (t *SortTable[B, V, PV, H, PH, S, PS]) SetBackend(tableBackend B) {
	t.Backend = tableBackend
}

// Validates the properties of a sort key against the table schema
// func (t *SortTable[B, V, PV, H, PH, S, PS]) ValidateSortKey(sortKey S) error {
// 	foundEmpty := false
// 	fieldComparators := sortKey.GetComparators()
// 	sortOrder := t.SortKeyFactory.GetSortOrder()

// 	for _, fieldName := range sortOrder {
// 		comparators, ok := fieldComparators[fieldName]
// 		if !ok || comparators == nil {
// 			foundEmpty = true
// 			continue
// 		} else if foundEmpty {
// 			return errors.New("All SortKey fields on the left side must be included in compare")
// 		}

// 		for _, comparator := range comparators {
// 			if !isComparator(comparator) {
// 				return errors.New("Found invalid type in comparators")
// 			}
// 		}
// 	}

// 	return nil
// }

type SortTableScan[V any, PV DataRow[V], H any, PH HashKey[H], S any, PS SortKey[S]] struct {
	DataRow PV
	HashKey PH
	SortKey PS
}

// Scans the entire sort table, holding batchSize data rows in memory at a time
func (t *SortTable[B, V, PV, H, PH, S, PS]) Scan(batchSize int) (chan *SortTableScan[V, PV, H, PH, S, PS], chan error) {
	fieldsChan, errChan := t.Backend.Scan(batchSize)
	return scan(fieldsChan, errChan, func(fields *ScanFields) (*SortTableScan[V, PV, H, PH, S, PS], error) {
		var err error
		res := &SortTableScan[V, PV, H, PH, S, PS]{}

		if res.DataRow, err = t.DataRowFactory.CreateFromFields(fields.DataRow); err != nil {
			return nil, err
		} else if res.HashKey, err = t.HashKeyFactory.CreateFromFields(fields.HashKey); err != nil {
			return nil, err
		} else if res.SortKey, err = t.SortKeyFactory.CreateFromFields(fields.HashKey); err != nil {
			return nil, err
		}
		return res, nil
	})
}

// Retrieves the values with the specified hash and sort key
func (t *SortTable[B, V, PV, H, PH, S, PS]) Get(hashKey PH, sortKey PS) (PV, error) {
	vals, err := t.GetMultiple([]PH{hashKey}, []PS{sortKey})
	if err != nil {
		return t.DataRowFactory.Create(), err
	}

	return vals[0], nil
}

// Retrieves the values with the specified hash and sort keys
func (t *SortTable[B, V, PV, H, PH, S, PS]) GetMultiple(hashKeys []PH, sortKeys []PS) ([]PV, error) {
	if len(hashKeys) != len(sortKeys) {
		return nil, InputLengthMismatchError
	}

	fieldValues, err := t.Backend.GetMultiple(
		t.HashKeyFactory.CreateFieldValuesList(hashKeys),
		t.SortKeyFactory.CreateFieldValuesList(sortKeys),
	)
	if err != nil {
		return nil, err
	} else if len(fieldValues) != len(hashKeys) {
		return nil, OutputLengthMismatchError
	}

	return t.DataRowFactory.CreateFromFieldsList(fieldValues)
}

// Adds a value with the specified hash and sort key
func (t *SortTable[B, V, PV, H, PH, S, PS]) Add(hashKey PH, sortKey PS, data PV) (PH, PS, error) {
	hashKeys, sortKeys, err := t.AddMultiple([]PH{hashKey}, []PS{sortKey}, []PV{data})
	if err != nil {
		return t.HashKeyFactory.Create(), t.SortKeyFactory.Create(), err
	}

	return hashKeys[0], sortKeys[0], nil
}

// Adds multiple values with the specified hash and sort keys
func (t *SortTable[B, V, PV, H, PH, S, PS]) AddMultiple(hashKeys []PH, sortKeys []PS, data []PV) ([]PH, []PS, error) {
	if len(hashKeys) != len(sortKeys) || len(hashKeys) != len(data) {
		return nil, nil, InputLengthMismatchError
	}

	hashKeyFieldsList, sortKeyFieldsList, err := t.Backend.AddMultiple(
		t.HashKeyFactory.CreateFieldValuesList(hashKeys),
		t.SortKeyFactory.CreateFieldValuesList(sortKeys),
		t.DataRowFactory.CreateFieldValuesList(data),
	)
	if err != nil {
		return nil, nil, err
	} else if len(hashKeyFieldsList) != len(hashKeys) || len(sortKeyFieldsList) != len(sortKeys) {
		return nil, nil, OutputLengthMismatchError
	}

	hashKeyResults, err := t.HashKeyFactory.CreateFromFieldsList(hashKeyFieldsList)
	if err != nil {
		return nil, nil, err
	}

	sortKeyResults, err := t.SortKeyFactory.CreateFromFieldsList(sortKeyFieldsList)
	return hashKeyResults, sortKeyResults, err
}

// Updates a value with the specified hash and sort key
func (t *SortTable[B, V, PV, H, PH, S, PS]) Update(hashKey PH, sortKey PS, data PV) error {
	return t.UpdateMultiple([]PH{hashKey}, []PS{sortKey}, []PV{data})
}

// Updates multiple values with the specified hash and sort keys
func (t *SortTable[B, V, PV, H, PH, S, PS]) UpdateMultiple(hashKeys []PH, sortKeys []PS, data []PV) error {
	if len(hashKeys) != len(data) {
		return InputLengthMismatchError
	}

	return t.Backend.UpdateMultiple(
		t.HashKeyFactory.CreateFieldValuesList(hashKeys),
		t.SortKeyFactory.CreateFieldValuesList(sortKeys),
		t.DataRowFactory.CreateFieldValuesList(data),
	)
}

// Deletes value with the specified hash and sort key
func (t *SortTable[B, V, PV, H, PH, S, PS]) Delete(hashKey PH, sortKey PS) error {
	return t.DeleteMultiple([]PH{hashKey}, []PS{sortKey})
}

// Deletes multiple value with the specified hash and sort keys
func (t *SortTable[B, V, PV, H, PH, S, PS]) DeleteMultiple(hashKeys []PH, sortKeys []PS) error {
	if len(hashKeys) != len(sortKeys) {
		return InputLengthMismatchError
	}

	return t.Backend.DeleteMultiple(
		t.HashKeyFactory.CreateFieldValuesList(hashKeys),
		t.SortKeyFactory.CreateFieldValuesList(sortKeys),
	)
}

// TODO: Everything below here

// Retrieves multiple values by the specified hash key, filtering by the sort key and its comparators
func (t *SortTable[B, V, PV, H, PH, S, PS]) GetWithSortKey(hashKey H, sortKey S) ([]V, []H, error) {
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

// Updates multiple values by the specified hash key, filtering by the sort key and its comparators
func (t *SortTable[B, V, PV, H, PH, S, PS]) UpdateWithSortKey(hashKey H, sortKey S, data V) error {
	err := t.ValidateSortKey(sortKey)
	if err != nil {
		return err
	}

	return t.Backend.UpdateWithSortKey(t.getSchema(), hashKey, sortKey, data)
}

// Deletes multiple values by the specified hash key, filtering by the sort key and its comparators
func (t *SortTable[B, V, PV, H, PH, S, PS]) DeleteWithSortKey(hashKey H, sortKey S) error {
	err := t.ValidateSortKey(sortKey)
	if err != nil {
		return err
	}

	return t.Backend.DeleteWithSortKey(t.getSchema(), hashKey, sortKey)
}

// Transfers the data from this sort table to another sort table of the same type
func (t *SortTable[B, V, PV, H, PH, S, PS]) TransferTo(newTable *SortTable[B, V, PV, H, PH, S, PS], batchSize int) error {
	dataChan, errorChan := t.Scan(batchSize)

	for {
		dataBuf := make([]V, 0, batchSize)
		hashKeyBuf := make([]H, 0, batchSize)
		sortKeyBuf := make([]S, 0, batchSize)

	makeBuf:
		for {
			select {
			case err, more := <-errorChan:
				if !more {
					errorChan = nil
					break makeBuf
				}

				return err
			case data, more := <-dataChan:
				if !more {
					dataChan = nil
					break makeBuf
				}

				dataBuf = append(dataBuf, data.DataRow)
				hashKeyBuf = append(hashKeyBuf, data.HashKey)
				sortKeyBuf = append(sortKeyBuf, data.SortKey)
			}

			if len(dataBuf) == batchSize {
				break
			}
		}

		if len(dataBuf) > 0 {
			_, _, err := newTable.AddMultiple(hashKeyBuf, sortKeyBuf, dataBuf)
			if err != nil {
				return err
			}
		}

		if dataChan == nil && errorChan == nil {
			return nil
		}
	}
}
