package datastore

import "errors"

func validateFieldTypes(dataRowFields DataRowFields, fieldTypes FieldTypesFactory) error {
	dataFieldTypes := fieldTypes.GetFieldTypes()

	if len(dataFieldTypes) != len(dataRowFields) {
		return errors.New("DataRow must have same number of type definitions and data fields")
	}

	for fieldName, field := range dataRowFields {
		fieldType, ok := dataFieldTypes[fieldName]
		if !ok {
			return errors.New("DataRow.GetFieldTypes() must return all data fields")
		} else if fieldType == nil {
			return errors.New("DataRow must define types for every field")
		} else if !fieldType.IsType(field) {
			return errors.New("Each DataRow field must satisfy its type")
		}
	}

	return nil
}

func validateFieldOptions(fieldOptions HashKeySchemaFactory, supported SupportedFieldOptions) error {
	dataFieldOptions := fieldOptions.GetFieldOptions()
	dataFieldTypes := fieldOptions.GetFieldTypes()

	if len(dataFieldOptions) > len(dataFieldTypes) {
		return errors.New("DataRow.GetFieldOptions() must not return more fields than have types in DataRow")
	}

	for fieldName, fieldOption := range dataFieldOptions {
		fieldType, ok := dataFieldTypes[fieldName]
		if !ok {
			return errors.New("DataRow.GetFieldTypes() must return all data fields")
		} else if fieldType == nil {
			return errors.New("DataRow must define types for every field")
		} else if !isFieldOptionSupportedForType(fieldType, fieldOption, supported) {
			return errors.New("FieldOption is not supported for specified type")
		}
	}

	return nil
}
