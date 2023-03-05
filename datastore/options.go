package datastore

type Options map[string][]string

type OptionTypes map[string]bool

type SupportedOptions map[FieldType]OptionTypes

func isOptionSupportedForType(fieldType FieldType, fieldOption string, supported SupportedOptions) bool {
	supportedMap, ok := supported[fieldType]
	if !ok {
		return false
	}

	return supportedMap[fieldOption]
}
