package logger

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
)

const (
	LogsDirectoryPath = "logs"
)

type Logger interface {
	// Debug(msg string, fields ...interface{})
	Debug(lf *LogFields)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type LogWrapper struct {
	logger *log.Logger
}

type LogFields struct {
	Message string
	Fields  interface{}
}

type HttpFields struct {
	Route  string
	Method string
	URL    *url.URL
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

	logger := log.New(filePath, "", log.LstdFlags)

	newLogger := &LogWrapper{
		logger: logger,
	}

	return newLogger
}

//	func (lw *LogWrapper) Debug(msg string, fields ...interface{}) {
//		fmt.Println("fields passed: ", fields)
//		lw.logger.Printf("level=DEBUG, message=%s, fields=%s", msg, fields)
//	}
func (lw *LogWrapper) Debug(lf *LogFields) {
	lw.logger.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("fields passed: ", lf)
	lw.logger.Printf("level=DEBUG, message=%s, fields=%+v", lf.Message, lf.Fields)
}

func (lw *LogWrapper) Info(msg string) {
	lw.logger.SetFlags(log.LstdFlags)
	lw.logger.Printf("level=INFO, message=%s", msg)
}

func (lw *LogWrapper) Warn(msg string) {
	lw.logger.SetFlags(log.LstdFlags)
	lw.logger.Printf("level=WARN, message=%s", msg)
}

func (lw *LogWrapper) Error(msg string) {
	lw.logger.SetFlags(log.LstdFlags)
	lw.logger.Printf("level=ERROR, message=%s", msg)
}
