package shipping

import "github.com/sophielizg/go-libs/datastore"

type Topic = datastore.Topic[DataRow, *DataRow]

func NewShippedTopic() *Topic {
	return &Topic{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Shipped"),
			datastore.WithDataRowSettings(&DataRowSettings),
		),
	}
}
