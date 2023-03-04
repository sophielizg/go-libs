package datastore

import "errors"

type HashTableSchema struct {
	DataRowSchemaFactory
	HashKeySchemaFactory
	Name string
}

func (s *HashTableSchema) Validate() error {
	dataFieldTypes := s.DataRowSchemaFactory.GetFieldTypes()
	for hashField := range s.HashKeySchemaFactory.GetFieldTypes() {
		if _, ok := dataFieldTypes[hashField]; ok {
			return errors.New("DataRow and HashKey cannot share any field names")
		}
	}

	return nil
}

func (s *HashTableSchema) validateDataRowFields(dataRow DataRowFields) error {
	return validateFieldTypes(dataRow, s.DataRowSchemaFactory)
}

func (s *HashTableSchema) validateHashKeyFields(hashKey DataRowFields) error {
	return validateFieldTypes(hashKey, s.HashKeySchemaFactory)
}

type SortTableSchema struct {
	DataRowSchemaFactory
	HashKeySchemaFactory
	SortKeySchemaFactory
	Name string
}
