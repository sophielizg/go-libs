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
}
