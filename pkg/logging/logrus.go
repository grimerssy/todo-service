package logging

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type ConfigLogrus struct {
	Format string
}

type Logrus struct {
	logger *logrus.Logger
}

func DefaultLogrus() *Logrus {
	return &Logrus{
		logger: logrus.StandardLogger(),
	}
}

func NewLogrus(cfg ConfigLogrus) *Logrus {
	logger := logrus.New()

	switch strings.ToUpper(cfg.Format) {
	case "JSON":
		logger.Formatter = &logrus.JSONFormatter{}
	case "TEXT":
		logger.Formatter = &logrus.TextFormatter{}
	default:
		logger.Fatal("invalid logger format")
	}

	return &Logrus{
		logger: logger,
	}
}

func (l *Logrus) Log(level Level, args ...any) {
	l.logger.Log(logrus.Level(level), args...)
}

func (l *Logrus) Logln(level Level, args ...any) {
	l.logger.Logln(logrus.Level(level), args...)
}

func (l *Logrus) Logf(level Level, format string, args ...any) {
	l.logger.Logf(logrus.Level(level), format, args...)
}

func (l *Logrus) LogFields(level Level, fields Fields, args ...any) {
	l.logger.WithFields(logrus.Fields(fields)).Log(logrus.Level(level), args...)
}

func (l *Logrus) LogFieldsln(level Level, fields Fields, args ...any) {
	l.logger.WithFields(logrus.Fields(fields)).Logln(logrus.Level(level), args...)
}

func (l *Logrus) LogFieldsf(level Level, fields Fields, format string, args ...any) {
	l.logger.WithFields(logrus.Fields(fields)).Logf(logrus.Level(level), format, args...)
}
