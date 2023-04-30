package purchase

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	CustomerNameKey = "CustomerName"
	PurchaseTimeKey = "PurchaseTime"
	ItemBrandKey    = "ItemBrand"
	ItemNameKey     = "ItemName"
)

type Key struct {
	CustomerName fields.String
	PurchaseTime fields.Time
	ItemBrand    fields.String
	ItemName     fields.String
}

func (k *Key) Mutator() *mutator.FieldMutator {
	return mutator.NewFieldMutator(
		mutator.WithAddress(CustomerNameKey, &k.CustomerName),
		mutator.WithAddress(PurchaseTimeKey, &k.PurchaseTime),
		mutator.WithAddress(ItemBrandKey, &k.ItemBrand),
		mutator.WithAddress(ItemNameKey, &k.ItemName),
	)
}

var KeySettings = &fields.RowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(CustomerNameKey, 123),
		fields.WithNumBytes(ItemBrandKey, 63),
		fields.WithNumBytes(ItemNameKey, 255),
	),
	FieldOrder: fields.OrderedFieldKeys{CustomerNameKey, PurchaseTimeKey, ItemBrandKey, ItemNameKey},
}
