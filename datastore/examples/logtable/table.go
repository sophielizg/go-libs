package logtable

import "github.com/sophielizg/go-libs/datastore"

func NewLogTable() *datastore.AppendTable[LogDataRow, *LogDataRow] {
	return &datastore.AppendTable[LogDataRow, *LogDataRow]{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Log"),
		),
	}
}
