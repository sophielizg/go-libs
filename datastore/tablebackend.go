package datastore

import "github.com/sophielizg/go-libs/datastore/fields"

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
	AddMultiple(data []fields.MappedFieldValues) error
}

type AppendTableBackend[C Connection] interface {
	TableBackend[C]
	AppendTableBackendOps
}

type HashTableBackendOps interface {
	// Data operations
	Scan(batchSize int) (chan *ScanFields, chan error)

	GetMultiple(hashKeys []fields.MappedFieldValues) ([]fields.MappedFieldValues, error)
	AddMultiple(hashKeys []fields.MappedFieldValues, data []fields.MappedFieldValues) ([]fields.MappedFieldValues, error)
	UpdateMultiple(hashKeys []fields.MappedFieldValues, data []fields.MappedFieldValues) error
	DeleteMultiple(hashKeys []fields.MappedFieldValues) error
}

type HashTableBackend[C Connection] interface {
	TableBackend[C]
	HashTableBackendOps
}

type SortTableBackendOps interface {
	// Data operations
	Scan(batchSize int) (chan *ScanFields, chan error)

	GetMultiple(hashKeys []fields.MappedFieldValues, sortKeys []fields.MappedFieldValues) ([]fields.MappedFieldValues, error)
	AddMultiple(hashKeys []fields.MappedFieldValues, sortKeys []fields.MappedFieldValues, data []fields.MappedFieldValues) ([]fields.MappedFieldValues, []fields.MappedFieldValues, error)
	UpdateMultiple(hashKeys []fields.MappedFieldValues, sortKeys []fields.MappedFieldValues, data []fields.MappedFieldValues) error
	DeleteMultiple(hashKeys []fields.MappedFieldValues, sortKeys []fields.MappedFieldValues) error

	GetWithSortKey(hashKey fields.MappedFieldValues, sortKey fields.MappedFieldValues) ([]fields.MappedFieldValues, []fields.MappedFieldValues, error)
	UpdateWithSortKey(hashKey fields.MappedFieldValues, sortKey fields.MappedFieldValues, data fields.MappedFieldValues) error
	DeleteWithSortKey(hashKey fields.MappedFieldValues, sortKey fields.MappedFieldValues) error
}

type SortTableBackend[C Connection] interface {
	TableBackend[C]
	SortTableBackendOps
}
