package compare

type Comparator[T Comparable] struct {
	Op     Operator
	Values []T
}

func (c1 *Comparator[T]) Equals(c2 *Comparator[T]) bool {
	if c1.Op != c2.Op {
		return false
	}

	if c1.Values == nil && c2.Values == nil {
		return true
	} else if c1.Values == nil || c2.Values == nil {
		return false
	} else if len(c1.Values) != len(c2.Values) {
		return false
	}

	for i := range c1.Values {
		if c1.Values[i] != c2.Values[i] {
			return false
		}
	}

	return true
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
func Btw[T Comparable](lower, upper T) *Comparator[T] {
	return &Comparator[T]{
		Op:     BTW,
		Values: []T{lower, upper},
	}
}
