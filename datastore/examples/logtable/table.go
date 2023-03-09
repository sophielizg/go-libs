package logtable

import "github.com/sophielizg/go-libs/datastore"

func LogTable(backend datastore.AppendTableBackend) datastore.AppendTable[*LogDataRow] {
	return datastore.AppendTable[*LogDataRow]{
		Name:           "Log",
		DataRowFactory: &LogDataRowFactory{},
		Backend:        backend,
	}
}
