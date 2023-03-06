package datastore

import (
	"errors"
)

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

func (t *HashTable[V, H]) ValidateSchema() error {
	err := t.getSchema().Validate()
	if err != nil {
		return err
	}

	return t.Backend.ValidateSchema(t.getSchema())
}

func (t *HashTable[V, H]) CreateOrUpdateSchema() error {
	err := t.ValidateSchema()
	if err != nil {
		return err
	}

	return t.Backend.CreateOrUpdateSchema(t.getSchema())
}

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

func (t *HashTable[V, H]) Update(hashKey H, data V) error {
	return t.UpdateMultiple([]H{hashKey}, []V{data})
}

func (t *HashTable[V, H]) UpdateMultiple(hashKeys []H, data []V) error {
	if len(hashKeys) != len(data) {
		return errors.New("The number of HashKeys must match the number of data values")
	}

	genericData := convertDataRowToInterface(data...)
	genericKeys := convertHashKeyToInterface(hashKeys...)

	return t.Backend.UpdateMultiple(t.getSchema(), genericKeys, genericData)
}

func (t *HashTable[V, H]) Delete(hashKeys ...H) error {
	genericKeys := convertHashKeyToInterface(hashKeys...)
	return t.Backend.DeleteMultiple(t.getSchema(), genericKeys)
}
