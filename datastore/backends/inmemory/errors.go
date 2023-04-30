package inmemory

import "errors"

var AutoGenerateNotSupportedError = errors.New("auto generate fields are not supported for inmemory backends")

var KeyExistsError = errors.New("cannot add a key that already exists")

var KeyDoesNotExistError = errors.New("cannot update a key that does not already exist")
