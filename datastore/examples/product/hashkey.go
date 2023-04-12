package product

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	brandKey = "Brand"
	nameKey  = "Name"
)

type HashKey struct {
	Brand        fields.String
	Name         fields.String
	fieldMutator *mutator.FieldMutator
}

func (v *HashKey) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(brandKey, &v.Brand),
			mutator.WithAddress(nameKey, &v.Name),
		)
	}

	return v.fieldMutator
}

var HashKeySettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(brandKey, 63),
		fields.WithNumBytes(nameKey, 255),
	),
	FieldOrder: fields.OrderedFieldKeys{brandKey, nameKey},
}
