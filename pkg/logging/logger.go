package logging

import (
	log "github.com/sirupsen/logrus"
)

func SetupLogger(contextualFields map[string]interface{}) {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields(contextualFields))
}
