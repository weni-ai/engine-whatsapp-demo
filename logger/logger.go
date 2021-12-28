package logger

import (
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
	"github.com/weni/whatsapp-router/config"
)

func init() {
	logrus.SetOutput(os.Stdout)
	level, err := logrus.ParseLevel("debug")
	if err != nil {
		logrus.Fatalf("Invalid log level '%s'", level)
	}
	logrus.SetLevel(level)

	if config.GetConfig().Server.SentryDSN != "" {
		hook, err := logrus_sentry.NewSentryHook(config.GetConfig().Server.SentryDSN, []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel})
		hook.Timeout = 0
		hook.StacktraceConfiguration.Enable = true
		hook.StacktraceConfiguration.Skip = 4
		hook.StacktraceConfiguration.Context = 5
		if err != nil {
			logrus.Fatalf("invalid sentry DSN: '%s': %s", config.GetConfig().Server.SentryDSN, err)
		}
		logrus.StandardLogger().Hooks.Add(hook)
	}
}

func Info(message string) {
	logrus.Info(message)
}

func Debug(message string) {
	logrus.Debug(message)
}

func Error(message string) {
	logrus.Error(message)
}
