package purchase

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	DepartmentKey  = "Department"
	PriceKey       = "Price"
	QuantityKey    = "Quantity"
	LastUpdatedKey = "LastUpdated"
)

type DataRow struct {
	Department   fields.String
	Price        fields.Float
	Quantity     fields.Int
	LastUpdated  fields.Time
	fieldMutator *mutator.FieldMutator
}

func (v *DataRow) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(DepartmentKey, &v.Department),
			mutator.WithAddress(PriceKey, &v.Price),
			mutator.WithAddress(QuantityKey, &v.Quantity),
			mutator.WithAddress(LastUpdatedKey, &v.LastUpdated),
		)
	}

	return v.fieldMutator
}

var DataRowSettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(DepartmentKey, 63),
	),
	FieldOrder: fields.OrderedFieldKeys{DepartmentKey, PriceKey, QuantityKey, LastUpdatedKey},
}
