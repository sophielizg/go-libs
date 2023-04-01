package datastore

// A key-value table that supports a partition key to sort values
type SortTable[B SortTableBackendOps, V any, PV DataRow[V], H any, PH HashKey[H], S any, PS SortKey[S], C SortKeyComparator] struct {
	Backend        B
	Settings       *TableSettings
	DataRowFactory DataRowFactory[V, PV]
	HashKeyFactory DataRowFactory[H, PH]
	SortKeyFactory DataRowFactory[S, PS]
}

func (t *SortTable[B, V, PV, H, PH, S, PS, C]) GetSettings() *TableSettings {
	return t.Settings
}

func (t *SortTable[B, V, PV, H, PH, S, PS, C]) SetBackend(tableBackend B) {
	t.Backend = tableBackend
}

// Validates the properties of a sort key against the table schema
// func (t *SortTable[B, V, PV, H, PH, S, PS, C]) ValidateSortKey(sortKey S) error {
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
func (t *SortTable[B, V, PV, H, PH, S, PS, C]) Scan(batchSize int) (chan *SortTableScan[V, PV, H, PH, S, PS], chan error) {
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
func (t *SortTable[B, V, PV, H, PH, S, PS, C]) Get(hashKey PH, sortKey PS) (PV, error) {
	vals, err := t.GetMultiple([]PH{hashKey}, []PS{sortKey})
	if err != nil {
		return t.DataRowFactory.Create(), err
	}

	return vals[0], nil
}

// Retrieves the values with the specified hash and sort keys
func (t *SortTable[B, V, PV, H, PH, S, PS, C]) GetMultiple(hashKeys []PH, sortKeys []PS) ([]PV, error) {
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
func (t *SortTable[B, V, PV, H, PH, S, PS, C]) Add(hashKey PH, sortKey PS, data PV) (PH, PS, error) {
	hashKeys, sortKeys, err := t.AddMultiple([]PH{hashKey}, []PS{sortKey}, []PV{data})
	if err != nil {
		return t.HashKeyFactory.Create(), t.SortKeyFactory.Create(), err
	}

	return hashKeys[0], sortKeys[0], nil
}

// Adds multiple values with the specified hash and sort keys
func (t *SortTable[B, V, PV, H, PH, S, PS, C]) AddMultiple(hashKeys []PH, sortKeys []PS, data []PV) ([]PH, []PS, error) {
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
func (t *SortTable[B, V, PV, H, PH, S, PS, C]) Update(hashKey PH, sortKey PS, data PV) error {
	return t.UpdateMultiple([]PH{hashKey}, []PS{sortKey}, []PV{data})
}

// Updates multiple values with the specified hash and sort keys
func (t *SortTable[B, V, PV, H, PH, S, PS, C]) UpdateMultiple(hashKeys []PH, sortKeys []PS, data []PV) error {
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
func (t *SortTable[B, V, PV, H, PH, S, PS, C]) Delete(hashKey PH, sortKey PS) error {
	return t.DeleteMultiple([]PH{hashKey}, []PS{sortKey})
}

// Deletes multiple value with the specified hash and sort keys
func (t *SortTable[B, V, PV, H, PH, S, PS, C]) DeleteMultiple(hashKeys []PH, sortKeys []PS) error {
	if len(hashKeys) != len(sortKeys) {
		return InputLengthMismatchError
	}

	return t.Backend.DeleteMultiple(
		t.HashKeyFactory.CreateFieldValuesList(hashKeys),
		t.SortKeyFactory.CreateFieldValuesList(sortKeys),
	)
}

func (t *SortTable[B, V, PV, H, PH, S, PS, C]) validateComparator(comparator C) error {
	foundEmpty := false

	for _, fieldName := range t.Settings.SortKeySettings.FieldOrder {
		if comparator.Mutator().GetField(fieldName) != nil {
			continue
		} else if foundEmpty {
			return ComparatorMissingFieldsError
		}

		foundEmpty = true
	}

	return nil
}

func (t *SortTable[B, V, PV, H, PH, S, PS, C]) GetWithSortComparator(hashKey PH, comparator C) ([]PV, []PS, error) {
	if err := t.validateComparator(comparator); err != nil {
		return nil, nil, err
	}

	dataRowFieldsList, sortKeyFieldsList, err := t.Backend.GetWithSortComparator(
		t.HashKeyFactory.CreateFieldValues(hashKey),
		comparator.Mutator().GetFields(),
	)
	if err != nil {
		return nil, nil, err
	} else if len(dataRowFieldsList) != len(sortKeyFieldsList) {
		return nil, nil, OutputLengthMismatchError
	}

	dataRowResults, err := t.DataRowFactory.CreateFromFieldsList(dataRowFieldsList)
	if err != nil {
		return nil, nil, err
	}

	sortKeyResults, err := t.SortKeyFactory.CreateFromFieldsList(sortKeyFieldsList)
	return dataRowResults, sortKeyResults, err
}

func (t *SortTable[B, V, PV, H, PH, S, PS, C]) UpdateWithSortKey(hashKey PH, comparator C, data PV) error {
	if err := t.validateComparator(comparator); err != nil {
		return err
	}

	return t.Backend.UpdateWithSortComparator(
		t.HashKeyFactory.CreateFieldValues(hashKey),
		comparator.Mutator().GetFields(),
		t.DataRowFactory.CreateFieldValues(data),
	)
}

func (t *SortTable[B, V, PV, H, PH, S, PS, C]) DeleteWithSortKey(hashKey PH, comparator C) error {
	if err := t.validateComparator(comparator); err != nil {
		return err
	}

	return t.Backend.DeleteWithSortComparator(
		t.HashKeyFactory.CreateFieldValues(hashKey),
		comparator.Mutator().GetFields(),
	)
}

// Transfers the data from this sort table to another sort table of the same type
func (t *SortTable[B, V, PV, H, PH, S, PS, C]) TransferTo(newTable *SortTable[B, V, PV, H, PH, S, PS, C], batchSize int) error {
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
