package datastore

type HashTableBackend interface {
	SupportedFieldOptions() SupportedFieldOptions
	ValidateHashTableSchema(schema *HashTableSchema) error
	CreateOrUpdateHashTableSchema(schema *HashTableSchema) error
	GetMultiple(schema *HashTableSchema, hashKeys []HashKey) ([]DataRowFields, error)
	AddMultiple(schema *HashTableSchema, hashKeys []HashKey, data []DataRow, options WriteOptions) ([]DataRowFields, error)
	UpdateMultiple(schema *HashTableSchema, hashKeys []HashKey, data []DataRow, options UpdateOptions) error
	DeleteMutiple(schema *HashTableSchema, hashKeys []HashKey) error
}

type SortTableBackend interface {
	HashTableBackend
	ValidateSortTableSchema(schema *SortTableSchema) error
	CreateOrUpdateSortTableSchema(schema *SortTableSchema) error
	GetWithSortKey(schema *SortTableSchema, hashKey *HashKey, sortKey *SortKey) ([]DataRowFields, error)
}
