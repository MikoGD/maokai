package main

import (
	"bufio"
	"log"
	"os"
	"time"
	"fmt"
)

type Logger interface {
	CreateLog(body string) err
}

type LoggerContext struct {
	logFilePath string

}

func (l *Logger) CreateLog(body string) err {
	currDateTime := time.Now().UTC()

	return fmt.Sprintf("[%s]\n%s\n", currDateTime.Format(time.RFC3339), body)
}

func CreateLogger() Logger {
	defaultLogFilePath := "/var/logs/maokai-logs"
	defaultLogFileName := "maokai-log.log"

	logger := Logger{}

	return logger
}

func main() {
	file, err := os.OpenFile("example.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	var bytesWritten int

	if bytesWritten, err = writer.WriteString(createLog("Log example")); err != nil {
		panic(err)
	}

	if err := writer.Flush(); err != nil {
		panic(err)
	}

	log.Printf("Bytes written %d\n", bytesWritten)
}
