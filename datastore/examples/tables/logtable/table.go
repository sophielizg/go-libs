package logtable

import "github.com/sophielizg/go-libs/datastore"

func LogWriterTable(backend datastore.AppendTableBackend) datastore.AppendTable[*LogDataRow] {
	return datastore.AppendTable[*LogDataRow]{
		Name:           "Log",
		DataRowFactory: &LogDataRowFactory{},
		Backend:        backend,
	}
}
