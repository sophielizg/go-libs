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
			ScanTableSchema: ScanTableSchema{
				Name:                 t.Name,
				DataRowSchemaFactory: t.DataRowFactory,
			},
		}
	}

	return t.schema
}

func (t *AppendTable[V]) ValidateSchema() error {
	return t.Backend.ValidateAppendTableSchema(t.getSchema())
}

func (t *AppendTable[V]) AppendMultiple(data []V) error {
	genericData := make([]DataRow, len(data))
	for i := range data {
		genericData[i] = data[i]
	}

	return t.Backend.AppendMultiple(t.getSchema(), genericData)
}
