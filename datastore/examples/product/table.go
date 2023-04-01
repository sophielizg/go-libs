package product

import "github.com/sophielizg/go-libs/datastore"

func NewProductTable() *datastore.HashTable[ProductDataRow, *ProductDataRow, ProductHashKey, *ProductHashKey] {
	return &datastore.HashTable[ProductDataRow, *ProductDataRow, ProductHashKey, *ProductHashKey]{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Product"),
		),
	}
}
