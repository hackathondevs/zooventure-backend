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
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	if os.Getenv("ENVIRONMENT") == "production" {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	return log
}
