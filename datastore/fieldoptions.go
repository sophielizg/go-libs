package datastore

const (
	autoGenerateOption Option = iota
)

var (
	FieldOptions = struct {
		AutoGenerateOption Option
	}{
		autoGenerateOption,
	}

	DefaultSupportedFieldOptions = SupportedOptions{
		(&IntField{}).TypeName():    OptionTypes{autoGenerateOption: true},
		(&StringField{}).TypeName(): OptionTypes{autoGenerateOption: true},
	}
)

func isOption(x Option) bool {
	return autoGenerateOption == x
}
