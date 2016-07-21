package util

import (
	"io"
	"log"
	"os"
)

var logger *log.Logger

// NewLogger return a new instance of a logger
func NewLogger(name string) (*log.Logger, error) {

	logName := "./logs/" + name + ".log"

	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file: %v", err)
		return nil, err
	}

	multi := io.MultiWriter(file, os.Stdout)

	instance := log.New(multi,
		name+": ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return instance, nil
}

// Logger return main logger instance
func Logger() *log.Logger {

	if logger == nil {
		appLogger, err := NewLogger("app")
		if err != nil {
			log.Fatalf("Cannot load default logger\n %v", err)
			panic(err)
		}
		logger = appLogger
	}

	return logger
}
