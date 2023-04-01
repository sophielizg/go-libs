package logtable

import "github.com/sophielizg/go-libs/datastore"

type LogTable = datastore.AppendTable[LogDataRow, *LogDataRow]

func NewLogTable() *LogTable {
	return &LogTable{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Log"),
			datastore.WithDataRowSettings(&LogDataRowSettings),
		),
	}
}
