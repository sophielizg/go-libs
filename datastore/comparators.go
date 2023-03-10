package datastore

const (
	EqualComparator Option = iota
	LessThanComparator
	GreaterThanComparator
)

func isComparator(x Option) bool {
	return EqualComparator <= x && GreaterThanComparator >= x
}
