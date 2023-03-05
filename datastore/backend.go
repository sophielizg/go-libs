package datastore

type AppendTableBackend interface {
	// Schema
	ValidateSchema(schema *AppendTableSchema) error
	CreateOrUpdateSchema(schema *AppendTableSchema) error

	// Data operations
	Scan(schema *AppendTableSchema, batchSize int) (chan *DataRowScanFields, chan error)
	AppendMultiple(schema *AppendTableSchema, data []DataRow) error
}

type HashTableBackend interface {
	// Configuration
	SupportedFieldOptions() SupportedOptions

	// Schema
	ValidateSchema(schema *HashTableSchema) error
	CreateOrUpdateSchema(schema *HashTableSchema) error

	// Data operations
	Scan(schema *HashTableSchema, batchSize int) (chan *HashTableScanFields, chan error)

	GetMultiple(schema *HashTableSchema, hashKeys []HashKey) ([]DataRowFields, error)
	AddMultiple(schema *HashTableSchema, hashKeys []HashKey, data []DataRow) ([]DataRowFields, error)
	UpdateMultiple(schema *HashTableSchema, hashKeys []HashKey, data []DataRow) error
	DeleteMultiple(schema *HashTableSchema, hashKeys []HashKey) error
}

type SortTableBackend interface {
	// Configuration
	SupportedFieldOptions() SupportedOptions

	// Schema
	ValidateSchema(schema *SortTableSchema) error
	CreateOrUpdateSchema(schema *SortTableSchema) error

	// Data operations
	Scan(schema *SortTableSchema, batchSize int) (chan *SortTableScanFields, chan error)

	GetMultiple(schema *SortTableSchema, hashKeys []HashKey, sortKeys []SortKey) ([]DataRowFields, error)
	AddMultiple(schema *SortTableSchema, hashKeys []HashKey, sortKeys []SortKey, data []DataRow) ([]DataRowFields, []DataRowFields, error)
	UpdateMultiple(schema *SortTableSchema, hashKeys []HashKey, sortKeys []SortKey, data []DataRow) error
	DeleteMultiple(schema *SortTableSchema, hashKeys []HashKey, sortKeys []SortKey) error

	GetWithSortKey(schema *SortTableSchema, hashKey HashKey, sortKey SortKey) ([]DataRowFields, []DataRowFields, error)
	UpdateWithSortKey(schema *SortTableSchema, hashKey HashKey, sortKey SortKey, data DataRow) error
	DeleteWithSortKey(schema *SortTableSchema, hashKey HashKey, sortKey SortKey) error
}
