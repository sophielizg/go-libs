package purchase

import (
	"time"

	"github.com/sophielizg/go-libs/datastore"
)

type PurchaseDataRow struct {
	Department  string
	Price       float32
	Quantity    int
	LastUpdated time.Time
}

func (d *PurchaseDataRow) GetFields() datastore.DataRowFields {
	return datastore.DataRowFields{
		"Department":  d.Department,
		"Price":       d.Price,
		"Quantity":    d.Quantity,
		"LastUpdated": d.LastUpdated,
	}
}

type PurchaseDataRowFactory struct{}

func (f *PurchaseDataRowFactory) CreateDefault() *PurchaseDataRow {
	return &PurchaseDataRow{}
}

func (f *PurchaseDataRowFactory) CreateFromFields(fields datastore.DataRowFields) (*PurchaseDataRow, error) {
	return &PurchaseDataRow{
		Department:  fields["Department"].(string),
		Price:       fields["Price"].(float32),
		Quantity:    fields["Quantity"].(int),
		LastUpdated: fields["LastUpdated"].(time.Time),
	}, nil
}

func (f *PurchaseDataRowFactory) GetFieldTypes() datastore.DataRowFieldTypes {
	return datastore.DataRowFieldTypes{
		"Department":  &datastore.StringField{NumChars: 64},
		"Price":       &datastore.FloatField{},
		"Quantity":    &datastore.IntField{},
		"LastUpdated": &datastore.TimeField{},
	}
}
