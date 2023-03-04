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

	// Configuration
	SupportedWriteOptions() SupportedOptions[WriteOption]

	// Schema
	ValidateAppendTableSchema(schema *AppendTableSchema) error
	CreateOrUpdateAppendTableSchema(schema *AppendTableSchema) error

	// Data operations
	AppendMultiple(schema *AppendTableSchema, data []DataRow, options Options[WriteOption]) error
}

type HashTableBackend interface {
	ScanTableBackend

	// Configuration
	SupportedFieldOptions() SupportedOptions[FieldOption]
	SupportedWriteOptions() SupportedOptions[WriteOption]

	// Schema
	ValidateHashTableSchema(schema *HashTableSchema) error
	CreateOrUpdateHashTableSchema(schema *HashTableSchema) error

	// Data operations
	GetMultiple(schema *HashTableSchema, hashKeys []HashKey) ([]DataRowFields, error)
	AddMultiple(schema *HashTableSchema, hashKeys []HashKey, data []DataRow, options Options[WriteOption]) ([]DataRowFields, error)
	UpdateMultiple(schema *HashTableSchema, hashKeys []HashKey, data []DataRow, options Options[UpdateOption]) error
	DeleteMutiple(schema *HashTableSchema, hashKeys []HashKey) error
}

type SortTableBackend interface {
	HashTableBackend

	// Schema
	ValidateSortTableSchema(schema *SortTableSchema) error
	CreateOrUpdateSortTableSchema(schema *SortTableSchema) error

	// Data operations
	GetWithSortKey(schema *SortTableSchema, hashKey HashKey, sortKey SortKey) ([]DataRowFields, []DataRowFields, error)
	UpdateWithSortKey(schema *SortTableSchema, hashKey HashKey, sortKey SortKey, data DataRow, options Options[UpdateOption]) error
	DeleteWithSortKey(schema *SortTableSchema, hashKey HashKey, sortKey SortKey) error
}
