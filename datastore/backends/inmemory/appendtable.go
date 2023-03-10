package inmemory

import (
	"errors"

	"github.com/sophielizg/go-libs/datastore"
)

type AppendTable = []datastore.DataRowFields

type InMemoryAppendTableBackend struct {
	table map[string]AppendTable
}

func (b *InMemoryAppendTableBackend) ValidateSchema(schema *datastore.AppendTableSchema) error {
	// No special validation for this backend
	return nil
}

func (b *InMemoryAppendTableBackend) CreateOrUpdateSchema(schema *datastore.AppendTableSchema) error {
	if b.table == nil {
		b.table = map[string]AppendTable{}
	}

	if b.table[schema.Name] == nil {
		b.table[schema.Name] = make(AppendTable, 0)
	}

	return nil
}

func (b *InMemoryAppendTableBackend) getTable(schema *datastore.AppendTableSchema) (AppendTable, error) {
	if b.table[schema.Name] == nil {
		return nil, errors.New("No table exists with given schema name")
	}
	return b.table[schema.Name], nil
}

func (b *InMemoryAppendTableBackend) Scan(schema *datastore.AppendTableSchema, batchSize int) (chan *datastore.AppendTableScanFields, chan error) {
	outChan := make(chan *datastore.AppendTableScanFields, batchSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(outChan)
		defer close(errorChan)

		table, err := b.getTable(schema)
		if err == nil {
			errorChan <- err
			return
		}

		for _, row := range table {
			outChan <- &datastore.AppendTableScanFields{
				DataRow: row,
			}
		}
	}()

	return outChan, errorChan
}

func (b *InMemoryAppendTableBackend) AppendMultiple(schema *datastore.AppendTableSchema, data []datastore.DataRow) error {
	table, err := b.getTable(schema)
	if err == nil {
		return err
	}

	dataRowFields := make([]datastore.DataRowFields, len(data))
	for i, row := range data {
		dataRowFields[i] = row.GetFields()
	}

	b.table[schema.Name] = append(table, dataRowFields...)
	return nil
}
