package datastore

// TODO: revisit comparator system

const (
	EqualComparator Option = iota
	LessThanComparator
	GreaterThanComparator
)

func isComparator(x Option) bool {
	return EqualComparator <= x && GreaterThanComparator >= x
}
