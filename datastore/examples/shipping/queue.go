package shipping

import "github.com/sophielizg/go-libs/datastore"

type Queue = datastore.Queue[Message, *Message]

func NewPendingShipmentQueue() *Queue {
	return &Queue{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("PendingShipment"),
			datastore.WithDataSettings(DataSettings),
		),
	}
}
