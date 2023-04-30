package fields

type OrderedFieldKeys = []string

type RowSettings struct {
	FieldSettings FieldSettings
	FieldOrder    OrderedFieldKeys
}

type SortFieldNames = []string
