package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	LogsDirectoryPath = "logs"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
}

type LogWrapper struct {
	logger *log.Logger
}

func New() *LogWrapper {
	if os.Getenv("APP_ENV") == "" {
		return &LogWrapper{
			logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		}
	}

	os.Mkdir(LogsDirectoryPath, 0777)

	year, month, day := time.Now().Date()
	fileName := fmt.Sprintf("%v-%v-%v.log", year, month, day)
	filePath, err := os.OpenFile(LogsDirectoryPath+"/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error creating logger file: ", err)
		return nil
	}

	logger := log.New(filePath, "", log.LstdFlags|log.Lshortfile)

	newLogger := &LogWrapper{
		logger: logger,
	}

	return newLogger
}

func (lw *LogWrapper) Debug(msg string) {
	lw.logger.Printf("level=DEBUG, message=%s", msg)
}

func (lw *LogWrapper) Info(msg string) {
	lw.logger.Printf("level=INFO, message=%s", msg)
}

func (lw *LogWrapper) Warn(msg string) {
	lw.logger.Printf("level=WARN, message=%s", msg)
}

func (lw *LogWrapper) Error(msg string) {
	lw.logger.Printf("level=ERROR, message=%s", msg)
}
