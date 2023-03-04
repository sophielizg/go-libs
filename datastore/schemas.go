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
}

type HashTableSchema struct {
	ScanTableSchema
	HashKeySchemaFactory
	SupportedFieldOptions SupportedOptions[FieldOption]
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

func (s *HashTableSchema) validateOptionsForFieldTypes(options Options[FieldOption], fieldTypesList ...DataRowFieldTypes) error {
	merged := DataRowFieldTypes{}
	for _, fieldTypes := range fieldTypesList {
		for k, v := range fieldTypes {
			merged[k] = v
		}
	}

	return validateOptions(merged, options, s.SupportedFieldOptions)
}

func (s *HashTableSchema) validateFieldOptions() error {
	return s.validateOptionsForFieldTypes(
		s.HashKeySchemaFactory.GetFieldOptions(),
		s.DataRowSchemaFactory.GetFieldTypes(),
		s.HashKeySchemaFactory.GetFieldTypes(),
	)
}

type SortTableSchema struct {
	HashTableSchema
	SortKeySchemaFactory
}

func (s *SortTableSchema) Validate() error {
	seen := map[string]bool{}

	for dataField := range s.DataRowSchemaFactory.GetFieldTypes() {
		seen[dataField] = true
	}

	hashFieldTypes := s.HashKeySchemaFactory.GetFieldTypes()
	sortFieldTypes := s.SortKeySchemaFactory.GetFieldTypes()
	for _, fieldTypes := range []DataRowFieldTypes{hashFieldTypes, sortFieldTypes} {
		for fieldName := range fieldTypes {
			if seen[fieldName] {
				return errors.New("DataRow, HashKey, and SortKey cannot share any field names")
			}
			seen[fieldName] = true
		}
	}

	return s.validateFieldOptions()
}

func (s *SortTableSchema) validateFieldOptions() error {
	return s.validateOptionsForFieldTypes(
		s.HashKeySchemaFactory.GetFieldOptions(),
		s.DataRowSchemaFactory.GetFieldTypes(),
		s.HashKeySchemaFactory.GetFieldTypes(),
		s.SortKeySchemaFactory.GetFieldTypes(),
	)
}

func (s *SortTableSchema) validateSortKeyFields(sortKey DataRowFields) error {
	return validateFieldTypes(sortKey, s.SortKeySchemaFactory)
}
