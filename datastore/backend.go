package datastore

import "github.com/sophielizg/go-libs/datastore/mutator"

type TableBackend[C Connection] interface {
	// Configuration
	SetSettings(settings *TableSettings)
	SetConnection(conn C)

	// Table
	Register() error
	Drop() error
}

type AppendTableBackendOps interface {
	// Data operations
	Scan(batchSize int) (chan *ScanFields, chan error)
	AddMultiple(data []mutator.MappedFieldValues) error
}

type AppendTableBackend[C Connection] interface {
	TableBackend[C]
	AppendTableBackendOps
}

type HashTableBackendOps interface {
	// Data operations
	Scan(batchSize int) (chan *ScanFields, chan error)

	GetMultiple(hashKeys []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error)
	AddMultiple(hashKeys []mutator.MappedFieldValues, data []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error)
	UpdateMultiple(hashKeys []mutator.MappedFieldValues, data []mutator.MappedFieldValues) error
	DeleteMultiple(hashKeys []mutator.MappedFieldValues) error
}

type HashTableBackend[C Connection] interface {
	TableBackend[C]
	HashTableBackendOps
}

type SortTableBackendOps interface {
	// Data operations
	Scan(batchSize int) (chan *ScanFields, chan error)

	GetMultiple(hashKeys []mutator.MappedFieldValues, sortKeys []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, error)
	AddMultiple(hashKeys []mutator.MappedFieldValues, sortKeys []mutator.MappedFieldValues, data []mutator.MappedFieldValues) ([]mutator.MappedFieldValues, []mutator.MappedFieldValues, error)
	UpdateMultiple(hashKeys []mutator.MappedFieldValues, sortKeys []mutator.MappedFieldValues, data []mutator.MappedFieldValues) error
	DeleteMultiple(hashKeys []mutator.MappedFieldValues, sortKeys []mutator.MappedFieldValues) error

	GetWithSortComparator(hashKey mutator.MappedFieldValues, comparator mutator.MappedFieldValues) ([]mutator.MappedFieldValues, []mutator.MappedFieldValues, error)
	UpdateWithSortComparator(hashKey mutator.MappedFieldValues, comparator mutator.MappedFieldValues, data mutator.MappedFieldValues) error
	DeleteWithSortComparator(hashKey mutator.MappedFieldValues, comparator mutator.MappedFieldValues) error
}

type SortTableBackend[C Connection] interface {
	TableBackend[C]
	SortTableBackendOps
}

type QueueBackendOps interface {
	// Data operations
	Size() (int, error)
	Push(messages []mutator.MappedFieldValues) error
	Pop() (string, mutator.MappedFieldValues, error)
	AckSuccess(messageId string) error
	AckFailure(messageId string) error
}

type QueueBackend[C Connection] interface {
	TableBackend[C]
	QueueBackendOps
}

type TopicBackendOps interface {
	// Topic operations
	Publish(messages []mutator.MappedFieldValues) error
	Subscribe(subscriptionId string) error
	Unsubscribe(subscriptionId string) error

	// Subscription operations
	HasMessage(subscriptionId string) (bool, error)
	RecieveMessage(subscriptionId string) (string, mutator.MappedFieldValues, error)
	AckSuccess(subscriptionId string, messageId string) error
	AckFailure(subscriptionId string, messageId string) error
}

type TopicBackend[C Connection] interface {
	TableBackend[C]
	TopicBackendOps
}
