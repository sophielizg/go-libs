package shipping

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	DepartmentKey      = "Department"
	BrandKey           = "Brand"
	NameKey            = "Name"
	PurchaseTimeKey    = "PurchaseTime"
	QuantityKey        = "Quantity"
	ShipmentTimeKey    = "ShipmentTime"
	ShippingAddressKey = "ShippingAddress"
)

type DataRow struct {
	Department      fields.String
	Brand           fields.String
	Name            fields.String
	PurchaseTime    fields.Time
	Quantity        fields.Int
	ShipmentTime    fields.NullTime
	ShippingAddress fields.String
	fieldMutator    *mutator.FieldMutator
}

func (v *DataRow) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(DepartmentKey, &v.Department),
			mutator.WithAddress(BrandKey, &v.Brand),
			mutator.WithAddress(NameKey, &v.Name),
			mutator.WithAddress(PurchaseTimeKey, &v.PurchaseTime),
			mutator.WithAddress(QuantityKey, &v.Quantity),
			mutator.WithAddress(ShipmentTimeKey, &v.ShipmentTime),
			mutator.WithAddress(ShippingAddressKey, &v.ShippingAddress),
		)
	}

	return v.fieldMutator
}

var DataRowSettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(DepartmentKey, 63),
		fields.WithNumBytes(BrandKey, 63),
		fields.WithNumBytes(NameKey, 255),
		fields.WithNumBytes(ShippingAddressKey, 255),
	),
	FieldOrder: fields.OrderedFieldKeys{DepartmentKey, BrandKey, NameKey, PurchaseTimeKey, QuantityKey, ShipmentTimeKey, ShippingAddressKey},
}
