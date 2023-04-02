package logtable

import (
	"github.com/sophielizg/go-libs/datastore/fields"
	"github.com/sophielizg/go-libs/datastore/mutator"
)

const (
	messageKey     = "Message"
	sourceKey      = "Source"
	levelKey       = "Level"
	createdTimeKey = "CreatedTime"
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
			mutator.WithAddress(messageKey, &v.Message),
			mutator.WithAddress(sourceKey, &v.Source),
			mutator.WithAddress(levelKey, &v.Level),
			mutator.WithAddress(createdTimeKey, &v.CreatedTime),
		)
	}

	return v.fieldMutator
}

var LogDataRowSettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(messageKey, 255),
		fields.WithNumBytes(sourceKey, 63),
		fields.WithNumBytes(levelKey, 7),
	),
	FieldOrder: fields.OrderedFieldKeys{sourceKey, levelKey, createdTimeKey, messageKey},
}
