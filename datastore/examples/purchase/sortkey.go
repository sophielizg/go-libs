package purchase

import (
	"github.com/sophielizg/go-libs/datastore/compare"
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	PurchaseTimeKey = "PurchaseTime"
	ItemBrandKey    = "ItemBrand"
	ItemNameKey     = "ItemName"
)

type SortKey struct {
	PurchaseTime fields.Time
	ItemBrand    fields.String
	ItemName     fields.String
	fieldMutator *mutator.FieldMutator
}

func (v *SortKey) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(PurchaseTimeKey, &v.PurchaseTime),
			mutator.WithAddress(ItemBrandKey, &v.ItemBrand),
			mutator.WithAddress(ItemNameKey, &v.ItemName),
		)
	}

	return v.fieldMutator
}

type SortKeyComparator struct {
	PurchaseTime *compare.Comparator[fields.Time]
	ItemBrand    *compare.Comparator[fields.String]
	ItemName     *compare.Comparator[fields.String]
	fieldMutator *mutator.FieldMutator
}

func (c *SortKeyComparator) Mutator() *mutator.FieldMutator {
	if c.fieldMutator == nil {
		c.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(PurchaseTimeKey, &c.PurchaseTime),
			mutator.WithAddress(ItemBrandKey, &c.ItemBrand),
			mutator.WithAddress(ItemNameKey, &c.ItemName),
		)
	}

	return c.fieldMutator
}

var SortKeySettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(ItemBrandKey, 63),
		fields.WithNumBytes(ItemNameKey, 255),
	),
	FieldOrder: fields.OrderedFieldKeys{ItemBrandKey, ItemNameKey},
}
