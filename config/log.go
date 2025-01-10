package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func InitLogging() {
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true, // Add colors to logs
	})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel) // Enable debug-level logging
}
