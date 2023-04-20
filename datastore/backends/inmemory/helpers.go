package inmemory

import (
	"encoding/json"

	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

func validateAutoGenerateSettings(settings *fields.DataRowSettings) error {
	if settings == nil {
		return nil
	}

	for _, fieldSetting := range settings.FieldSettings {
		if fieldSetting.AutoGenerate {
			return AutoGenerateNotSupportedError
		}
	}

	return nil
}

func stringifyKey(key mutator.MappedFieldValues) (string, error) {
	bytes, err := json.Marshal(key)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func unstringifyKey(stringified string) (mutator.MappedFieldValues, error) {
	bytes := []byte(stringified)
	key := mutator.MappedFieldValues{}
	err := json.Unmarshal(bytes, &key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
