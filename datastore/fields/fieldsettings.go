package fields

type Setting = int8

const (
	AutoGenerateOption Setting = iota
)

type FieldSettings = map[string]*FieldSetting

type FieldSetting struct {
	NumBytes     int
	AutoGenerate bool
}

func NewFieldSettings(options ...func(FieldSettings)) FieldSettings {
	settings := FieldSettings{}

	for _, option := range options {
		option(settings)
	}

	return settings
}

func settingForFieldName(settings FieldSettings, fieldName string) *FieldSetting {
	if settings[fieldName] == nil {
		settings[fieldName] = &FieldSetting{}
	}

	return settings[fieldName]
}

func WithAutoGenerate(fieldName string) func(settings FieldSettings) {
	return func(settings FieldSettings) {
		setting := settingForFieldName(settings, fieldName)
		setting.AutoGenerate = true
	}
}

func WithNumBytes(fieldName string, numBytes int) func(settings FieldSettings) {
	return func(settings FieldSettings) {
		setting := settingForFieldName(settings, fieldName)
		setting.NumBytes = numBytes
	}
}
