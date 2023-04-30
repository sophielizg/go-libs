package product

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	BrandKey = "Brand"
	NameKey  = "Name"
)

type Key struct {
	Brand fields.String
	Name  fields.String
}

func (k *Key) Mutator() *mutator.FieldMutator {
	return mutator.NewFieldMutator(
		mutator.WithAddress(BrandKey, &k.Brand),
		mutator.WithAddress(NameKey, &k.Name),
	)
}

var KeySettings = &fields.RowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(BrandKey, 63),
		fields.WithNumBytes(NameKey, 255),
	),
	FieldOrder: fields.OrderedFieldKeys{BrandKey, NameKey},
}
