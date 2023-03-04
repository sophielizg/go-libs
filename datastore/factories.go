package datastore

type FieldTypesFactory interface {
	GetFieldTypes() DataRowFieldTypes
}

type FieldOptionsFactory interface {
	GetFieldOptions() Options[FieldOption]
}

type SortOrderFactory interface {
	GetSortOrder() []string
}

type DataRowSchemaFactory interface {
	FieldTypesFactory
}

type DataRowFactory[V DataRow] interface {
	DataRowSchemaFactory
	CreateFromFields(fields DataRowFields) (V, error)
}

type HashKeySchemaFactory interface {
	DataRowSchemaFactory
	FieldOptionsFactory
}

type HashKeyFactory[H HashKey] interface {
	DataRowFactory[H]
	HashKeySchemaFactory
}

type SortKeySchemaFactory interface {
	HashKeySchemaFactory
	SortOrderFactory
}

type SortKeyFactory[S SortKey] interface {
	HashKeyFactory[S]
	SortKeySchemaFactory
}
