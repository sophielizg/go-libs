package datastore

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

type TableSettings struct {
	Name           string
	EmptyValues    mutator.MappedFieldValues
	DataSettings   *fields.RowSettings
	KeySettings    *fields.RowSettings
	SortFieldNames fields.SortFieldNames
}

func (s *TableSettings) ApplyOption(option func(*TableSettings)) {
	option(s)
}

func NewTableSettings(options ...func(*TableSettings)) *TableSettings {
	settings := &TableSettings{}

	for _, option := range options {
		option(settings)
	}

	return settings
}

func WithTableName(name string) func(*TableSettings) {
	return func(settings *TableSettings) {
		settings.Name = name
	}
}

func WithDataSettings(rowSettings *fields.RowSettings) func(*TableSettings) {
	return func(settings *TableSettings) {
		settings.DataSettings = rowSettings
	}
}

func WithKeySettings(rowSettings *fields.RowSettings) func(*TableSettings) {
	return func(settings *TableSettings) {
		settings.KeySettings = rowSettings
	}
}

func WithSortFieldNames(SortFieldNames fields.SortFieldNames) func(*TableSettings) {
	return func(settings *TableSettings) {
		settings.SortFieldNames = SortFieldNames
	}
}

func WithEntry[E any, PE mutator.Mutatable[E]]() func(*TableSettings) {
	return func(settings *TableSettings) {
		empty := mutator.MutatableFactory[E, PE]{}.Create()
		settings.EmptyValues = empty.Mutator().GetFields()
	}
}
