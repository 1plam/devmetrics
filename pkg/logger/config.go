package logger

import (
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level            zapcore.Level
	Development      bool
	Encoding         string
	OutputPaths      []string
	ErrorOutputPaths []string
	EncoderConfig    zapcore.EncoderConfig
}

func DefaultConfig() Config {
	return Config{
		Level:       zapcore.InfoLevel,
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}
