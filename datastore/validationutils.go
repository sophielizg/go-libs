package datastore

import "errors"

func validateDataRowFieldName(fieldName string, dataRowFields DataRowFields) error {
	_, ok := dataRowFields[fieldName]
	if !ok {
		return errors.New("fieldName must be a field in DataRow")
	}

	return nil
}

func validateFields(dataRowFields DataRowFields, fieldTypes FieldTypesFactory) error {
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

func validateKeyFieldTypes(fieldTypes DataRowFieldTypes) error {
	for _, fieldType := range fieldTypes {
		if !fieldType.IsComparable() {
			return errors.New("All types in key must be comparable")
		}
	}

	return nil
}

func validateOptions[O Option](dataFieldTypes DataRowFieldTypes, dataFieldOptions Options[O], supported SupportedOptions[O]) error {
	if len(dataFieldOptions) > len(dataFieldTypes) {
		return errors.New("DataRow.GetFieldOptions() must not return more fields than have types in DataRow")
	}

	for fieldName, fieldOption := range dataFieldOptions {
		fieldType, ok := dataFieldTypes[fieldName]
		if !ok {
			return errors.New("DataRow.GetFieldTypes() must return all data fields")
		} else if fieldType == nil {
			return errors.New("DataRow must define types for every field")
		} else if !isOptionSupportedForType(fieldType, fieldOption, supported) {
			return errors.New("FieldOption is not supported for specified type")
		}
	}

	return nil
}
