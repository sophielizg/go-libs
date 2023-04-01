package purchase

import "github.com/sophielizg/go-libs/datastore"

type PurchaseTable = datastore.SortTable[PurchaseDataRow, *PurchaseDataRow, PurchaseHashKey, *PurchaseHashKey, PurchaseSortKey, *PurchaseSortKey, PurchaseSortKeyComparator, *PurchaseSortKeyComparator]

func NewPurchaseTable() *PurchaseTable {
	return &PurchaseTable{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Purchase"),
			datastore.WithDataRowSettings(&PurchaseDataRowSettings),
			datastore.WithHashKeySettings(&ProductHashKeySettings),
			datastore.WithSortKeySettings(&PurchaseSortKeySettings),
		),
	}
}
