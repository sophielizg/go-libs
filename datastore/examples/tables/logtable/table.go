package logtable

import "github.com/sophielizg/go-libs/datastore"

func LogReaderTable(backend datastore.ScanTableBackend) datastore.ScanTable[*LogDataRow] {
	return datastore.ScanTable[*LogDataRow]{
		Name:           "Log",
		DataRowFactory: &LogDataRowFactory{},
		Backend:        backend,
	}
}

func LogWriterTable(backend datastore.AppendTableBackend) datastore.AppendTable[*LogDataRow] {
	return datastore.AppendTable[*LogDataRow]{
		Name:           "Log",
		DataRowFactory: &LogDataRowFactory{},
		Backend:        backend,
	}
}
