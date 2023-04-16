package logtable

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	MessageKey     = "Message"
	SourceKey      = "Source"
	LevelKey       = "Level"
	CreatedTimeKey = "CreatedTime"
)

type LogDataRow struct {
	Message      fields.String
	Source       fields.String
	Level        fields.String
	CreatedTime  fields.Time
	fieldMutator *mutator.FieldMutator
}

func (v *LogDataRow) Mutator() *mutator.FieldMutator {
	if v.fieldMutator == nil {
		v.fieldMutator = mutator.NewFieldMutator(
			mutator.WithAddress(MessageKey, &v.Message),
			mutator.WithAddress(SourceKey, &v.Source),
			mutator.WithAddress(LevelKey, &v.Level),
			mutator.WithAddress(CreatedTimeKey, &v.CreatedTime),
		)
	}

	return v.fieldMutator
}

var LogDataRowSettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(MessageKey, 255),
		fields.WithNumBytes(SourceKey, 63),
		fields.WithNumBytes(LevelKey, 7),
	),
	FieldOrder: fields.OrderedFieldKeys{SourceKey, LevelKey, CreatedTimeKey, MessageKey},
}
