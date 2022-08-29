package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log ILogger

type ILogger interface {
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
}

type (
	Field  zap.Field
	logger struct {
		zaplog *zap.Logger
	}
)

func InitDefaultLogger(logfile string, debug bool, initialFields map[string]interface{}) {
	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}

	logConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{ //or zap.NewProductionEncoderConfig()
			TimeKey:       "ts",
			LevelKey:      "level",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", logfile},
		ErrorOutputPaths: []string{"stderr", logfile},
		InitialFields:    initialFields,
	}

	log, err := logConfig.Build()
	if err != nil {
		panic(err)
	}

	Log = &logger{log}
}

func (l *logger) Info(msg string, fields ...Field) {
	l.zaplog.Info(msg, withFields(fields...)...)
	defer l.zaplog.Sync()
}

func (l *logger) Warn(msg string, fields ...Field) {
	l.zaplog.With(withFields(fields...)...).Warn(msg)
	defer l.zaplog.Sync()
}

func (l *logger) Error(msg string, fields ...Field) {
	l.zaplog.With(withFields(fields...)...).Error(msg)
	defer l.zaplog.Sync()
}

func (l *logger) Debug(msg string, fields ...Field) {
	l.zaplog.With(withFields(fields...)...).Debug(msg)
	defer l.zaplog.Sync()
}

func withFields(fields ...Field) []zap.Field {
	zapFields := make([]zap.Field, 0)
	for _, f := range fields {
		zapFields = append(zapFields, zap.Field(f))
	}
	return zapFields
}

func String(key string, val string) Field {
	return Field(zap.String(key, val))
}
