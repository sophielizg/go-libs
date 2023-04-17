package datastore

import (
	"github.com/sophielizg/go-libs/datastore/mutator"
)

// A simple key-value table
type HashTable[V any, PV mutator.Mutatable[V], H any, PH mutator.Mutatable[H]] struct {
	Backend        HashTableBackendOps
	Settings       *TableSettings
	DataRowFactory mutator.MutatableFactory[V, PV]
	HashKeyFactory mutator.MutatableFactory[H, PH]
}

func (t *HashTable[V, PV, H, PH]) Init() {
	t.Settings.ApplyOption(WithDataRow[V, PV]())
	t.Settings.ApplyOption(WithHashKey[H, PH]())
}

func (t *HashTable[V, PV, H, PH]) GetSettings() *TableSettings {
	return t.Settings
}

func (t *HashTable[V, PV, H, PH]) SetBackend(tableBackend HashTableBackendOps) {
	t.Backend = tableBackend
}

type HashTableScan[V any, PV mutator.Mutatable[V], H any, PH mutator.Mutatable[H]] struct {
	DataRow PV
	HashKey PH
}

// Scans the entire hash table, holding batchSize data rows in memory at a time
func (t *HashTable[V, PV, H, PH]) Scan(batchSize int) (chan *HashTableScan[V, PV, H, PH], chan error) {
	fieldsChan, errChan := t.Backend.Scan(batchSize)
	return scan(fieldsChan, errChan, func(fields *ScanFields) (*HashTableScan[V, PV, H, PH], error) {
		var err error
		res := &HashTableScan[V, PV, H, PH]{}

		if res.DataRow, err = t.DataRowFactory.CreateFromFields(fields.DataRow); err != nil {
			return nil, err
		} else if res.HashKey, err = t.HashKeyFactory.CreateFromFields(fields.HashKey); err != nil {
			return nil, err
		}
		return res, nil
	})
}

// Retrieves the values with the specified hash keys
func (t *HashTable[V, PV, H, PH]) Get(hashKeys ...PH) ([]PV, error) {
	data, err := t.Backend.GetMultiple(t.HashKeyFactory.CreateFieldValuesList(hashKeys))
	if err != nil {
		return nil, err
	}

	return t.DataRowFactory.CreateFromFieldsList(data)
}

// Adds a value with the specified hash key
func (t *HashTable[V, PV, H, PH]) Add(hashKey PH, data PV) (PH, error) {
	hashKeys, err := t.AddMultiple([]PH{hashKey}, []PV{data})
	if err != nil {
		return t.HashKeyFactory.Create(), err
	}

	return hashKeys[0], nil
}

// Adds multiple values with specified hash keys
func (t *HashTable[V, PV, H, PH]) AddMultiple(hashKeys []PH, data []PV) ([]PH, error) {
	if len(hashKeys) != len(data) {
		return nil, InputLengthMismatchError
	}

	fieldValues, err := t.Backend.AddMultiple(
		t.HashKeyFactory.CreateFieldValuesList(hashKeys),
		t.DataRowFactory.CreateFieldValuesList(data),
	)

	if err != nil {
		return nil, err
	} else if len(fieldValues) != len(hashKeys) {
		return nil, OutputLengthMismatchError
	}

	return t.HashKeyFactory.CreateFromFieldsList(fieldValues)
}

// Updates a value with the specified hash key
func (t *HashTable[V, PV, H, PH]) Update(hashKey PH, data PV) error {
	return t.UpdateMultiple([]PH{hashKey}, []PV{data})
}

// Updates multiple values with specified hash keys
func (t *HashTable[V, PV, H, PH]) UpdateMultiple(hashKeys []PH, data []PV) error {
	if len(hashKeys) != len(data) {
		return InputLengthMismatchError
	}

	return t.Backend.UpdateMultiple(
		t.HashKeyFactory.CreateFieldValuesList(hashKeys),
		t.DataRowFactory.CreateFieldValuesList(data),
	)
}

// Deletes values with the specified hash keys
func (t *HashTable[V, PV, H, PH]) Delete(hashKeys ...PH) error {
	return t.Backend.DeleteMultiple(t.HashKeyFactory.CreateFieldValuesList(hashKeys))
}

func (t *HashTable[V, PV, H, PH]) TransferTo(newTable *HashTable[V, PV, H, PH], batchSize int) error {
	dataChan, errorChan := t.Scan(batchSize)

	return transfer(batchSize, dataChan, errorChan, func(buf []*HashTableScan[V, PV, H, PH]) error {
		dataRows := make([]PV, len(buf))
		hashKeys := make([]PH, len(buf))

		for i, val := range buf {
			dataRows[i] = val.DataRow
			hashKeys[i] = val.HashKey
		}

		_, err := newTable.AddMultiple(hashKeys, dataRows)
		return err
	})
}
