package datastore

const (
	AutoGenerateFieldOption Option = iota
)

var (
	DefaultSupportedFieldOptions = SupportedOptions{
		intFieldId:    OptionTypes{AutoGenerateFieldOption: true},
		stringFieldId: OptionTypes{AutoGenerateFieldOption: true},
	}
)

func isFieldOption(x Option) bool {
	return AutoGenerateFieldOption == x
}
