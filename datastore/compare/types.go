package compare

import "github.com/sophielizg/go-libs/datastore/fields"

type Operator = int8

const (
	EQ Operator = iota
	LT
	LTE
	GT
	GTE
	BTW
)

type Comparable interface {
	fields.IntField | fields.FloatField | fields.TimeField | fields.BoolField | fields.StringField
}
