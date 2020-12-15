package logger

import log "github.com/sirupsen/logrus"

func SetLogger() {
	log.SetFormatter(&log.JSONFormatter{})
}
