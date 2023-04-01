package datastore

import "github.com/sophielizg/go-libs/datastore/fields"

type Table[B any] interface {
	Init()
	GetSettings() *TableSettings
	SetBackend(tableBackend B)
}

type TableSettings struct {
	Name            string
	DataRowSettings fields.DataRowSettings
	HashKeySettings fields.DataRowSettings
	SortKeySettings fields.DataRowSettings
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

func WithDataRow[V any, PV DataRow[V]]() func(*TableSettings) {
	return func(settings *TableSettings) {
		dataRow := DataRowFactory[V, PV]{}.Create()
		settings.DataRowSettings.EmptyValues = dataRow.Mutator().GetFields()
	}
}

func WithHashKey[H any, PH HashKey[H]]() func(*TableSettings) {
	return func(settings *TableSettings) {
		hashKey := DataRowFactory[H, PH]{}.Create()
		settings.HashKeySettings.EmptyValues = hashKey.Mutator().GetFields()
	}
}

func WithSortKey[S any, PS SortKey[S]]() func(*TableSettings) {
	return func(settings *TableSettings) {
		sortKey := DataRowFactory[S, PS]{}.Create()
		settings.SortKeySettings.EmptyValues = sortKey.Mutator().GetFields()
	}
}
