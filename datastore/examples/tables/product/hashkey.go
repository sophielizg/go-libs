package product

import (
	"github.com/sophielizg/go-libs/datastore"
)

type ProductHashKey struct {
	Brand string
	Name  string
}

func (d *ProductHashKey) GetFields() datastore.DataRowFields {
	return datastore.DataRowFields{
		"Brand": d.Brand,
		"Name":  d.Name,
	}
}

type ProductHashKeyFactory struct{}

func (f *ProductHashKeyFactory) CreateDefault() *ProductHashKey {
	return &ProductHashKey{}
}

func (f *ProductHashKeyFactory) CreateFromFields(fields datastore.DataRowFields) (*ProductHashKey, error) {
	return &ProductHashKey{
		Brand: fields["Brand"].(string),
		Name:  fields["Name"].(string),
	}, nil
}

func (f *ProductHashKeyFactory) GetFieldTypes() datastore.DataRowFieldTypes {
	return datastore.DataRowFieldTypes{
		"Brand": &datastore.StringField{NumChars: 64},
		"Name":  &datastore.StringField{NumChars: 256},
	}
}

func (f *ProductHashKeyFactory) GetFieldOptions() datastore.Options {
	return datastore.Options{}
}

func (f *ProductHashKeyFactory) GetSortOrder() []string {
	return []string{"Brand", "Name"}
}
