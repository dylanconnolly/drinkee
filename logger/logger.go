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
	err := os.Mkdir(LogsDirectoryPath, 0777)
	if err != nil && err != os.ErrExist {
		fmt.Printf("error creating logger: %+v", err)
	}

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
