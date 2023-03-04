package datastore

type DataRowFields map[string]interface{}

type DataRowFieldTypes map[string]FieldType

type KeyFieldOptions map[string]FieldOption

type DataRow interface {
	GetFields() DataRowFields
}

type HashKey interface {
	DataRow
}

type SortKey interface {
	HashKey
}
