package zap

import (
	"github.com/sophielizg/go-libs/logger"
	"go.uber.org/zap"
)

type Logger struct {
	sugar *zap.SugaredLogger
}

func (l *Logger) WithFields(fields logger.LogFields) logger.Logger {
	mergedFields := make([]interface{}, len(fields)*2)
	i := 0
	for key, value := range fields {
		mergedFields[i] = key
		mergedFields[i+1] = value
		i += 2
	}

	return &Logger{
		sugar: l.sugar.With(mergedFields...),
	}
}

func (l *Logger) Debug(msg string) {
	l.sugar.Debug(msg)
}

func (l *Logger) Info(msg string) {
	l.sugar.Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.sugar.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.sugar.Error(msg)
}

func (l *Logger) Fatal(msg string) {
	l.sugar.Fatal(msg)
}

func (l *Logger) Close() error {
	return l.sugar.Sync()
}
