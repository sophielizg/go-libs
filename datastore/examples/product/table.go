package product

import "github.com/sophielizg/go-libs/datastore"

type Table = datastore.HashTable[Key, *Key, Entry, *Entry]

func NewTable() *Table {
	return &Table{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Product"),
			datastore.WithDataSettings(DataSettings),
			datastore.WithKeySettings(KeySettings),
		),
	}
}
