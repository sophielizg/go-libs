package zap

import (
	"os"

	"github.com/sophielizg/go-libs/datastore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateWithBackends(name string, backends ...datastore.AppendTableBackend) *Logger {
	cores := make([]zapcore.Core, len(backends))

	for i, backend := range backends {
		cores[i] = logTableCore(backend)
	}

	return CreateWithCores(name, cores...)
}

func CreateWithCores(name string, cores ...zapcore.Core) *Logger {
	consoleHighPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	consoleLowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	cores = append(cores, zapcore.NewCore(consoleEncoder, consoleErrors, consoleHighPriority))
	cores = append(cores, zapcore.NewCore(consoleEncoder, consoleDebugging, consoleLowPriority))

	core := zapcore.NewTee(cores...)

	logger := zap.New(core).Named(name)
	return &Logger{
		sugar: logger.Sugar(),
	}
}
