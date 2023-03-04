package datastore

import "errors"

type ScanTableSchema struct {
	DataRowSchemaFactory
	Name string
}

func (s *ScanTableSchema) validateDataRowFields(dataRow DataRowFields) error {
	return validateFieldTypes(dataRow, s.DataRowSchemaFactory)
}

type AppendTableSchema struct {
	ScanTableSchema
	SupportedWriteOptions SupportedOptions[WriteOption]
}

func (s *AppendTableSchema) validateWriteOptions(writeOptions Options[WriteOption]) error {
	return validateOptions(s.DataRowSchemaFactory.GetFieldTypes(), writeOptions, s.SupportedWriteOptions)
}

type HashTableSchema struct {
	ScanTableSchema
	HashKeySchemaFactory
	SupportedFieldOptions SupportedOptions[FieldOption]
	SupportedWriteOptions SupportedOptions[WriteOption]
}

func (s *HashTableSchema) Validate() error {
	dataFieldTypes := s.DataRowSchemaFactory.GetFieldTypes()
	for hashField := range s.HashKeySchemaFactory.GetFieldTypes() {
		if _, ok := dataFieldTypes[hashField]; ok {
			return errors.New("DataRow and HashKey cannot share any field names")
		}
	}

	return s.validateFieldOptions()
}

func (s *HashTableSchema) validateHashKeyFields(hashKey DataRowFields) error {
	return validateFieldTypes(hashKey, s.HashKeySchemaFactory)
}

func (s *HashTableSchema) validateFieldOptions() error {
	return validateOptions(s.HashKeySchemaFactory.GetFieldTypes(), s.HashKeySchemaFactory.GetFieldOptions(), s.SupportedFieldOptions)
}

func (s *HashTableSchema) validateWriteOptions(writeOptions Options[WriteOption]) error {
	return validateOptions(s.DataRowSchemaFactory.GetFieldTypes(), writeOptions, s.SupportedWriteOptions)
}

func (s *HashTableSchema) validateUpdateOptions(updateOptions Options[UpdateOption]) error {
	writeOptions := Options[WriteOption]{}
	includeExcludeOptions := Options[*IncludeExcludeOption]{}

	for fieldName, option := range updateOptions {
		if v, ok := option.(*IncludeExcludeOption); ok {
			includeExcludeOptions[fieldName] = v
		} else {
			writeOptions[fieldName] = option
		}
	}

	err := validateIncludeExcludeOptions(includeExcludeOptions)
	if err != nil {
		return err
	}

	return validateOptions(s.HashKeySchemaFactory.GetFieldTypes(), writeOptions, s.SupportedWriteOptions)
}

type SortTableSchema struct {
	HashTableSchema
	SortKeySchemaFactory
}
