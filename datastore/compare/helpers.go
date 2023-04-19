package compare

import "github.com/sophielizg/go-libs/datastore/fields"

func IsNilComparator(comparator any) bool {
	switch val := comparator.(type) {
	case *Comparator[fields.Int]:
		return val == nil
	case *Comparator[fields.NullInt]:
		return val == nil
	case *Comparator[fields.UInt]:
		return val == nil
	case *Comparator[fields.NullUInt]:
		return val == nil
	case *Comparator[fields.BigInt]:
		return val == nil
	case *Comparator[fields.NullBigInt]:
		return val == nil
	case *Comparator[fields.BigUInt]:
		return val == nil
	case *Comparator[fields.NullBigUInt]:
		return val == nil
	case *Comparator[fields.SmallFloat]:
		return val == nil
	case *Comparator[fields.NullSmallFloat]:
		return val == nil
	case *Comparator[fields.Float]:
		return val == nil
	case *Comparator[fields.NullFloat]:
		return val == nil
	case *Comparator[fields.String]:
		return val == nil
	case *Comparator[fields.NullString]:
		return val == nil
	case *Comparator[fields.Bool]:
		return val == nil
	case *Comparator[fields.NullBool]:
		return val == nil
	case *Comparator[fields.Time]:
		return val == nil
	case *Comparator[fields.NullTime]:
		return val == nil
	default:
		return false
	}
}
