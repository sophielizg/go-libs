package datastore

// A simple key-value table
type HashTable[B HashTableBackendOps, V any, PV DataRow[V], H any, PH HashKey[H]] struct {
	Backend        B
	Settings       *TableSettings
	DataRowFactory DataRowFactory[V, PV]
	HashKeyFactory DataRowFactory[H, PH]
}

func (t *HashTable[B, V, PV, H, PH]) GetSettings() *TableSettings {
	return t.Settings
}

func (t *HashTable[B, V, PV, H, PH]) SetBackend(tableBackend B) {
	t.Backend = tableBackend
}

type HashTableScan[V any, PV DataRow[V], H any, PH HashKey[H]] struct {
	DataRow PV
	HashKey PH
}

// Scans the entire hash table, holding batchSize data rows in memory at a time
func (t *HashTable[B, V, PV, H, PH]) Scan(batchSize int) (chan *HashTableScan[V, PV, H, PH], chan error) {
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
func (t *HashTable[B, V, PV, H, PH]) Get(hashKeys ...PH) ([]PV, error) {
	data, err := t.Backend.GetMultiple(t.HashKeyFactory.CreateFieldValuesList(hashKeys))
	if err != nil {
		return nil, err
	}

	return t.DataRowFactory.CreateFromFieldsList(data)
}

// Adds a value with the specified hash key
func (t *HashTable[B, V, PV, H, PH]) Add(hashKey PH, data PV) (PH, error) {
	hashKeys, err := t.AddMultiple([]PH{hashKey}, []PV{data})
	if err != nil {
		return t.HashKeyFactory.Create(), err
	}

	return hashKeys[0], nil
}

// Adds multiple values with specified hash keys
func (t *HashTable[B, V, PV, H, PH]) AddMultiple(hashKeys []PH, data []PV) ([]PH, error) {
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
func (t *HashTable[B, V, PV, H, PH]) Update(hashKey PH, data PV) error {
	return t.UpdateMultiple([]PH{hashKey}, []PV{data})
}

// Updates multiple values with specified hash keys
func (t *HashTable[B, V, PV, H, PH]) UpdateMultiple(hashKeys []PH, data []PV) error {
	if len(hashKeys) != len(data) {
		return InputLengthMismatchError
	}

	return t.Backend.UpdateMultiple(
		t.HashKeyFactory.CreateFieldValuesList(hashKeys),
		t.DataRowFactory.CreateFieldValuesList(data),
	)
}

// Deletes values with the specified hash keys
func (t *HashTable[B, V, PV, H, PH]) Delete(hashKeys ...PH) error {
	return t.Backend.DeleteMultiple(t.HashKeyFactory.CreateFieldValuesList(hashKeys))
}
