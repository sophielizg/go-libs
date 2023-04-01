package product

import (
	"github.com/sophielizg/go-libs/datastore/fields"
)

const (
	brandKey = "Brand"
	nameKey  = "Name"
)

type ProductHashKey struct {
	Brand   fields.String
	Name    fields.String
	builder *fields.DataRowBuilder
}

func (v *ProductHashKey) Builder() *fields.DataRowBuilder {
	if v.builder == nil {
		v.builder = fields.NewDataRowBuilder(
			fields.WithAddress(brandKey, &v.Brand),
			fields.WithAddress(nameKey, &v.Name),
		)
	}

	return v.builder
}

var ProductHashKeySettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(brandKey, 63),
		fields.WithNumBytes(nameKey, 255),
	),
	FieldOrder: fields.OrderedFieldKeys{brandKey, nameKey},
}
