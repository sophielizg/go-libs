package product

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	BrandKey = "Brand"
	NameKey  = "Name"
)

type HashKey struct {
	Brand        fields.String
	Name         fields.String
	fieldMutator *mutator.FieldMutator
}

func (v *HashKey) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(BrandKey, &v.Brand),
			mutator.WithAddress(NameKey, &v.Name),
		)
	}

	return v.fieldMutator
}

var HashKeySettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(BrandKey, 63),
		fields.WithNumBytes(NameKey, 255),
	),
	FieldOrder: fields.OrderedFieldKeys{BrandKey, NameKey},
}
