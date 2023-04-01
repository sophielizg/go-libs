package purchase

import (
	"github.com/sophielizg/go-libs/datastore/compare"
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	purchaseTimeKey = "PurchaseTime"
	itemBrandKey    = "ItemBrand"
	itemNameKey     = "ItemName"
)

type PurchaseSortKeyComparator struct {
	PurchaseTime *compare.Comparator[fields.Time]
	ItemBrand    *compare.Comparator[fields.String]
	ItemName     *compare.Comparator[fields.String]
	fieldMutator *mutator.FieldMutator
}

func (c *PurchaseSortKeyComparator) Mutator() *mutator.FieldMutator {
	if c.fieldMutator == nil {
		c.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(purchaseTimeKey, &c.PurchaseTime),
			mutator.WithAddress(itemBrandKey, &c.ItemBrand),
			mutator.WithAddress(itemNameKey, &c.ItemName),
		)
	}

	return c.fieldMutator
}

type PurchaseSortKey struct {
	PurchaseTime fields.Time
	ItemBrand    fields.String
	ItemName     fields.String
	fieldMutator *mutator.FieldMutator
}

func (v *PurchaseSortKey) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(purchaseTimeKey, &v.PurchaseTime),
			mutator.WithAddress(itemBrandKey, &v.ItemBrand),
			mutator.WithAddress(itemNameKey, &v.ItemName),
		)
	}

	return v.fieldMutator
}

var PurchaseSortKeySettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(itemBrandKey, 63),
		fields.WithNumBytes(itemNameKey, 255),
	),
	FieldOrder: fields.OrderedFieldKeys{itemBrandKey, itemNameKey},
}
