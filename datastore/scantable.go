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
	return t.Backend.ValidateSchema(t.getSchema())
}

func (t *ScanTable[V]) CreateOrUpdateSchema() error {
	return t.Backend.CreateOrUpdateSchema(t.getSchema())
}

func (t *ScanTable[V]) Scan(batchSize int) (chan DataRowScan[V], chan error) {
	scanDataRowChan, scanErrorChan := t.Backend.Scan(t.getSchema(), batchSize)
	return scan(
		batchSize,
		scanDataRowChan,
		scanErrorChan,
		func(scanDataRow *DataRowScanFields) (DataRowScan[V], error) {
			var err error
			res := DataRowScan[V]{}
			res.DataRow, err = t.DataRowFactory.CreateFromFields(scanDataRow.DataRow)
			return res, err
		},
	)
}
