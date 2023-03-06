package purchase

import "github.com/sophielizg/go-libs/datastore"

func PurchaseTable(backend datastore.SortTableBackend) datastore.SortTable[*PurchaseDataRow, *PurchaseHashKey, *PurchaseSortKey] {
	return datastore.SortTable[*PurchaseDataRow, *PurchaseHashKey, *PurchaseSortKey]{
		Name:           "Purchase",
		DataRowFactory: &PurchaseDataRowFactory{},
		HashKeyFactory: &PurchaseHashKeyFactory{},
		SortKeyFactory: &PurchaseSortKeyFactory{},
		Backend:        backend,
	}
}
