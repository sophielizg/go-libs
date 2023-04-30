package product

import "github.com/sophielizg/go-libs/datastore/fields"

type Entry = fields.KeyedEntry[Key, *Key, Data, *Data]
