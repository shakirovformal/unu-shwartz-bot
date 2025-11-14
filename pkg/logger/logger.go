package logger

import (
	"log"
	"os"
)

type logger struct {
	*log.Logger
}

func NewLogger() *logger {
	return &logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}
