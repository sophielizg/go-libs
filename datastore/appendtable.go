package datastore

// A table to which data can be appended
type AppendTable[V DataRow] struct {
	Backend        AppendTableBackend
	DataRowFactory DataRowFactory[V]
	Name           string
	schema         *AppendTableSchema
}

func (t *AppendTable[V]) getSchema() *AppendTableSchema {
	if t.schema == nil {
		t.schema = &AppendTableSchema{
			BaseTableSchema: BaseTableSchema{
				Name:                 t.Name,
				DataRowSchemaFactory: t.DataRowFactory,
			},
		}
	}

	return t.schema
}

// Validates the schema of the table
func (t *AppendTable[V]) ValidateSchema() error {
	return t.Backend.ValidateSchema(t.getSchema())
}

// Creates or updates the schema of the table
func (t *AppendTable[V]) CreateOrUpdateSchema() error {
	return t.Backend.CreateOrUpdateSchema(t.getSchema())
}

// Scans the entire table, holding batchSize data rows in memory at a time
func (t *AppendTable[V]) Scan(batchSize int) (chan DataRowScan[V], chan error) {
	scanDataRowChan, scanErrorChan := t.Backend.Scan(t.getSchema(), batchSize)
	return scan(
		batchSize,
		scanDataRowChan,
		scanErrorChan,
		func(scanDataRow *AppendTableScanFields) (DataRowScan[V], error) {
			res := DataRowScan[V]{}
			rows, err := convertDataRowFieldsToInterface(
				[]DataRowFields{scanDataRow.DataRow},
				t.getSchema().validateDataRowFields,
				t.DataRowFactory,
			)

			if err != nil {
				return res, err
			}

			res.DataRow = rows[0]
			return res, err
		},
	)
}

// Appends data to the table
func (t *AppendTable[V]) Append(data ...V) error {
	genericData := convertDataRowToInterface(data...)
	return t.Backend.AppendMultiple(t.getSchema(), genericData)
}

// Deletes all the data from the table, holding batchSize data rows in memory at a time
func (t *AppendTable[V]) DeleteAll(batchSize int) error {
	return t.Backend.DeleteAll(t.getSchema(), batchSize)
}

// Transfers data from the current table to another table of the same type
func (t *AppendTable[V]) TransferTo(newTable *AppendTable[V], batchSize int) error {
	dataChan, errorChan := t.Scan(batchSize)

	for {
		buf := make([]V, 0, batchSize)

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

				buf = append(buf, data.DataRow)
			}

			if len(buf) == batchSize {
				break
			}
		}

		if len(buf) > 0 {
			err := newTable.Append(buf...)
			if err != nil {
				return err
			}
		}

		if dataChan == nil && errorChan == nil {
			return nil
		}
	}
}
