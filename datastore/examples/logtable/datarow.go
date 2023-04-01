package logtable

import (
	"github.com/sophielizg/go-libs/datastore/fields"
)

const (
	messageKey     = "MessageKey"
	sourceKey      = "Source"
	levelKey       = "Level"
	createdTimeKey = "CreatedTime"
)

type LogDataRow struct {
	Message     fields.String
	Source      fields.String
	Level       fields.String
	CreatedTime fields.Time
	builder     *fields.DataRowBuilder
}

func (v *LogDataRow) Builder() *fields.DataRowBuilder {
	if v.builder == nil {
		v.builder = fields.NewDataRowBuilder(
			fields.WithAddress(messageKey, &v.Message),
			fields.WithAddress(sourceKey, &v.Source),
			fields.WithAddress(levelKey, &v.Level),
			fields.WithAddress(createdTimeKey, &v.CreatedTime),
		)
	}

	return v.builder
}

var LogDataRowSettings = fields.DataRowSettings{
	FieldSettings: fields.NewFieldSettings(
		fields.WithNumBytes(messageKey, 255),
		fields.WithNumBytes(sourceKey, 63),
		fields.WithNumBytes(levelKey, 7),
	),
	FieldOrder: fields.OrderedFieldKeys{sourceKey, levelKey, createdTimeKey, messageKey},
}
