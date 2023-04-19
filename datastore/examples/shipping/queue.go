package shipping

import "github.com/sophielizg/go-libs/datastore"

type Queue = datastore.Queue[DataRow, *DataRow]

func NewPendingShipmentQueue() *Queue {
	return &Queue{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("PendingShipment"),
			datastore.WithDataRowSettings(&DataRowSettings),
		),
	}
}
