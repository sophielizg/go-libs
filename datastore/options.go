package datastore

type Option interface {
	// Unique name to identify the option
	Name() string
	OverrideSupported() bool
}

type Options[O Option] map[string]O

type SupportedOptions[O Option] map[FieldType][]O

func isOptionSupportedForType[O Option](fieldType FieldType, fieldOption O, supported SupportedOptions[O]) bool {
	if fieldOption.OverrideSupported() {
		return true
	}

	supportedList, ok := supported[fieldType]
	if !ok {
		return false
	}

	for _, option := range supportedList {
		if fieldOption.Name() == option.Name() {
			return true
		}
	}

	return false
}
