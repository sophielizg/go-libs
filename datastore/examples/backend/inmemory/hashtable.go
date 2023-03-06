package inmemory

import "github.com/sophielizg/go-libs/datastore"

type InMemoryHashTableBackend struct {
	Conn InMemoryDatastoreConnection
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
	return b.Conn.CreateOrUpdateSchema(schema.Name)
}

func (b *InMemoryHashTableBackend) Scan(schema *datastore.HashTableSchema, batchSize int) (chan *datastore.HashTableScanFields, chan error) {
	return b.Conn.Scan(schema.Name, batchSize)
}

func (b *InMemoryHashTableBackend) GetMultiple(schema *datastore.HashTableSchema, hashKeys []datastore.HashKey) ([]datastore.DataRowFields, error) {
	res := make([]datastore.DataRowFields, len(hashKeys))
	for i, hashKey := range hashKeys {
		dataRow, err := b.Conn.Get(schema.Name, hashKey.GetFields())
		if err != nil {
			return nil, err
		}

		if len(dataRow) == 0 {
			res[i] = nil
		} else {
			// This table is a hash table, so assume that only one row exists per hash key
			res[i] = dataRow[0]
		}
	}

	return res, nil
}

func (b *InMemoryHashTableBackend) AddMultiple(schema *datastore.HashTableSchema, hashKeys []datastore.HashKey, data []datastore.DataRow) ([]datastore.DataRowFields, error) {
	res := make([]datastore.DataRowFields, len(hashKeys))
	for i, hashKey := range hashKeys {
		hashKeyFields, err := applyKeyOptions(hashKey.GetFields(), schema.GetFieldTypes(), schema.HashKeySchemaFactory.GetFieldOptions())
		if err != nil {
			return nil, err
		}
		res[i] = hashKeyFields

		err = b.Conn.Add(schema.Name, hashKeyFields, data[i].GetFields())
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (b *InMemoryHashTableBackend) UpdateMultiple(schema *datastore.HashTableSchema, hashKeys []datastore.HashKey, data []datastore.DataRow) error {
	for i, hashKey := range hashKeys {
		err := b.Conn.Update(schema.Name, hashKey.GetFields(), data[i].GetFields())
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *InMemoryHashTableBackend) DeleteMultiple(schema *datastore.HashTableSchema, hashKeys []datastore.HashKey) error {
	for _, hashKey := range hashKeys {
		err := b.Conn.Delete(schema.Name, hashKey.GetFields())
		if err != nil {
			return err
		}
	}

	return nil
}
