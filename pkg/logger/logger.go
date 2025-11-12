package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Set log level based on environment
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}

	Log.SetLevel(level)

	// Set formatter
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "YYYY-MM-DD HH:mm:ss",
		DisableColors:   false,
	})

	// Set output
	Log.SetOutput(os.Stdout)

	// Add custom fields
	Log.WithFields(logrus.Fields{
		"service": "whatsapp-bot",
		"version": "1.0.0",
	})
}

func WithFields(fields Fields) *logrus.Entry {
	return Log.WithFields(logrus.ToLogrusFields(fields))
}

type Fields map[string]interface{}

func ToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for key, value := range fields {
		logrusFields[key] = value
	}
	return logrusFields
}

func Debug(msg string) {
	Log.Debug(msg)
}

func Info(msg string) {
	Log.Info(msg)
}

func Warn(msg string) {
	Log.Warn(msg)
}

func Error(msg string) {
	Log.Error(msg)
}

func Fatal(msg string) {
	Log.Fatal(msg)
}

func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

func WithError(err error) *logrus.Entry {
	return Log.WithError(err)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

func WithFields(fields Fields) *logrus.Entry {
	return Log.WithFields(ToLogrusFields(fields))
}