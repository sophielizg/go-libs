package datastoremysql

import (
	"fmt"

	"github.com/sophielizg/go-libs/datastore"
)

type ScanTableBackend struct {
	conn Connection
}

func (b *ScanTableBackend) ValidateScanTableSchema(schema *datastore.ScanTableSchema) error {
	return nil
}

func (b *ScanTableBackend) CreateOrUpdateScanTableSchema(schema *datastore.ScanTableSchema) error {
	// TODO
}

func (b *ScanTableBackend) Scan(schema *datastore.ScanTableSchema) (chan datastore.DataRowFields, chan error) {
	resultsChan := make(chan datastore.DataRowFields)
	errorsChan := make(chan error)

	go func(schema *datastore.ScanTableSchema, resultsChan chan datastore.DataRowFields, errorsChan chan error) {
		query := fmt.Sprintf("SELECT * FROM %s", schema.Name)
		rows, err := b.conn.db.Queryx(query)
		if err != nil {
			errorsChan <- err
			close(resultsChan)
			close(errorsChan)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var dataRowFields datastore.DataRowFields
			err = rows.MapScan(dataRowFields)

			if err != nil {
				errorsChan <- err
			} else {
				resultsChan <- dataRowFields
			}
		}

		close(resultsChan)
		close(errorsChan)
	}(schema, resultsChan, errorsChan)

	return resultsChan, errorsChan
}
