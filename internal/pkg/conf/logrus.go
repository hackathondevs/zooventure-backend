package conf

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	level, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	log := logrus.New()
	log.SetLevel(logrus.Level(level))
	return log
}
