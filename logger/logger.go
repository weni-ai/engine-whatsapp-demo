package logger

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/weni/whatsapp-router/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	var err error

	if config.GetConfig().Server.SentryDSN != "" {
		err = sentry.Init(sentry.ClientOptions{
			Dsn:              config.GetConfig().Server.SentryDSN,
			AttachStacktrace: true,
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("sentry.Init: %s", err))
		}
	}

	zconfig := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	zconfig.EncoderConfig = encoderConfig

	log, err = zconfig.Build(
		zap.AddCallerSkip(1),
		zap.Hooks(func(entry zapcore.Entry) error {
			if entry.Level == zapcore.ErrorLevel {
				defer sentry.Flush(2 * time.Second)
				sentry.CaptureMessage(fmt.Sprintf("%s, Line No: %d :: %s", entry.Caller.File, entry.Caller.Line, entry.Message))
			}
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}
}

func Info(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	log.Error(message, fields...)
}
