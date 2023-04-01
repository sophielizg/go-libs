package mutator

type MappedFieldValues = map[string]any

type Mutatable[M any] interface {
	*M
	Mutator() *FieldMutator
}
