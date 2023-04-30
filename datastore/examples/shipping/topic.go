package shipping

import "github.com/sophielizg/go-libs/datastore"

type Topic = datastore.Topic[Message, *Message]

func NewShippedTopic() *Topic {
	return &Topic{
		Settings: datastore.NewTableSettings(
			datastore.WithTableName("Shipped"),
			datastore.WithDataSettings(DataSettings),
		),
	}
}
