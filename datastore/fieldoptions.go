package datastore

type FieldOption interface {
	Option
}

var DefaultSupportedFieldOptions = SupportedOptions[FieldOption]{
	&IntField{}:    []FieldOption{&AutoGenerateOption{}},
	&StringField{}: []FieldOption{&AutoGenerateOption{}},
}

type AutoGenerateOption struct{}

func (o *AutoGenerateOption) Name() string {
	return "AutoGenerateOption"
}

func (o *AutoGenerateOption) OverrideSupported() bool {
	return false
}
