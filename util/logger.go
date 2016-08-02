package util

import (
	"io"
	"log"
	"os"
)

const defaultLoggerName = "app"

var loggers map[string]*log.Logger

// NewLogger return a new instance of a logger
func NewLogger(name string) (*log.Logger, error) {

	if loggers == nil {
		loggers = make(map[string]*log.Logger)
	}

	if loggers[name] != nil {
		return loggers[name], nil
	}

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

	loggers[name] = instance
	return instance, nil
}

// Logger return main logger instance
func Logger() *log.Logger {

	if loggers[defaultLoggerName] == nil {
		appLogger, err := NewLogger("app")
		if err != nil {
			log.Fatalf("Cannot load default logger\n %v", err)
			panic(err)
		}
		loggers[defaultLoggerName] = appLogger
	}

	return loggers[defaultLoggerName]
}
