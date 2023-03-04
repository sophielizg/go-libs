package datastore

type WriteOption interface {
	Option
}

var DefaultSupportedWriteOptions = SupportedOptions[WriteOption]{
	&TimeField{}: []WriteOption{&GenerateCurrentTimeOption{}},
}

type GenerateCurrentTimeOption struct{}

func (o *GenerateCurrentTimeOption) Name() string {
	return "GenerateCurrentTimeOption"
}

func (o *GenerateCurrentTimeOption) OverrideSupported() bool {
	return false
}
