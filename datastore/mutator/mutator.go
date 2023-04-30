package mutator

import (
	"errors"

	"github.com/sophielizg/go-libs/utils"
)

var SetFieldTypeError = errors.New("unable to set field: invalid type")

type gettersMap = map[string]func() any
type settersMap = map[string]func(any) error

type FieldMutator struct {
	fieldGetters gettersMap
	fieldSetters settersMap
}

func (m *FieldMutator) GetField(key string) any {
	return m.fieldGetters[key]()
}

func (m *FieldMutator) GetFields() MappedFieldValues {
	fields := MappedFieldValues{}

	for key, getter := range m.fieldGetters {
		fields[key] = getter()
	}

	return fields
}

func (m *FieldMutator) SetField(key string, value any) error {
	return m.fieldSetters[key](value)
}

func (m *FieldMutator) SetFields(fields MappedFieldValues) error {
	for key, value := range fields {
		if err := m.SetField(key, value); err != nil {
			return err
		}
	}

	return nil
}

func NewFieldMutator(options ...func(*FieldMutator)) *FieldMutator {
	builder := &FieldMutator{
		fieldGetters: gettersMap{},
		fieldSetters: settersMap{},
	}

	for _, option := range options {
		option(builder)
	}

	return builder
}

func WithAddress[T any](key string, address *T) func(m *FieldMutator) {
	return func(m *FieldMutator) {
		m.fieldSetters[key] = func(value any) error {
			if v, ok := value.(T); ok {
				*address = v
				return nil
			}

			return SetFieldTypeError
		}

		m.fieldGetters[key] = func() any {
			return *address
		}
	}
}

func MergeFieldMutators(mutators ...*FieldMutator) *FieldMutator {
	merged := NewFieldMutator()

	for _, mutator := range mutators {
		merged.fieldGetters = utils.MergeMaps(merged.fieldGetters, mutator.fieldGetters)
		merged.fieldSetters = utils.MergeMaps(merged.fieldSetters, mutator.fieldSetters)
	}

	return merged
}
