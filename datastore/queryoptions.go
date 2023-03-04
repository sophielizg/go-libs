package datastore

type WriteOption struct {
	GenerateCurrentTime bool
}

type WriteOptions map[string]WriteOption

type UpdateOption struct {
	WriteOption
	Include bool
	Exclude bool
}

type UpdateOptions map[string]UpdateOption
