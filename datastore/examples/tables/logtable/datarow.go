package logtable

import (
	"time"

	"github.com/sophielizg/go-libs/datastore"
)

type LogDataRow struct {
	Message     string
	Source      string
	Level       string
	CreatedTime time.Time
}

func (d *LogDataRow) GetFields() datastore.DataRowFields {
	return datastore.DataRowFields{
		"Message":     d.Message,
		"Source":      d.Source,
		"Level":       d.Level,
		"CreatedTime": d.CreatedTime,
	}
}

type LogDataRowFactory struct{}

func (f *LogDataRowFactory) CreateDefault() *LogDataRow {
	return &LogDataRow{}
}

func (f *LogDataRowFactory) CreateFromFields(fields datastore.DataRowFields) (*LogDataRow, error) {
	return &LogDataRow{
		Message:     fields["Message"].(string),
		Source:      fields["Source"].(string),
		Level:       fields["Level"].(string),
		CreatedTime: fields["CreatedTime"].(time.Time),
	}, nil
}

func (f *LogDataRowFactory) GetFieldTypes() datastore.DataRowFieldTypes {
	return datastore.DataRowFieldTypes{
		"Message":     &datastore.StringField{NumChars: 1024},
		"Source":      &datastore.StringField{NumChars: 64},
		"Level":       &datastore.StringField{NumChars: 8},
		"CreatedTime": &datastore.TimeField{},
	}
}
