package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
}

func Init(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		log.Warnf("Invalid log level %s, defaulting to info", level)
		lvl = logrus.InfoLevel
	}
	log.SetLevel(lvl)
}

func Debug(msg string, fields ...interface{}) {
	log.WithFields(parseFields(fields...)).Debug(msg)
}

func Info(msg string, fields ...interface{}) {
	log.WithFields(parseFields(fields...)).Info(msg)
}

func Warn(msg string, fields ...interface{}) {
	log.WithFields(parseFields(fields...)).Warn(msg)
}

func Error(msg string, fields ...interface{}) {
	log.WithFields(parseFields(fields...)).Error(msg)
}

func Fatal(msg string, fields ...interface{}) {
	log.WithFields(parseFields(fields...)).Fatal(msg)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}

func parseFields(fields ...interface{}) logrus.Fields {
	f := make(logrus.Fields)
	if len(fields)%2 != 0 {
		return f
	}
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		f[key] = fields[i+1]
	}
	return f
}