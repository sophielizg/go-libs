package datastore

var (
	FieldOptions = struct {
		AutoGenerateOption string
	}{
		"AutoGenerateOption",
	}

	FieldOptionTypes = OptionTypes{
		FieldOptions.AutoGenerateOption: true,
	}
)

var DefaultSupportedFieldOptions = SupportedOptions{
	(&IntField{}).TypeName():    FieldOptionTypes,
	(&StringField{}).TypeName(): FieldOptionTypes,
}
