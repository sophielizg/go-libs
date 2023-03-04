package datastore

type FieldOption interface {
	// Unique name to identify the option
	Name() string
}

type SupportedFieldOptions map[FieldType][]FieldOption

var DefaultSupportedFieldOptions = SupportedFieldOptions{
	&IntField{}:    []FieldOption{&AutoGenerateOption{}},
	&StringField{}: []FieldOption{&AutoGenerateOption{}},
}

func isFieldOptionSupportedForType(fieldType FieldType, fieldOption FieldOption, supported SupportedFieldOptions) bool {
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

type AutoGenerateOption struct{}

func (o *AutoGenerateOption) Name() string {
	return "AutoGenerateOption"
}
