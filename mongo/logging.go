package mongo

import (
	// External Imports
	"github.com/sirupsen/logrus"
)

const (
	logError       = "datastore error"
	logConflict    = "resource conflict"
	logNotFound    = "resource not found"
	logNotHashable = "unable to hash secret"
)

// logger provides the package scoped logger implementation.
var logger storeLogger

// storeLogger provides a wrapper around the logrus logger in order to implement
// required database library logging interfaces.
type storeLogger struct {
	*logrus.Logger
}

// SetDebug turns on debug level logging, including debug at the driver level.
// If false, disables driver level logging and sets logging to info level.
func SetDebug(isDebug bool) {
	if isDebug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
}

// SetLogger enables binding in your own customised logrus logger.
func SetLogger(log *logrus.Logger) {
	logger = storeLogger{
		Logger: log,
	}
}
