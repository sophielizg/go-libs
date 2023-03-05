package inmemory

import (
	"github.com/sophielizg/go-libs/datastore"
)

type InMemoryScanTableBackend struct{}

func (b *InMemoryScanTableBackend) ValidateSchema(schema *datastore.ScanTableSchema) error {
	// No special validation for this backend
	return nil
}

func (b *InMemoryScanTableBackend) CreateOrUpdateSchema(schema *datastore.ScanTableSchema) error {
	return createOrUpdateDbSchema(schema.Name)
}

func (b *InMemoryScanTableBackend) Scan(schema *datastore.ScanTableSchema, batchSize int) (chan *datastore.DataRowScanFields, chan error) {
	return scanDb(schema.Name, batchSize)
}
