package datastore

const (
	AutoGenerateFieldOption Option = iota
)

var (
	DefaultSupportedFieldOptions = SupportedOptions{
		(&IntField{}).TypeName():    OptionTypes{AutoGenerateFieldOption: true},
		(&StringField{}).TypeName(): OptionTypes{AutoGenerateFieldOption: true},
	}
)

func isFieldOption(x Option) bool {
	return AutoGenerateFieldOption == x
}
