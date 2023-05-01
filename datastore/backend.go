package datastore

import (
	"github.com/sophielizg/go-libs/datastore/mutator"
	"github.com/sophielizg/go-libs/datastore/queries"
)

type TableBackend[C Connection] interface {
	// Configuration
	SetSettings(settings *TableSettings)
	SetConnection(conn C)

	// Table
	Register() error
	Drop() error
}

type AppendTableBackendQueries interface {
	queries.ScanableBackend
	queries.AddableBackend
}

type AppendTableBackend[C Connection] interface {
	TableBackend[C]
	AppendTableBackendQueries
}

type HashTableBackendQueries interface {
	queries.ScanableBackend
	queries.CountableBackend
	queries.CRUDableBackend
}

type HashTableBackend[C Connection] interface {
	TableBackend[C]
	HashTableBackendQueries
}

type SortTableBackendQueries interface {
	queries.ScanableBackend
	queries.CountableBackend
	queries.CRUDableBackend
	queries.SortableBackend
}

type SortTableBackend[C Connection] interface {
	TableBackend[C]
	SortTableBackendQueries
}

type QueueBackendQueries interface {
	queries.CountableBackend
	queries.MessageReceiveableBackend
	SendMessage(messages []mutator.MappedFieldValues) error
}

type QueueBackend[C Connection] interface {
	TableBackend[C]
	QueueBackendQueries
}

type TopicBackendQueries interface {
	Publish(messages []mutator.MappedFieldValues) error
	Subscribe(subscriptionId string) (SubscriptionBackendQueries, error)
}

type SubscriptionBackendQueries interface {
	queries.MessageReceiveableBackend
	Unsubscribe() error
}

type TopicBackend[C Connection] interface {
	TableBackend[C]
	TopicBackendQueries
}
