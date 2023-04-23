package logtable

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	MessageKey     = "Message"
	LoggerNameKey  = "LoggerName"
	LevelKey       = "Level"
	StackKey       = "Stack"
	FieldsKey      = "Fields"
	CreatedTimeKey = "CreatedTime"
)

type DataRow struct {
	Message      fields.String
	LoggerName   fields.String
	Level        fields.String
	Stack        fields.String
	Fields       fields.JsonMap
	CreatedTime  fields.Time
	fieldMutator *mutator.FieldMutator
}

func (v *DataRow) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(MessageKey, &v.Message),
			mutator.WithAddress(LoggerNameKey, &v.LoggerName),
			mutator.WithAddress(LevelKey, &v.Level),
			mutator.WithAddress(StackKey, &v.Stack),
			mutator.WithAddress(FieldsKey, &v.Fields),
			mutator.WithAddress(CreatedTimeKey, &v.CreatedTime),
		)
	}

	return v.fieldMutator
}

var DataRowSettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(MessageKey, 255),
		fields.WithNumBytes(LoggerNameKey, 63),
		fields.WithNumBytes(LevelKey, 7),
		fields.WithNumBytes(StackKey, 1023),
	),
	FieldOrder: fields.OrderedFieldKeys{LoggerNameKey, LevelKey, CreatedTimeKey, MessageKey, StackKey, FieldsKey},
}
