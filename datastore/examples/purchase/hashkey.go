package purchase

import (
	"github.com/sophielizg/go-libs/datastore"
)

type PurchaseHashKey struct {
	CustomerName string
}

func (d *PurchaseHashKey) GetFields() datastore.DataRowFields {
	return datastore.DataRowFields{
		"CustomerName": d.CustomerName,
	}
}

type PurchaseHashKeyFactory struct{}

func (f *PurchaseHashKeyFactory) CreateDefault() *PurchaseHashKey {
	return &PurchaseHashKey{}
}

func (f *PurchaseHashKeyFactory) CreateFromFields(fields datastore.DataRowFields) (*PurchaseHashKey, error) {
	return &PurchaseHashKey{
		CustomerName: fields["CustomerName"].(string),
	}, nil
}

func (f *PurchaseHashKeyFactory) GetFieldTypes() datastore.DataRowFieldTypes {
	return datastore.DataRowFieldTypes{
		"CustomerName": &datastore.StringField{NumChars: 256},
	}
}

func (f *PurchaseHashKeyFactory) GetFieldOptions() datastore.Options {
	return datastore.Options{}
}

func (f *PurchaseHashKeyFactory) GetSortOrder() []string {
	return []string{"CustomerName"}
}
