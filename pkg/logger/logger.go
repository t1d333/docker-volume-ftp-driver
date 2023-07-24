package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Infow(msg string, args ...interface{})
	Warnw(msg string, args ...interface{})
	Errorw(msg string, args ...interface{})
	Fatalw(msg string, args ...interface{})

	Infof(s string, args ...interface{})
	Warnf(s string, args ...interface{})
	Errorf(s string, args ...interface{})
	Fatalf(s string, args ...interface{})
}

func NewLogger() Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000")
	config.EncoderConfig.StacktraceKey = ""

	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}

	log := logger.Sugar()

	return log
}
