package datastore

import "errors"

type SortTable[V DataRow, H HashKey, S SortKey] struct {
	HashTable[V, H]
	Backend        SortTableBackend
	SortKeyFactory SortKeyFactory[S]
	Name           string
	schema         *SortTableSchema
}

func (t *SortTable[V, H, S]) getSchema() *SortTableSchema {
	if t.schema == nil {
		t.schema = &SortTableSchema{
			Name:                 t.Name,
			DataRowSchemaFactory: t.DataRowFactory,
			HashKeySchemaFactory: t.HashKeyFactory,
			SortKeySchemaFactory: t.SortKeyFactory,
		}
	}

	return t.schema
}

func (t *SortTable[V, H, S]) ValidateSchema() error {
	err := t.HashTable.ValidateSchema()
	if err != nil {
		return err
	}

	sortKeyFieldTypes := t.SortKeyFactory.GetFieldTypes()
	sortOrder := t.SortKeyFactory.GetSortOrder()
	if len(sortKeyFieldTypes) != len(sortOrder) {
		return errors.New("SortKey field types and sort order must have the same number of fields")
	}

	for _, field := range sortOrder {
		fieldType, ok := sortKeyFieldTypes[field]
		if !ok {
			return errors.New("SortKey.GetFieldTypes() must return all fields in sort order")
		} else if fieldType == nil {
			return errors.New("SortKey must define types for all fields in sort order")
		} else if !fieldType.IsComparable() {
			return errors.New("All SortKey fields must be comparable")
		}
	}

	return t.Backend.ValidateSortTableSchema(t.getSchema())
}

func (t *SortTable[V, H, S]) ValidateSortKey(sortKey S) error {
	foundNil := false
	sortKeyFields := sortKey.GetFields()
	sortOrder := t.SortKeyFactory.GetSortOrder()

	for i, field := range sortOrder {
		val, ok := sortKeyFields[field]
		if !ok {
			return errors.New("SortKey.GetFields() must return all fields in sort order")
		} else if val == nil {
			if i == 0 {
				return errors.New("Leftmost SortKey field must be populated")
			}
			foundNil = true
		} else if foundNil {
			return errors.New("All SortKey fields on the left side must be populated")
		}
	}

	return nil
}

// func (t *SortTable[V, H, S]) GetWithSortKey(hashKey H, sortKey S) ([]V, error) {
// 	// call db implementation
// }

// func (t *SortTable[V, H, S]) UpdateWithSortKey(hashKey H, sortKey S, data V, options UpdateOptions) error {
// 	// call db implementation
// }

// func (t *SortTable[V, H, S]) DeleteWithSortKey(hashKey H, sortKey S) error {
// 	// call db implementation
// }
