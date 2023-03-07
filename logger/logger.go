package logger

var logger Logger

type LogFields map[string]interface{}

type Logger interface {
	WithFields(fields LogFields) Logger
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

func SetLogger(l Logger) {
	logger = l
}

func WithFields(fields LogFields) Logger {
	return logger.WithFields(fields)
}

func Debug(msg string) {
	logger.Debug(msg)
}

func Info(msg string) {
	logger.Info(msg)
}

func Warn(msg string) {
	logger.Warn(msg)
}

func Error(msg string) {
	logger.Error(msg)
}

func Fatal(msg string) {
	logger.Fatal(msg)
}
