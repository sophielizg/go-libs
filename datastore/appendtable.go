package datastore

// A table to which data can be appended
type AppendTable[V any, PV DataRow[V]] struct {
	Backend        AppendTableBackendOps
	Settings       *TableSettings
	DataRowFactory DataRowFactory[V, PV]
}

func (t *AppendTable[V, PV]) Init() {
	t.Settings.ApplyOption(WithDataRow[V, PV]())
}

func (t *AppendTable[V, PV]) GetSettings() *TableSettings {
	return t.Settings
}

func (t *AppendTable[V, PV]) SetBackend(tableBackend AppendTableBackendOps) {
	t.Backend = tableBackend
}

// Scans the entire table, holding batchSize data rows in memory at a time
func (t *AppendTable[V, PV]) Scan(batchSize int) (chan PV, chan error) {
	fieldsChan, errChan := t.Backend.Scan(batchSize)
	return scan(fieldsChan, errChan, func(fields *ScanFields) (PV, error) {
		return t.DataRowFactory.CreateFromFields(fields.DataRow)
	})
}

// Adds data to the table
func (t *AppendTable[V, PV]) Add(data ...PV) error {
	return t.Backend.AddMultiple(t.DataRowFactory.CreateFieldValuesList(data))
}
