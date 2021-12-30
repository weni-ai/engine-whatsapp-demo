package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"github.com/weni/whatsapp-router/config"
)

func init() {
	logrus.SetOutput(os.Stdout)
	level, err := logrus.ParseLevel(config.GetConfig().Server.LogLevel)
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

func MiddlewareLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var requestID string
		if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
			requestID = reqID.(string)
		}
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		latency := time.Since(start)

		fields := logrus.Fields{
			"status":                              ww.Status(),
			"took":                                latency,
			fmt.Sprintf("measure#%s.latency", ""): latency.Nanoseconds(),
			"remote":                              r.RemoteAddr,
			"request":                             r.RequestURI,
			"method":                              r.Method,
		}
		if requestID != "" {
			fields["request-id"] = requestID
		}
		logrus.WithFields(fields).Info("request completed")
	}
	return http.HandlerFunc(fn)
}
