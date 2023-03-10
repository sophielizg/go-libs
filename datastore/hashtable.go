package datastore

import (
	"errors"
)

// A simple key-value table
type HashTable[V DataRow, H HashKey] struct {
	Backend        HashTableBackend
	DataRowFactory DataRowFactory[V]
	HashKeyFactory KeyFactory[H]
	Name           string
	schema         *HashTableSchema
}

func (t *HashTable[V, H]) getSchema() *HashTableSchema {
	if t.schema == nil {
		t.schema = &HashTableSchema{
			BaseTableSchema: BaseTableSchema{
				Name:                 t.Name,
				DataRowSchemaFactory: t.DataRowFactory,
			},
			HashKeySchemaFactory:  t.HashKeyFactory,
			SupportedFieldOptions: t.getSupportedFieldOptions(),
		}
	}

	return t.schema
}

func (t *HashTable[V, H]) getSupportedFieldOptions() SupportedOptions {
	supported := t.Backend.SupportedFieldOptions()
	if supported == nil {
		supported = DefaultSupportedFieldOptions
	}

	return supported
}

// Validates the schema of the hash table
func (t *HashTable[V, H]) ValidateSchema() error {
	err := t.getSchema().Validate()
	if err != nil {
		return err
	}

	return t.Backend.ValidateSchema(t.getSchema())
}

// Creates or updates the schema of the hash table
func (t *HashTable[V, H]) CreateOrUpdateSchema() error {
	err := t.ValidateSchema()
	if err != nil {
		return err
	}

	return t.Backend.CreateOrUpdateSchema(t.getSchema())
}

// Scans the entire hash table, holding batchSize data rows in memory at a time
func (t *HashTable[V, H]) Scan(batchSize int) (chan HashTableScan[V, H], chan error) {
	scanDataRowChan, scanErrorChan := t.Backend.Scan(t.getSchema(), batchSize)
	return scan(
		batchSize,
		scanDataRowChan,
		scanErrorChan,
		func(scanDataRow *HashTableScanFields) (HashTableScan[V, H], error) {
			var err error
			res := HashTableScan[V, H]{}

			res.DataRow, err = t.DataRowFactory.CreateFromFields(scanDataRow.DataRow)

			if err == nil {
				res.HashKey, err = t.HashKeyFactory.CreateFromFields(scanDataRow.HashKey)
			}

			return res, err
		},
	)
}

// Retrieves the values with the specified hash keys
func (t *HashTable[V, H]) Get(hashKeys ...H) ([]V, error) {
	genericKeys := convertHashKeyToInterface(hashKeys...)
	dataRowFieldsList, err := t.Backend.GetMultiple(t.getSchema(), genericKeys)
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

// Adds a value with the specified hash key
func (t *HashTable[V, H]) Add(hashKey H, data V) (H, error) {
	hashKeys, err := t.AddMultiple([]H{hashKey}, []V{data})
	if err != nil {
		return t.HashKeyFactory.CreateDefault(), err
	}

	if len(hashKeys) != 1 {
		return t.HashKeyFactory.CreateDefault(), errors.New("Must return exactly one HashKey from Add")
	} else {
		return hashKeys[0], nil
	}
}

// Adds multiple values with specified hash keys
func (t *HashTable[V, H]) AddMultiple(hashKeys []H, data []V) ([]H, error) {
	if len(hashKeys) != len(data) {
		return nil, errors.New("The number of HashKeys must match the number of data values")
	}

	genericData := convertDataRowToInterface(data...)
	genericKeys := convertHashKeyToInterface(hashKeys...)

	hashKeyFieldsList, err := t.Backend.AddMultiple(t.getSchema(), genericKeys, genericData)
	if err != nil {
		return nil, err
	} else if len(hashKeyFieldsList) != len(hashKeys) {
		return nil, errors.New("Datastore constraint not satisfied, must return exactly the same number of HashKeys as were input")
	}

	return convertDataRowFieldsToInterface[H](
		hashKeyFieldsList,
		t.getSchema().validateHashKeyFields,
		t.HashKeyFactory,
	)
}

// Updates a value with the specified hash key
func (t *HashTable[V, H]) Update(hashKey H, data V) error {
	return t.UpdateMultiple([]H{hashKey}, []V{data})
}

// Updates multiple values with specified hash keys
func (t *HashTable[V, H]) UpdateMultiple(hashKeys []H, data []V) error {
	if len(hashKeys) != len(data) {
		return errors.New("The number of HashKeys must match the number of data values")
	}

	genericData := convertDataRowToInterface(data...)
	genericKeys := convertHashKeyToInterface(hashKeys...)

	return t.Backend.UpdateMultiple(t.getSchema(), genericKeys, genericData)
}

// Deletes values with the specified hash keys
func (t *HashTable[V, H]) Delete(hashKeys ...H) error {
	genericKeys := convertHashKeyToInterface(hashKeys...)
	return t.Backend.DeleteMultiple(t.getSchema(), genericKeys)
}

// Transfers the data from this hash table to another hash table of the same type
func (t *HashTable[V, H]) TransferTo(newTable *HashTable[V, H], batchSize int) error {
	dataChan, errorChan := t.Scan(batchSize)

	for {
		dataBuf := make([]V, 0, batchSize)
		hashKeyBuf := make([]H, 0, batchSize)

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
			}

			if len(dataBuf) == batchSize {
				break
			}
		}

		if len(dataBuf) > 0 {
			_, err := newTable.AddMultiple(hashKeyBuf, dataBuf)
			if err != nil {
				return err
			}
		}

		if dataChan == nil && errorChan == nil {
			return nil
		}
	}
}
