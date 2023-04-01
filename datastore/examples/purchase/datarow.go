package purchase

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	departmentKey  = "Department"
	priceKey       = "Price"
	quantityKey    = "Quantity"
	lastUpdatedKey = "LastUpdated"
)

type PurchaseDataRow struct {
	Department   fields.String
	Price        fields.Float
	Quantity     fields.Int
	LastUpdated  fields.Time
	fieldMutator *mutator.FieldMutator
}

func (v *PurchaseDataRow) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(departmentKey, &v.Department),
			mutator.WithAddress(priceKey, &v.Price),
			mutator.WithAddress(quantityKey, &v.Quantity),
			mutator.WithAddress(lastUpdatedKey, &v.LastUpdated),
		)
	}

	return v.fieldMutator
}

var PurchaseDataRowSettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(departmentKey, 63),
	),
	FieldOrder: fields.OrderedFieldKeys{departmentKey, priceKey, quantityKey, lastUpdatedKey},
}
