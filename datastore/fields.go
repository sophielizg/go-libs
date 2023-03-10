package datastore

import "time"

const (
	intFieldId Option = iota
	stringFieldId
	jsonFieldId
	boolFieldId
	timeFieldId
	floatFieldId
)

type FieldType interface {
	TypeId() Option
	IsNullable() bool
	IsComparable() bool
	IsType(x interface{}) bool
}

type IntField struct {
	Nullable bool
	Unsigned bool
	Large    bool
}

func (f *IntField) TypeId() Option {
	return intFieldId
}

func (f *IntField) IsNullable() bool {
	return f.Nullable
}

func (f *IntField) IsComparable() bool {
	return true
}

func (f *IntField) IsType(x interface{}) bool {
	isType := false

	switch true {
	case f.Nullable && f.Unsigned && f.Large:
		_, isType = x.(*uint64)
	case f.Nullable && f.Unsigned && !f.Large:
		_, isType = x.(*uint)
	case f.Nullable && !f.Unsigned && f.Large:
		_, isType = x.(*int64)
	case f.Nullable && !f.Unsigned && !f.Large:
		_, isType = x.(*int)
	case !f.Nullable && f.Unsigned && f.Large:
		_, isType = x.(uint64)
	case !f.Nullable && f.Unsigned && !f.Large:
		_, isType = x.(uint)
	case !f.Nullable && !f.Unsigned && f.Large:
		_, isType = x.(int64)
	case !f.Nullable && !f.Unsigned && !f.Large:
		_, isType = x.(int)
	}

	return isType
}

type StringField struct {
	Nullable bool
	NumChars int
}

func (f *StringField) TypeId() Option {
	return stringFieldId
}

func (f *StringField) IsNullable() bool {
	return f.Nullable
}

func (f *StringField) IsComparable() bool {
	return true
}

func (f *StringField) IsType(x interface{}) bool {
	var ok bool
	var val string
	var valPtr *string

	if f.Nullable {
		valPtr, ok = x.(*string)
	} else {
		val, ok = x.(string)
	}

	if !ok {
		return false
	}

	if f.Nullable {
		if valPtr == nil {
			return true
		}
		val = *valPtr
	}

	if len(val) > f.NumChars {
		return false
	}

	return true
}

type JsonField struct {
	Nullable bool
}

func (f *JsonField) TypeId() Option {
	return jsonFieldId
}

func (f *JsonField) IsNullable() bool {
	return f.Nullable
}

func (f *JsonField) IsComparable() bool {
	return false
}

func (f *JsonField) IsType(x interface{}) bool {
	if f.Nullable && x == nil {
		return true
	}

	_, isMap := x.(map[string]interface{})
	if isMap {
		return true
	}

	_, isList := x.([]interface{})
	if isList {
		return true
	}

	return false
}

type BoolField struct {
	Nullable bool
}

func (f *BoolField) TypeId() Option {
	return boolFieldId
}

func (f *BoolField) IsNullable() bool {
	return f.Nullable
}

func (f *BoolField) IsComparable() bool {
	return true
}

func (f *BoolField) IsType(x interface{}) bool {
	var isType bool

	if f.Nullable {
		_, isType = x.(*bool)
	} else {
		_, isType = x.(bool)
	}

	return isType
}

type TimeField struct {
	Nullable bool
}

func (f *TimeField) TypeId() Option {
	return timeFieldId
}

func (f *TimeField) IsNullable() bool {
	return f.Nullable
}

func (f *TimeField) IsComparable() bool {
	return true
}

func (f *TimeField) IsType(x interface{}) bool {
	var isType bool

	if f.Nullable {
		_, isType = x.(*time.Time)
	} else {
		_, isType = x.(time.Time)
	}

	return isType
}

type FloatField struct {
	Nullable bool
	Large    bool
}

func (f *FloatField) TypeId() Option {
	return floatFieldId
}

func (f *FloatField) IsNullable() bool {
	return f.Nullable
}

func (f *FloatField) IsComparable() bool {
	return true
}

func (f *FloatField) IsType(x interface{}) bool {
	var isType bool

	switch true {
	case f.Nullable && f.Large:
		_, isType = x.(*float64)
	case f.Nullable && !f.Large:
		_, isType = x.(*float32)
	case !f.Nullable && f.Large:
		_, isType = x.(float64)
	case !f.Nullable && !f.Large:
		_, isType = x.(float32)
	}

	return isType
}
