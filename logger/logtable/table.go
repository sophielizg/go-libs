package logtable

import "github.com/sophielizg/go-libs/datastore"

type LogTable = datastore.AppendTable[*LogDataRow]

func CreateLogTable(backend datastore.AppendTableBackend) LogTable {
	return datastore.AppendTable[*LogDataRow]{
		Name:           "Log",
		DataRowFactory: &logDataRowFactory{},
		Backend:        backend,
	}
}
