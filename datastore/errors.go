package datastore

import "errors"

var InputLengthMismatchError = errors.New("the number of keys and values input must match")

var OutputLengthMismatchError = errors.New("the number of keys or values output must exactly match how mamy were input")
