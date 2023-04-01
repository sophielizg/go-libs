package product

import (
	"github.com/sophielizg/go-libs/datastore/fields"
)

const (
	departmentKey  = "Department"
	priceKey       = "Price"
	quantityKey    = "Quantity"
	lastUpdatedKey = "LastUpdated"
)

type ProductDataRow struct {
	Department  fields.String
	Price       fields.Float
	Quantity    fields.Int
	LastUpdated fields.Time
	builder     *fields.DataRowBuilder
}

func (v *ProductDataRow) Builder() *fields.DataRowBuilder {
	if v.builder == nil {
		v.builder = fields.NewDataRowBuilder(
			fields.WithAddress(departmentKey, &v.Department),
			fields.WithAddress(priceKey, &v.Price),
			fields.WithAddress(quantityKey, &v.Quantity),
			fields.WithAddress(lastUpdatedKey, &v.LastUpdated),
		)
	}

	return v.builder
}

var ProductDataRowSettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(departmentKey, 63),
	),
	FieldOrder: fields.OrderedFieldKeys{departmentKey, priceKey, quantityKey, lastUpdatedKey},
}
