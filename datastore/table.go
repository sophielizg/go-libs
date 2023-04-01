package datastore

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

type Table[B any] interface {
	Init()
	GetSettings() *TableSettings
	SetBackend(tableBackend B)
}

type TableSettings struct {
	Name            string
	DataRowSettings *fields.DataRowSettings
	HashKeySettings *fields.DataRowSettings
	SortKeySettings *fields.DataRowSettings
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

func WithDataRowSettings(dataRowSettings *fields.DataRowSettings) func(*TableSettings) {
	return func(settings *TableSettings) {
		settings.DataRowSettings = dataRowSettings
	}
}

func WithHashKeySettings(dataRowSettings *fields.DataRowSettings) func(*TableSettings) {
	return func(settings *TableSettings) {
		settings.HashKeySettings = dataRowSettings
	}
}

func WithSortKeySettings(dataRowSettings *fields.DataRowSettings) func(*TableSettings) {
	return func(settings *TableSettings) {
		settings.SortKeySettings = dataRowSettings
	}
}

func WithDataRow[V any, PV mutator.Mutatable[V]]() func(*TableSettings) {
	return func(settings *TableSettings) {
		empty := mutator.MutatableFactory[V, PV]{}.Create()
		settings.DataRowSettings.EmptyValues = empty.Mutator().GetFields()
	}
}

func WithHashKey[H any, PH mutator.Mutatable[H]]() func(*TableSettings) {
	return func(settings *TableSettings) {
		empty := mutator.MutatableFactory[H, PH]{}.Create()
		settings.HashKeySettings.EmptyValues = empty.Mutator().GetFields()
	}
}

func WithSortKey[S any, PS mutator.Mutatable[S]]() func(*TableSettings) {
	return func(settings *TableSettings) {
		empty := mutator.MutatableFactory[S, PS]{}.Create()
		settings.SortKeySettings.EmptyValues = empty.Mutator().GetFields()
	}
}
