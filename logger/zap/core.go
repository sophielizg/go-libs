package zap

import (
	"encoding/json"
	"time"

	"github.com/sophielizg/go-libs/datastore"
	"github.com/sophielizg/go-libs/logger/logtable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logTableWriter struct {
	logTable *logtable.LogTable
}

func (w *logTableWriter) Write(bytes []byte) (int, error) {
	var entry map[string]interface{}
	err := json.Unmarshal(bytes, &entry)
	if err != nil {
		return 0, err
	}

	dataRow := &logtable.LogDataRow{}

	if val, ok := entry["Time"]; ok {
		dataRow.CreatedTime = val.(time.Time)
		delete(entry, "Time")
	}

	if val, ok := entry["Level"]; ok {
		dataRow.Level = val.(zapcore.Level).String()
		delete(entry, "Level")
	}

	if val, ok := entry["LoggerName"]; ok {
		dataRow.LoggerName = val.(string)
		delete(entry, "LoggerName")
	}

	if val, ok := entry["Message"]; ok {
		dataRow.Message = val.(string)
		delete(entry, "Message")
	}

	if val, ok := entry["Stack"]; ok {
		dataRow.Stack = val.(string)
		delete(entry, "Stack")
	}

	if _, ok := entry["Caller"]; ok {
		delete(entry, "Caller")
	}

	dataRow.Fields = entry
	w.logTable.Append(dataRow)

	return len(bytes), nil
}

func logTableCore(backend datastore.AppendTableBackend) zapcore.Core {
	tablePriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	logTable := &logTableWriter{
		logTable: logtable.CreateLogTable(backend),
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(logTable), tablePriority)
}
