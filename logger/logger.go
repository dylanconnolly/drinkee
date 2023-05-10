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
	New() *log.Logger
}

func New() *log.Logger {
	os.Mkdir(LogsDirectoryPath, 0777)

	year, month, day := time.Now().Date()
	fileName := fmt.Sprintf("%v-%v-%v.log", year, month, day)
	filePath, err := os.OpenFile(LogsDirectoryPath+"/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0750)
	if err != nil {
		fmt.Println("error creating logger file: ", err)
		return nil
	}

	logger := log.New(filePath, "", log.LstdFlags)
	return logger
}
