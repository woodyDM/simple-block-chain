package core

import "github.com/sirupsen/logrus"

var Log *logrus.Logger = logrus.New()

func init() {
	Log.SetLevel(logrus.DebugLevel)
}
