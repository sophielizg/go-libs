package inmemory

import (
	"errors"

	"github.com/sophielizg/go-libs/datastore"
)

type HashTable = map[string]datastore.DataRowFields

type InMemoryHashTableBackend struct {
	table map[string]HashTable
}

func (b *InMemoryHashTableBackend) SupportedFieldOptions() datastore.SupportedOptions {
	// No special field support for this backend
	return nil
}

func (b *InMemoryHashTableBackend) ValidateSchema(schema *datastore.HashTableSchema) error {
	// No special validation for this backend
	return nil
}

func (b *InMemoryHashTableBackend) CreateOrUpdateSchema(schema *datastore.HashTableSchema) error {
	if b.table == nil {
		b.table = map[string]HashTable{}
	}

	if b.table[schema.Name] == nil {
		b.table[schema.Name] = make(HashTable, 0)
	}

	return nil
}

func (b *InMemoryHashTableBackend) getTable(schema *datastore.HashTableSchema) (HashTable, error) {
	if b.table[schema.Name] == nil {
		return nil, errors.New("No table exists with given schema name")
	}
	return b.table[schema.Name], nil
}

func (b *InMemoryHashTableBackend) Scan(schema *datastore.HashTableSchema, batchSize int) (chan *datastore.HashTableScanFields, chan error) {
	outChan := make(chan *datastore.HashTableScanFields, batchSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(outChan)
		defer close(errorChan)

		table, err := b.getTable(schema)
		if err != nil {
			errorChan <- err
			return
		}

		for key, row := range table {
			hashKey, err := unstringifyKey(key)
			if err != nil {
				errorChan <- err
				continue
			}

			outChan <- &datastore.HashTableScanFields{
				AppendTableScanFields: datastore.AppendTableScanFields{
					DataRow: row,
				},
				HashKey: hashKey,
			}
		}
	}()

	return outChan, errorChan
}

func (b *InMemoryHashTableBackend) GetMultiple(schema *datastore.HashTableSchema, hashKeys []datastore.HashKey) ([]datastore.DataRowFields, error) {
	res := make([]datastore.DataRowFields, len(hashKeys))

	table, err := b.getTable(schema)
	if err != nil {
		return nil, err
	}

	for i, hashKey := range hashKeys {
		key, err := stringifyKey(hashKey.GetFields())
		if err != nil {
			return nil, err
		}

		res[i] = table[key]
	}

	return res, nil
}

func (b *InMemoryHashTableBackend) AddMultiple(schema *datastore.HashTableSchema, hashKeys []datastore.HashKey, data []datastore.DataRow) ([]datastore.DataRowFields, error) {
	res := make([]datastore.DataRowFields, len(hashKeys))

	table, err := b.getTable(schema)
	if err != nil {
		return nil, err
	}

	for i, hashKey := range hashKeys {
		hashKeyFields, err := generateUniqueKey(table, hashKey, schema.HashKeySchemaFactory)
		if err != nil {
			return nil, err
		}

		key, err := stringifyKey(hashKeyFields)
		if err != nil {
			return nil, err
		}

		res[i] = hashKeyFields
		table[key] = data[i].GetFields()
	}

	return res, nil
}

func (b *InMemoryHashTableBackend) UpdateMultiple(schema *datastore.HashTableSchema, hashKeys []datastore.HashKey, data []datastore.DataRow) error {
	table, err := b.getTable(schema)
	if err != nil {
		return err
	}

	for i, hashKey := range hashKeys {
		key, err := stringifyKey(hashKey.GetFields())
		if err != nil {
			return err
		}

		table[key] = data[i].GetFields()
	}

	return nil
}

func (b *InMemoryHashTableBackend) DeleteMultiple(schema *datastore.HashTableSchema, hashKeys []datastore.HashKey) error {
	table, err := b.getTable(schema)
	if err != nil {
		return err
	}

	for _, hashKey := range hashKeys {
		key, err := stringifyKey(hashKey.GetFields())
		if err != nil {
			return err
		}

		delete(table, key)
	}

	return nil
}
