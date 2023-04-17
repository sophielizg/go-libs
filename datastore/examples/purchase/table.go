package purchase

import "github.com/sophielizg/go-libs/datastore"

type PurchaseTable = datastore.SortTable[DataRow, *DataRow, HashKey, *HashKey, SortKey, *SortKey, SortKeyComparator, *SortKeyComparator]

type TableScan = datastore.SortTableScan[DataRow, *DataRow, HashKey, *HashKey, SortKey, *SortKey]

func NewTable() *PurchaseTable {
	return &PurchaseTable{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Purchase"),
			datastore.WithDataRowSettings(&DataRowSettings),
			datastore.WithHashKeySettings(&HashKeySettings),
			datastore.WithSortKeySettings(&SortKeySettings),
		),
	}
}
