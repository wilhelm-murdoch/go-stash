package tools

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetOutput(os.Stdout)
}

func Error(args ...any) {
	logger.Error(args...)
}

func Fatal(args ...any) {
	logger.Fatal(args...)
}

func Warning(args ...any) {
	logger.Warning(args...)
}

func Info(args ...any) {
	logger.Info(args...)
}

func Debug(args ...any) {
	logger.Debug(args...)
}
