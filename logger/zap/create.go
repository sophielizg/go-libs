package zap

import (
	"os"

	"github.com/sophielizg/go-libs/logger/logtable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapConfig struct {
	Name  string
	Cores []zapcore.Core
}

func Create(options ...func(*zapConfig)) *Logger {
	config := &zapConfig{}

	for _, option := range options {
		option(config)
	}

	consoleHighPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	consoleLowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	config.Cores = append(config.Cores, zapcore.NewCore(consoleEncoder, consoleErrors, consoleHighPriority))
	config.Cores = append(config.Cores, zapcore.NewCore(consoleEncoder, consoleDebugging, consoleLowPriority))

	core := zapcore.NewTee(config.Cores...)

	logger := zap.New(core).Named(config.Name)
	return &Logger{
		sugar: logger.Sugar(),
	}
}

func WithName(name string) func(config *zapConfig) {
	return func(config *zapConfig) {
		config.Name = name
	}
}

func WithLogTable(table *logtable.LogTable) func(config *zapConfig) {
	return func(config *zapConfig) {
		core := logTableCore(table)
		config.Cores = append(config.Cores, core)
	}
}

func WithCore(core zapcore.Core) func(config *zapConfig) {
	return func(config *zapConfig) {
		config.Cores = append(config.Cores, core)
	}
}
