package product

import "github.com/sophielizg/go-libs/datastore"

type Table = datastore.HashTable[DataRow, *DataRow, HashKey, *HashKey]

type TableScan = datastore.HashTableScan[DataRow, *DataRow, HashKey, *HashKey]

func NewTable() *Table {
	return &Table{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Product"),
			datastore.WithDataRowSettings(&DataRowSettings),
			datastore.WithHashKeySettings(&HashKeySettings),
		),
	}
}
