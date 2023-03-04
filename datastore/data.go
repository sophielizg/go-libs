package datastore

type DataRowFields map[string]interface{}

type DataRowFieldTypes map[string]FieldType

type DataRow interface {
	GetFields() DataRowFields
}

type HashKey interface {
	DataRow
}

type SortKey interface {
	HashKey
}
