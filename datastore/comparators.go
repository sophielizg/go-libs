package datastore

const (
	equalComparator Option = iota
	lessThanComparator
	greaterThanComparator
)

var Comparators = struct {
	Equal       Option
	LessThan    Option
	GreaterThan Option
}{
	equalComparator,
	lessThanComparator,
	greaterThanComparator,
}

func isComparator(x Option) bool {
	return equalComparator <= x && greaterThanComparator >= x
}
