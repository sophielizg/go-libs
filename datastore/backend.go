package datastore

type ScanTableBackend interface {
	// Schema
	ValidateScanTableSchema(schema *ScanTableSchema) error
	CreateOrUpdateScanTableSchema(schema *ScanTableSchema) error

	// Data operations
	Scan(schema *ScanTableSchema) (chan DataRowFields, chan error)
}

type AppendTableBackend interface {
	ScanTableBackend

	// Schema
	ValidateAppendTableSchema(schema *AppendTableSchema) error
	CreateOrUpdateAppendTableSchema(schema *AppendTableSchema) error

	// Data operations
	AppendMultiple(schema *AppendTableSchema, data []DataRow) error
}

type HashTableBackend interface {
	ScanTableBackend

	// Configuration
	SupportedFieldOptions() SupportedOptions[FieldOption]

	// Schema
	ValidateHashTableSchema(schema *HashTableSchema) error
	CreateOrUpdateHashTableSchema(schema *HashTableSchema) error

	// Data operations
	GetMultiple(schema *HashTableSchema, hashKeys []HashKey) ([]DataRowFields, error)
	AddMultiple(schema *HashTableSchema, hashKeys []HashKey, data []DataRow) ([]DataRowFields, error)
	UpdateMultiple(schema *HashTableSchema, hashKeys []HashKey, data []DataRow) error
	DeleteMultiple(schema *HashTableSchema, hashKeys []HashKey) error
}

type SortTableBackend interface {
	ScanTableBackend

	// Configuration
	SupportedFieldOptions() SupportedOptions[FieldOption]

	// Schema
	ValidateSortTableSchema(schema *SortTableSchema) error
	CreateOrUpdateSortTableSchema(schema *SortTableSchema) error

	// Data operations
	GetMultiple(schema *SortTableSchema, hashKeys []HashKey, sortKeys []SortKey) ([]DataRowFields, error)
	AddMultiple(schema *SortTableSchema, hashKeys []HashKey, sortKeys []SortKey, data []DataRow) ([]DataRowFields, []DataRowFields, error)
	UpdateMultiple(schema *SortTableSchema, hashKeys []HashKey, sortKeys []SortKey, data []DataRow) error
	DeleteMultiple(schema *SortTableSchema, hashKeys []HashKey, sortKeys []SortKey) error

	GetWithSortKey(schema *SortTableSchema, hashKey HashKey, sortKey SortKey) ([]DataRowFields, []DataRowFields, error)
	UpdateWithSortKey(schema *SortTableSchema, hashKey HashKey, sortKey SortKey, data DataRow) error
	DeleteWithSortKey(schema *SortTableSchema, hashKey HashKey, sortKey SortKey) error
}
