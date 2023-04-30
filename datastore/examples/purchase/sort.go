package purchase

import (
	"github.com/sophielizg/go-libs/datastore/compare"
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

type SortComparator struct {
	PurchaseTime *compare.Comparator[fields.Time]
	ItemBrand    *compare.Comparator[fields.String]
	ItemName     *compare.Comparator[fields.String]
}

func (c *SortComparator) Mutator() *mutator.FieldMutator {
	return mutator.NewFieldMutator(
		mutator.WithAddress(PurchaseTimeKey, &c.PurchaseTime),
		mutator.WithAddress(ItemBrandKey, &c.ItemBrand),
		mutator.WithAddress(ItemNameKey, &c.ItemName),
	)
}

var SortFieldNames = fields.SortFieldNames{PurchaseTimeKey, ItemBrandKey, ItemNameKey}
