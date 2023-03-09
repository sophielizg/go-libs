package logtable

import "github.com/sophielizg/go-libs/datastore"

func CreateLogTable(backend datastore.AppendTableBackend) *LogTable {
	return &LogTable{
		table: datastore.AppendTable[*LogDataRow]{
			Name:           "Log",
			DataRowFactory: &logDataRowFactory{},
			Backend:        backend,
		},
	}
}

type LogTable struct {
	table datastore.AppendTable[*LogDataRow]
}

func (t *LogTable) Append(data ...*LogDataRow) error {
	return t.table.Append(data...)
}
