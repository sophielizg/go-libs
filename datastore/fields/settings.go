package fields

import "github.com/sophielizg/go-libs/datastore/mutator"

type Setting = int8

const (
	AutoGenerateOption Setting = iota
)

type DataRowSettings struct {
	EmptyValues   mutator.MappedFieldValues
	FieldSettings FieldSettings
	FieldOrder    OrderedFieldKeys
}

type FieldSettings = map[string]*FieldSetting

type FieldSetting struct {
	NumBytes     int
	AutoGenerate bool
}

func NewFieldSettings(options ...func(FieldSettings)) FieldSettings {
	settings := FieldSettings{}

	for _, option := range options {
		option(settings)
	}

	return settings
}

func settingForFieldName(settings FieldSettings, fieldName string) *FieldSetting {
	if settings[fieldName] == nil {
		settings[fieldName] = &FieldSetting{}
	}

	return settings[fieldName]
}

func WithAutoGenerate(fieldName string) func(settings FieldSettings) {
	return func(settings FieldSettings) {
		setting := settingForFieldName(settings, fieldName)
		setting.AutoGenerate = true
	}
}

func WithNumBytes(fieldName string, numBytes int) func(settings FieldSettings) {
	return func(settings FieldSettings) {
		setting := settingForFieldName(settings, fieldName)
		setting.NumBytes = numBytes
	}
}

func MergeSettings(settingsList ...FieldSettings) FieldSettings {
	merged := FieldSettings{}

	for _, settings := range settingsList {
		for fieldName, setting := range settings {
			merged[fieldName] = setting
		}
	}

	return merged
}
