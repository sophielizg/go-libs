package inmemory

import "github.com/sophielizg/go-libs/datastore"

type InMemoryAppendTableBackend struct {
	Conn InMemoryDatastoreConnection
}

func (b *InMemoryAppendTableBackend) ValidateSchema(schema *datastore.AppendTableSchema) error {
	// No special validation for this backend
	return nil
}

func (b *InMemoryAppendTableBackend) CreateOrUpdateSchema(schema *datastore.AppendTableSchema) error {
	return b.Conn.CreateOrUpdateSchema(schema.Name)
}

func (b *InMemoryAppendTableBackend) Scan(schema *datastore.AppendTableSchema, batchSize int) (chan *datastore.DataRowScanFields, chan error) {
	outResChan := make(chan *datastore.DataRowScanFields, 1)
	inResChan, errorChan := b.Conn.Scan(schema.Name, batchSize)

	go func() {
		defer close(outResChan)

		for {
			res, more := <-inResChan
			if !more {
				inResChan = nil
				break
			}

			outResChan <- &res.DataRowScanFields
		}
	}()

	return outResChan, errorChan
}

func (b *InMemoryAppendTableBackend) AppendMultiple(schema *datastore.AppendTableSchema, data []datastore.DataRow) error {
	for _, dataRow := range data {
		hashKey := datastore.DataRowFields{
			"Id": generateStringKey(64),
		}
		err := b.Conn.Add(schema.Name, hashKey, dataRow.GetFields())
		if err != nil {
			return err
		}
	}

	return nil
}
