package datastore

import "github.com/sophielizg/go-libs/datastore/mutator"

// A key-value table that supports a partition key to sort values
type SortTable[V any, PV mutator.Mutatable[V], H any, PH mutator.Mutatable[H], S any, PS mutator.Mutatable[S], C any, PC mutator.Mutatable[C]] struct {
	Backend        SortTableBackendOps
	Settings       *TableSettings
	DataRowFactory mutator.MutatableFactory[V, PV]
	HashKeyFactory mutator.MutatableFactory[H, PH]
	SortKeyFactory mutator.MutatableFactory[S, PS]
}

func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) Init() {
	t.Settings.ApplyOption(WithDataRow[V, PV]())
	t.Settings.ApplyOption(WithHashKey[H, PH]())
	t.Settings.ApplyOption(WithSortKey[S, PS]())
}

func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) GetSettings() *TableSettings {
	return t.Settings
}

func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) SetBackend(tableBackend SortTableBackendOps) {
	t.Backend = tableBackend
}

type SortTableScan[V any, PV mutator.Mutatable[V], H any, PH mutator.Mutatable[H], S any, PS mutator.Mutatable[S]] struct {
	DataRow PV
	HashKey PH
	SortKey PS
}

// Scans the entire sort table, holding batchSize data rows in memory at a time
func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) Scan(batchSize int) (chan *SortTableScan[V, PV, H, PH, S, PS], chan error) {
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
func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) Get(hashKey PH, sortKey PS) (PV, error) {
	vals, err := t.GetMultiple([]PH{hashKey}, []PS{sortKey})
	if err != nil {
		return t.DataRowFactory.Create(), err
	}

	return vals[0], nil
}

// Retrieves the values with the specified hash and sort keys
func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) GetMultiple(hashKeys []PH, sortKeys []PS) ([]PV, error) {
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
func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) Add(hashKey PH, sortKey PS, data PV) (PH, PS, error) {
	hashKeys, sortKeys, err := t.AddMultiple([]PH{hashKey}, []PS{sortKey}, []PV{data})
	if err != nil {
		return t.HashKeyFactory.Create(), t.SortKeyFactory.Create(), err
	}

	return hashKeys[0], sortKeys[0], nil
}

// Adds multiple values with the specified hash and sort keys
func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) AddMultiple(hashKeys []PH, sortKeys []PS, data []PV) ([]PH, []PS, error) {
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
func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) Update(hashKey PH, sortKey PS, data PV) error {
	return t.UpdateMultiple([]PH{hashKey}, []PS{sortKey}, []PV{data})
}

// Updates multiple values with the specified hash and sort keys
func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) UpdateMultiple(hashKeys []PH, sortKeys []PS, data []PV) error {
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
func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) Delete(hashKey PH, sortKey PS) error {
	return t.DeleteMultiple([]PH{hashKey}, []PS{sortKey})
}

// Deletes multiple value with the specified hash and sort keys
func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) DeleteMultiple(hashKeys []PH, sortKeys []PS) error {
	if len(hashKeys) != len(sortKeys) {
		return InputLengthMismatchError
	}

	return t.Backend.DeleteMultiple(
		t.HashKeyFactory.CreateFieldValuesList(hashKeys),
		t.SortKeyFactory.CreateFieldValuesList(sortKeys),
	)
}

func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) validateComparator(comparator PC) error {
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

func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) GetWithSortComparator(hashKey PH, comparator PC) ([]PV, []PS, error) {
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

func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) UpdateWithSortKey(hashKey PH, comparator PC, data PV) error {
	if err := t.validateComparator(comparator); err != nil {
		return err
	}

	return t.Backend.UpdateWithSortComparator(
		t.HashKeyFactory.CreateFieldValues(hashKey),
		comparator.Mutator().GetFields(),
		t.DataRowFactory.CreateFieldValues(data),
	)
}

func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) DeleteWithSortKey(hashKey PH, comparator PC) error {
	if err := t.validateComparator(comparator); err != nil {
		return err
	}

	return t.Backend.DeleteWithSortComparator(
		t.HashKeyFactory.CreateFieldValues(hashKey),
		comparator.Mutator().GetFields(),
	)
}

// Transfers the data from this sort table to another sort table of the same type
func (t *SortTable[V, PV, H, PH, S, PS, C, PC]) TransferTo(newTable *SortTable[V, PV, H, PH, S, PS, C, PC], batchSize int) error {
	dataChan, errorChan := t.Scan(batchSize)

	return transfer(batchSize, dataChan, errorChan, func(buf []*SortTableScan[V, PV, H, PH, S, PS]) error {
		dataRows := make([]PV, len(buf))
		hashKeys := make([]PH, len(buf))
		sortKeys := make([]PS, len(buf))

		for i, val := range buf {
			dataRows[i] = val.DataRow
			hashKeys[i] = val.HashKey
			sortKeys[i] = val.SortKey
		}

		_, _, err := newTable.AddMultiple(hashKeys, sortKeys, dataRows)
		return err
	})
}
