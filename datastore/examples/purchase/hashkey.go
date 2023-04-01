package purchase

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	customerNameKey = "CustomerName"
)

type PurchaseHashKey struct {
	CustomerName fields.String
	fieldMutator *mutator.FieldMutator
}

func (v *PurchaseHashKey) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(customerNameKey, &v.CustomerName),
		)
	}

	return v.fieldMutator
}

var ProductHashKeySettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(customerNameKey, 123),
	),
	FieldOrder: fields.OrderedFieldKeys{customerNameKey},
}
