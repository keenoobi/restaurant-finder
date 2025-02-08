package logger

import (
	"Go_Day03/internal/interfaces/logger" // Импортируем пакет с интерфейсом

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func NewLogrusLogger() *LogrusLogger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetLevel(logrus.InfoLevel)
	return &LogrusLogger{logger: l}
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LogrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *LogrusLogger) WithFields(fields map[string]interface{}) logger.Logger {
	return &LogrusLogger{logger: l.logger.WithFields(fields).Logger}
}
