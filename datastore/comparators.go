package datastore

var (
	Comparators = struct {
		EqualComparator       string
		LessThanComparator    string
		GreaterThanComparator string
	}{
		"Equal",
		"LessThan",
		"GreaterThan",
	}

	ComparatorTypes = OptionTypes{
		Comparators.EqualComparator:       true,
		Comparators.LessThanComparator:    true,
		Comparators.GreaterThanComparator: true,
	}
)
