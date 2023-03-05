package datastore

import "errors"

type BaseTableSchema struct {
	DataRowSchemaFactory
	Name string
}

func (s *BaseTableSchema) validateDataRowFields(dataRow DataRowFields) error {
	return validateFields(dataRow, s.DataRowSchemaFactory)
}

type AppendTableSchema struct {
	BaseTableSchema
}

type HashTableSchema struct {
	BaseTableSchema
	HashKeySchemaFactory  KeySchemaFactory
	SupportedFieldOptions SupportedOptions
}

func (s *HashTableSchema) Validate() error {
	dataFieldTypes := s.DataRowSchemaFactory.GetFieldTypes()
	hashFieldTypes := s.HashKeySchemaFactory.GetFieldTypes()

	for hashField := range hashFieldTypes {
		if _, ok := dataFieldTypes[hashField]; ok {
			return errors.New("DataRow and HashKey cannot share any field names")
		}
	}

	err := validateKeyFieldTypes(hashFieldTypes)
	if err != nil {
		return err
	}

	return s.validateFieldOptions()
}

func (s *HashTableSchema) validateHashKeyFields(hashKey DataRowFields) error {
	return validateFields(hashKey, s.HashKeySchemaFactory)
}

func (s *HashTableSchema) validateOptionsForFieldTypes(options Options, fieldTypesList ...DataRowFieldTypes) error {
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
	SortKeySchemaFactory KeySchemaFactory
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

	err := validateKeyFieldTypes(hashFieldTypes)
	if err != nil {
		return err
	}

	err = validateKeyFieldTypes(sortFieldTypes)
	if err != nil {
		return err
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
	return validateFields(sortKey, s.SortKeySchemaFactory)
}
