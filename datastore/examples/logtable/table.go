package logtable

import "github.com/sophielizg/go-libs/datastore"

type Table = datastore.AppendTable[Entry, *Entry]

func New() *Table {
	return &Table{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Log"),
			datastore.WithDataSettings(DataSettings),
		),
	}
}
