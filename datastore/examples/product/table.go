package product

import "github.com/sophielizg/go-libs/datastore"

type ProductTable = datastore.HashTable[ProductDataRow, *ProductDataRow, ProductHashKey, *ProductHashKey]

func NewProductTable() *ProductTable {
	return &ProductTable{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Product"),
			datastore.WithDataRowSettings(&ProductDataRowSettings),
			datastore.WithHashKeySettings(&ProductHashKeySettings),
		),
	}
}
