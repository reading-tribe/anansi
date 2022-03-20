package logging

import (
	log "github.com/sirupsen/logrus"
)

type AnansiLogger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type anansiLogger struct {
	Fields map[string]interface{}
}

func NewLogger(contextualFields map[string]interface{}) AnansiLogger {
	log.SetFormatter(&log.JSONFormatter{})

	return anansiLogger{
		Fields: contextualFields,
	}
}

func (a anansiLogger) Info(args ...interface{}) {
	log.WithFields(a.Fields).Info(args)
}
func (a anansiLogger) Debug(args ...interface{}) {
	log.WithFields(a.Fields).Debug(args)
}
func (a anansiLogger) Warn(args ...interface{}) {
	log.WithFields(a.Fields).Warn(args)
}
func (a anansiLogger) Error(args ...interface{}) {
	log.WithFields(a.Fields).Error(args)
}
