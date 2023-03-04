package datastore

import (
	"errors"
)

type HashTable[V DataRow, H HashKey] struct {
	Backend        HashTableBackend
	DataRowFactory DataRowFactory[V]
	HashKeyFactory HashKeyFactory[H]
	Name           string
	schema         *HashTableSchema
}

func (t *HashTable[V, H]) getSchema() *HashTableSchema {
	if t.schema == nil {
		t.schema = &HashTableSchema{
			Name:                 t.Name,
			DataRowSchemaFactory: t.DataRowFactory,
			HashKeySchemaFactory: t.HashKeyFactory,
		}
	}

	return t.schema
}

func (t *HashTable[V, H]) getSupportedFieldOptions() SupportedFieldOptions {
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

	err = validateFieldOptions(t.HashKeyFactory, t.getSupportedFieldOptions())
	if err != nil {
		return err
	}

	return t.Backend.ValidateHashTableSchema(t.getSchema())
}

func (t *HashTable[V, H]) CreateOrUpdateSchema() error {
	err := t.ValidateSchema()
	if err != nil {
		return err
	}

	return t.Backend.CreateOrUpdateHashTableSchema(t.getSchema())
}

func (t *HashTable[V, H]) Get(hashKeys ...H) ([]V, error) {
	genericKeys := make([]HashKey, len(hashKeys))
	for i := range hashKeys {
		genericKeys[i] = hashKeys[i]
	}

	dataRowFieldsList, err := t.Backend.GetMultiple(t.getSchema(), genericKeys)
	if err != nil {
		return nil, err
	} else if len(dataRowFieldsList) > len(hashKeys) {
		return nil, errors.New("Datastore constraint not satisfied, more values than keys returned")
	}

	dataRowResults := make([]V, len(dataRowFieldsList))
	for i, dataRowFields := range dataRowFieldsList {
		err = t.getSchema().validateDataRowFields(dataRowFields)
		if err != nil {
			return nil, err
		}

		dataRowResults[i], err = t.DataRowFactory.CreateFromFields(dataRowFields)
		if err != nil {
			return nil, err
		}
	}

	return dataRowResults, nil
}

func (t *HashTable[V, H]) Add(hashKey H, data V, options WriteOptions) (*H, error) {
	vals, err := t.AddMultiple([]H{hashKey}, []V{data}, options)
	if err != nil {
		return nil, err
	}

	if len(vals) == 0 {
		return nil, nil
	} else {
		return &vals[0], nil
	}
}

func (t *HashTable[V, H]) AddMultiple(hashKeys []H, data []V, options WriteOptions) ([]H, error) {
	if len(hashKeys) != len(data) {
		return nil, errors.New("The number of HashKeys must match the number of data values")
	}

	genericKeys := make([]HashKey, len(hashKeys))
	genericData := make([]DataRow, len(data))
	for i := range hashKeys {
		genericKeys[i] = hashKeys[i]
		genericData[i] = data[i]
	}

	hashKeyFieldsList, err := t.Backend.AddMultiple(t.getSchema(), genericKeys, genericData, options)
	if err != nil {
		return nil, err
	} else if len(hashKeyFieldsList) > len(hashKeys) {
		return nil, errors.New("Datastore constraint not satisfied, more values than keys returned")
	}

	hashKeyResults := make([]H, len(hashKeyFieldsList))
	for i, hashKeyFields := range hashKeyFieldsList {
		err = t.getSchema().validateHashKeyFields(hashKeyFields)
		if err != nil {
			return nil, err
		}

		hashKeyResults[i], err = t.HashKeyFactory.CreateFromFields(hashKeyFields)
		if err != nil {
			return nil, err
		}
	}

	return hashKeyResults, nil
}

func (t *HashTable[V, H]) Update(hashKey H, data V, options UpdateOptions) error {
	return t.UpdateMultiple([]H{hashKey}, []V{data}, options)
}

func (t *HashTable[V, H]) UpdateMultiple(hashKeys []H, data []V, options UpdateOptions) error {
	if len(hashKeys) != len(data) {
		return errors.New("The number of HashKeys must match the number of data values")
	}

	genericKeys := make([]HashKey, len(hashKeys))
	genericData := make([]DataRow, len(data))
	for i := range hashKeys {
		genericKeys[i] = hashKeys[i]
		genericData[i] = data[i]
	}

	return t.Backend.UpdateMultiple(t.getSchema(), genericKeys, genericData, options)
}

func (t *HashTable[V, H]) Delete(hashKeys ...H) error {
	genericKeys := make([]HashKey, len(hashKeys))
	for i := range hashKeys {
		genericKeys[i] = hashKeys[i]
	}

	return t.Backend.DeleteMutiple(t.getSchema(), genericKeys)
}
