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

type productHashKeyFactory struct{}

func (f *productHashKeyFactory) CreateDefault() *ProductHashKey {
	return &ProductHashKey{}
}

func (f *productHashKeyFactory) CreateFromFields(fields datastore.DataRowFields) (*ProductHashKey, error) {
	return &ProductHashKey{
		Brand: fields["Brand"].(string),
		Name:  fields["Name"].(string),
	}, nil
}

func (f *productHashKeyFactory) GetFieldTypes() datastore.DataRowFieldTypes {
	return datastore.DataRowFieldTypes{
		"Brand": &datastore.StringField{NumChars: 64},
		"Name":  &datastore.StringField{NumChars: 256},
	}
}

func (f *productHashKeyFactory) GetFieldOptions() datastore.Options {
	return datastore.Options{}
}

func (f *productHashKeyFactory) GetSortOrder() []string {
	return []string{"Brand", "Name"}
}
