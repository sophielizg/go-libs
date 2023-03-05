package datastore

type FieldTypesFactory interface {
	GetFieldTypes() DataRowFieldTypes
}

type FieldOptionsFactory interface {
	GetFieldOptions() Options
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
	CreateDefault() V
}

type KeySchemaFactory interface {
	DataRowSchemaFactory
	FieldOptionsFactory
	SortOrderFactory
}

type KeyFactory[H HashKey] interface {
	DataRowFactory[H]
	KeySchemaFactory
}
