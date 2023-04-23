package logtable

import "github.com/sophielizg/go-libs/datastore"

type LogTable = datastore.AppendTable[DataRow, *DataRow]

func New(tableName string) *LogTable {
	return &LogTable{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName(tableName),
			datastore.WithDataRowSettings(&DataRowSettings),
		),
	}
}
