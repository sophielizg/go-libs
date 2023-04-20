package inmemory

import "errors"

var AutoGenerateNotSupportedError = errors.New("auto generate fields are not supported for inmemory backends")

var HashKeyExistsError = errors.New("cannot add a hash key that already exists")

var HashKeyDoesNotExistError = errors.New("cannot update a hash key that does not already exist")
