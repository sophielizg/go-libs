package purchase

import (
	"time"

	"github.com/sophielizg/go-libs/datastore"
)

type PurchaseSortKeyComparators struct {
	PurchaseTime []string
	ItemBrand    []string
	ItemName     []string
}

type PurchaseSortKey struct {
	PurchaseTime time.Time
	ItemBrand    string
	ItemName     string
	Comparators  PurchaseSortKeyComparators
}

func (d *PurchaseSortKey) GetFields() datastore.DataRowFields {
	return datastore.DataRowFields{
		"PurchaseTime": d.PurchaseTime,
		"ItemBrand":    d.ItemBrand,
		"ItemName":     d.ItemName,
	}
}

func (d *PurchaseSortKey) GetComparators() datastore.Options {
	return datastore.Options{
		"PurchaseTime": d.Comparators.PurchaseTime,
		"ItemBrand":    d.Comparators.ItemBrand,
		"ItemName":     d.Comparators.ItemName,
	}
}

type PurchaseSortKeyFactory struct{}

func (f *PurchaseSortKeyFactory) CreateDefault() *PurchaseSortKey {
	return &PurchaseSortKey{}
}

func (f *PurchaseSortKeyFactory) CreateFromFields(fields datastore.DataRowFields) (*PurchaseSortKey, error) {
	return &PurchaseSortKey{
		PurchaseTime: fields["PurchaseTime"].(time.Time),
		ItemBrand:    fields["ItemBrand"].(string),
		ItemName:     fields["ItemName"].(string),
	}, nil
}

func (f *PurchaseSortKeyFactory) GetFieldTypes() datastore.DataRowFieldTypes {
	return datastore.DataRowFieldTypes{
		"PurchaseTime": &datastore.TimeField{},
		"ItemBrand":    &datastore.StringField{NumChars: 64},
		"ItemName":     &datastore.StringField{NumChars: 256},
	}
}

func (f *PurchaseSortKeyFactory) GetFieldOptions() datastore.Options {
	return datastore.Options{}
}

func (f *PurchaseSortKeyFactory) GetSortOrder() []string {
	return []string{"PurchaseTime", "ItemBrand", "ItemName"}
}
