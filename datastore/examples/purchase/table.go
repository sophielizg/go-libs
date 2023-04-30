package purchase

import "github.com/sophielizg/go-libs/datastore"

type PurchaseTable = datastore.SortTable[Key, *Key, Entry, *Entry, SortComparator, *SortComparator]

func NewTable() *PurchaseTable {
	return &PurchaseTable{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Purchase"),
			datastore.WithDataSettings(DataSettings),
			datastore.WithKeySettings(KeySettings),
			datastore.WithSortFieldNames(SortFieldNames),
		),
	}
}
