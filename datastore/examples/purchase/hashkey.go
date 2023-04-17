package purchase

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	CustomerNameKey = "CustomerName"
)

type HashKey struct {
	CustomerName fields.String
	fieldMutator *mutator.FieldMutator
}

func (v *HashKey) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(CustomerNameKey, &v.CustomerName),
		)
	}

	return v.fieldMutator
}

var HashKeySettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(CustomerNameKey, 123),
	),
	FieldOrder: fields.OrderedFieldKeys{CustomerNameKey},
}
