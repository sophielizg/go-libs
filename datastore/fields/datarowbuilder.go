package fields

import (
	"errors"
)

var SetFieldTypeError = errors.New("unable to set field: invalid type")

type DataRowBuilder struct {
	fieldGetters map[string]func() any
	fieldSetters map[string]func(any) error
}

func (b *DataRowBuilder) GetField(key string) any {
	return b.fieldGetters[key]()
}

func (b *DataRowBuilder) GetFields() MappedFieldValues {
	fields := MappedFieldValues{}

	for key, getter := range b.fieldGetters {
		fields[key] = getter()
	}

	return fields
}

func (b *DataRowBuilder) SetField(key string, value any) error {
	return b.fieldSetters[key](value)
}

func (b *DataRowBuilder) SetFields(fields MappedFieldValues) error {
	for key, value := range fields {
		if err := b.fieldSetters[key](value); err != nil {
			return err
		}
	}

	return nil
}

func NewDataRowBuilder(options ...func(*DataRowBuilder)) *DataRowBuilder {
	builder := &DataRowBuilder{}

	for _, option := range options {
		option(builder)
	}

	return builder
}

func WithAddress[T FieldType](key string, address *T) func(b *DataRowBuilder) {
	return func(b *DataRowBuilder) {
		b.fieldSetters[key] = func(value any) error {
			if v, ok := value.(T); ok {
				*address = v
				return nil
			}

			return SetFieldTypeError
		}

		b.fieldGetters[key] = func() any {
			return *address
		}
	}
}
