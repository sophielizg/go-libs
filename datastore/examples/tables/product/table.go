package product

import "github.com/sophielizg/go-libs/datastore"

func ProductTable(backend datastore.HashTableBackend) datastore.HashTable[*ProductDataRow, *ProductHashKey] {
	return datastore.HashTable[*ProductDataRow, *ProductHashKey]{
		Name:           "Product",
		DataRowFactory: &productDataRowFactory{},
		HashKeyFactory: &productHashKeyFactory{},
		Backend:        backend,
	}
}
