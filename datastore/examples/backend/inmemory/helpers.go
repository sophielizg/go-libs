package inmemory

import (
	"encoding/json"
	"math/rand"

	"github.com/sophielizg/go-libs/datastore"
)

// Used to calculate random key
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateStringKey(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateIntKey() int {
	return rand.Int()
}

func stringifyHashKey(hashKey datastore.DataRowFields) (string, error) {
	bytes, err := json.Marshal(hashKey)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func unstringifyHashKey(stringified string) (datastore.DataRowFields, error) {
	bytes := []byte(stringified)
	hashKey := datastore.DataRowFields{}
	err := json.Unmarshal(bytes, &hashKey)
	if err != nil {
		return nil, err
	}

	return hashKey, nil
}

func applyKeyOptions(key datastore.DataRowFields, fieldTypes datastore.DataRowFieldTypes, optionsForFields datastore.Options) (datastore.DataRowFields, error) {
	for fieldName, options := range optionsForFields {
		for _, option := range options {
			if option == datastore.FieldOptions.AutoGenerateOption {
				if _, ok := fieldTypes[fieldName].(*datastore.IntField); ok {
					key[fieldName] = generateIntKey()
				} else if stringType, ok := fieldTypes[fieldName].(*datastore.StringField); ok {
					key[fieldName] = generateStringKey(stringType.NumChars)
				}
			}
		}
	}

	return key, nil
}

func shouldApplyKeyOptions(key datastore.DataRowFields, fieldTypes datastore.DataRowFieldTypes, optionsForFields datastore.Options) bool {
	for fieldName := range optionsForFields {
		if key[fieldName] == nil {
			return true
		}

		_, ok := fieldTypes[fieldName].(*datastore.IntField)
		if ok && key[fieldName] == 0 {
			return true
		}

		_, ok = fieldTypes[fieldName].(*datastore.StringField)
		if ok && key[fieldName] == "" {
			return true
		}
	}

	return false
}
