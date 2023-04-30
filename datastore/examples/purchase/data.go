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

type Data struct {
	Department  fields.String
	Price       fields.Float
	Quantity    fields.Int
	LastUpdated fields.Time
}

func (d *Data) Mutator() *mutator.FieldMutator {
	return mutator.NewFieldMutator(
		mutator.WithAddress(DepartmentKey, &d.Department),
		mutator.WithAddress(PriceKey, &d.Price),
		mutator.WithAddress(QuantityKey, &d.Quantity),
		mutator.WithAddress(LastUpdatedKey, &d.LastUpdated),
	)
}

var DataSettings = &fields.RowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(DepartmentKey, 63),
	),
	FieldOrder: fields.OrderedFieldKeys{DepartmentKey, PriceKey, QuantityKey, LastUpdatedKey},
}
