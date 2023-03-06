package product

import (
	"time"

	"github.com/sophielizg/go-libs/datastore"
)

type ProductDataRow struct {
	Department  string
	Price       float32
	Quantity    int
	LastUpdated time.Time
}

func (d *ProductDataRow) GetFields() datastore.DataRowFields {
	return datastore.DataRowFields{
		"Department":  d.Department,
		"Price":       d.Price,
		"Quantity":    d.Quantity,
		"LastUpdated": d.LastUpdated,
	}
}

type ProductDataRowFactory struct{}

func (f *ProductDataRowFactory) CreateDefault() *ProductDataRow {
	return &ProductDataRow{}
}

func (f *ProductDataRowFactory) CreateFromFields(fields datastore.DataRowFields) (*ProductDataRow, error) {
	return &ProductDataRow{
		Department:  fields["Department"].(string),
		Price:       fields["Price"].(float32),
		Quantity:    fields["Quantity"].(int),
		LastUpdated: fields["LastUpdated"].(time.Time),
	}, nil
}

func (f *ProductDataRowFactory) GetFieldTypes() datastore.DataRowFieldTypes {
	return datastore.DataRowFieldTypes{
		"Department":  &datastore.StringField{NumChars: 64},
		"Price":       &datastore.FloatField{},
		"Quantity":    &datastore.IntField{},
		"LastUpdated": &datastore.TimeField{},
	}
}
