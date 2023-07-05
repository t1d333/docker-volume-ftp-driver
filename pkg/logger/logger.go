package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitializeNewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = new(logrus.TextFormatter)
	logger.Formatter.(*logrus.TextFormatter).DisableColors = false
	logger.Level = logrus.TraceLevel
	logger.Out = os.Stdout
	return logger
}
