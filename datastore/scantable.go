package datastore

type ScanTable[V DataRow] struct {
	Backend        ScanTableBackend
	DataRowFactory DataRowFactory[V]
	Name           string
	schema         *ScanTableSchema
}

func (t *ScanTable[V]) getSchema() *ScanTableSchema {
	if t.schema == nil {
		t.schema = &ScanTableSchema{
			Name:                 t.Name,
			DataRowSchemaFactory: t.DataRowFactory,
		}
	}

	return t.schema
}

func (t *ScanTable[V]) ValidateSchema() error {
	return t.Backend.ValidateScanTableSchema(t.getSchema())
}

func (t *ScanTable[V]) Scan() (chan V, chan error) {
	dataRowFieldsChan, scanErrorChan := t.Backend.Scan(t.getSchema())

	dataRowChan := make(chan V)
	errorChan := make(chan error)
	go func(dataRowFieldsChan chan DataRowFields, scanErrorChan chan error, dataRowChan chan V, errorChan chan error) {
		select {
		case err := <-scanErrorChan:
			errorChan <- err
		case dataRowFields, more := <-dataRowFieldsChan:
			dataRow, err := t.DataRowFactory.CreateFromFields(dataRowFields)

			if err != nil {
				errorChan <- err
			} else {
				dataRowChan <- dataRow
			}

			if !more {
				close(dataRowChan)
			}
		}
		close(errorChan)
	}(dataRowFieldsChan, scanErrorChan, dataRowChan, errorChan)

	return dataRowChan, errorChan
}
