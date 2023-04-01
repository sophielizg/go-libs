package mutator

import (
	"errors"
)

var SetFieldTypeError = errors.New("unable to set field: invalid type")

type FieldMutator struct {
	fieldGetters map[string]func() any
	fieldSetters map[string]func(any) error
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
		if err := m.fieldSetters[key](value); err != nil {
			return err
		}
	}

	return nil
}

func NewFieldMutator(options ...func(*FieldMutator)) *FieldMutator {
	builder := &FieldMutator{}

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
