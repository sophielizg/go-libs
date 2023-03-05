package inmemory

import "github.com/sophielizg/go-libs/datastore"

type InMemoryAppendTableBackend struct{}

func (b *InMemoryAppendTableBackend) ValidateSchema(schema *datastore.AppendTableSchema) error {
	// No special validation for this backend
	return nil
}

func (b *InMemoryAppendTableBackend) CreateOrUpdateSchema(schema *datastore.AppendTableSchema) error {
	return createOrUpdateDbSchema(schema.Name)
}

func (b *InMemoryAppendTableBackend) Scan(schema *datastore.AppendTableSchema, batchSize int) (chan *datastore.DataRowScanFields, chan error) {
	return scanDb(schema.Name, batchSize)
}

func (b *InMemoryAppendTableBackend) AppendMultiple(schema *datastore.AppendTableSchema, data []datastore.DataRow) error {
	for _, dataRow := range data {
		hashKey := datastore.DataRowFields{
			"Id": generateStringKey(64),
		}
		err := insertWithHashKey(schema.Name, hashKey, dataRow.GetFields())
		if err != nil {
			return err
		}
	}

	return nil
}
