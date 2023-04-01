package compare

type Comparator[T Comparable] struct {
	Op     Operator
	Values []T
}

func Eq[T Comparable](value T) *Comparator[T] {
	return &Comparator[T]{
		Op:     EQ,
		Values: []T{value},
	}
}

func Lt[T Comparable](value T) *Comparator[T] {
	return &Comparator[T]{
		Op:     LT,
		Values: []T{value},
	}
}

func Lte[T Comparable](value T) *Comparator[T] {
	return &Comparator[T]{
		Op:     LTE,
		Values: []T{value},
	}
}

func Gt[T Comparable](value T) *Comparator[T] {
	return &Comparator[T]{
		Op:     GT,
		Values: []T{value},
	}
}

func Gte[T Comparable](value T) *Comparator[T] {
	return &Comparator[T]{
		Op:     GTE,
		Values: []T{value},
	}
}

// between is inclusive
func Btw[T Comparable](value T) *Comparator[T] {
	return &Comparator[T]{
		Op:     BTW,
		Values: []T{value},
	}
}
