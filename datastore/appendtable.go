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

func (t *AppendTable[V]) getSupportedWriteOptions() SupportedOptions[WriteOption] {
	supported := t.Backend.SupportedWriteOptions()
	if supported == nil {
		supported = DefaultSupportedWriteOptions
	}

	return supported
}

func (t *AppendTable[V]) ValidateSchema() error {
	return t.Backend.ValidateAppendTableSchema(t.getSchema())
}

func (t *AppendTable[V]) AppendMultiple(data []V, options Options[WriteOption]) error {
	err := t.getSchema().validateWriteOptions(options)
	if err != nil {
		return err
	}

	genericData := make([]DataRow, len(data))
	for i := range data {
		genericData[i] = data[i]
	}

	return t.Backend.AppendMultiple(t.getSchema(), genericData, options)
}
