package queries

import "errors"

var ComparatorMissingFieldsError = errors.New("all SortKey fields on the left side must be included in comparator")
