package logtable

import (
	"time"

	"github.com/sophielizg/go-libs/datastore"
)

type LogDataRow struct {
	Message     string
	LoggerName  string
	Level       string
	Stack       string
	Fields      map[string]interface{}
	CreatedTime time.Time
}

func (d *LogDataRow) GetFields() datastore.DataRowFields {
	return datastore.DataRowFields{
		"Message":     d.Message,
		"LoggerName":  d.LoggerName,
		"Level":       d.Level,
		"Stack":       d.Stack,
		"Fields":      d.Fields,
		"CreatedTime": d.CreatedTime,
	}
}

type logDataRowFactory struct {
	fieldTypes datastore.DataRowFieldTypes
}

func (f *logDataRowFactory) CreateDefault() *LogDataRow {
	return nil
}

func (f *logDataRowFactory) CreateFromFields(fields datastore.DataRowFields) (*LogDataRow, error) {
	return &LogDataRow{
		Message:     fields["Message"].(string),
		LoggerName:  fields["LoggerName"].(string),
		Level:       fields["Level"].(string),
		Stack:       fields["Stack"].(string),
		Fields:      fields["Fields"].(map[string]interface{}),
		CreatedTime: fields["CreatedTime"].(time.Time),
	}, nil
}

func (f *logDataRowFactory) GetFieldTypes() datastore.DataRowFieldTypes {
	if f.fieldTypes == nil {
		f.fieldTypes = datastore.DataRowFieldTypes{
			"Message":     &datastore.StringField{NumChars: 256},
			"LoggerName":  &datastore.StringField{NumChars: 64},
			"Level":       &datastore.StringField{NumChars: 8},
			"Stack":       &datastore.StringField{NumChars: 1024},
			"Fields":      &datastore.JsonField{},
			"CreatedTime": &datastore.TimeField{},
		}
	}

	return f.fieldTypes
}
