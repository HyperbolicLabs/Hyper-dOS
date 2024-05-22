package helper

import (
	"github.com/sirupsen/logrus"
)

func SetLogLevel(level string) {
	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		// error
		logrus.Fatalf("invalid loglevel: %v", level)
	}

	logrus.Infof("loglevel is set to: %v", level)
}
