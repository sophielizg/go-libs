package datastore

type Option = int8

type Options map[string][]Option

type OptionTypes map[Option]bool

type SupportedOptions map[Option]OptionTypes

func isOptionSupportedForType(fieldType FieldType, fieldOption Option, supported SupportedOptions) bool {
	supportedMap, ok := supported[fieldType.TypeId()]
	if !ok {
		return false
	}

	return supportedMap[fieldOption]
}
