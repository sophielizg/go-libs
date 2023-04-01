package fields

import (
	"time"
)

type FieldType interface {
	IntField | FloatField | StringField | JsonField | BoolField | TimeField
}

type IntField interface {
	Int | NullInt | UInt | NullUInt
	BigInt | NullBigInt | BigUInt | NullBigUInt
}

type Int = int
type NullInt = *int
type UInt = uint
type NullUInt = *uint
type BigInt = int64
type NullBigInt = *int64
type BigUInt = uint64
type NullBigUInt = *uint64

type FloatField interface {
	Float | NullFloat | BigFloat | NullBigFloat
}

type Float = float32
type NullFloat = *float32
type BigFloat = float64
type NullBigFloat = *float64

type StringField interface {
	String | NullString
}

type String = string
type NullString = *string

type JsonField interface {
	JsonMap | JsonList
}

type JsonMap = map[string]any
type JsonList = []any

type BoolField interface {
	Bool | NullBool
}

type Bool = bool
type NullBool = *bool

type TimeField interface {
	Time | NullTime
}

type Time = time.Time
type NullTime = *time.Time
