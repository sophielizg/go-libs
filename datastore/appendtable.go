package datastore

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

func (t *AppendTable[V]) ValidateSchema() error {
	return t.Backend.ValidateSchema(t.getSchema())
}

func (t *AppendTable[V]) CreateOrUpdateSchema() error {
	return t.Backend.CreateOrUpdateSchema(t.getSchema())
}

func (t *AppendTable[V]) Scan(batchSize int) (chan DataRowScan[V], chan error) {
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

func (t *AppendTable[V]) Append(data ...V) error {
	genericData := convertDataRowToInterface(data...)
	return t.Backend.AppendMultiple(t.getSchema(), genericData)
}
