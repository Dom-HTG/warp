package utils

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	logfile, err := os.OpenFile("warplogs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}

	//Initialize logger with appropriate settings.
	logrus.SetLevel(logrus.DebugLevel)
	log.SetOutput(logfile)
}
