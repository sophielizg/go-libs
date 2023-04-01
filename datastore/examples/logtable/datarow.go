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

type A struct {
	*LogDataRow
}

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
		fields.WithNumBytes(sourceKey, 255),
		fields.WithNumBytes(levelKey, 255),
	),
	FieldOrder: fields.OrderedFieldKeys{messageKey, sourceKey, levelKey, createdTimeKey},
}
