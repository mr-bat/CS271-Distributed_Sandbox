package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = logrus.New()

func initLogger()  {
	Logger.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	Logger.Level = logrus.InfoLevel
	Logger.Out = os.Stdout
}
