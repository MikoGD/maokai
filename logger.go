package maokai

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

type Logger interface {
	CreateLog(body string) error
}

type LoggerConfig struct {
	logDirectoryPath string
	logName          string
}

type FileLogger struct {
	file   *os.File
	writer *bufio.Writer
}

func (fl *FileLogger) CreateLog(body string) error {
	currDateTime := time.Now().UTC()

	formattedLog := fmt.Sprintf("[%s]\n%s\n", currDateTime.Format(time.RFC3339), body)

	if _, err := fl.writer.WriteString(formattedLog); err != nil {
		return err
	}

	if err := fl.writer.Flush(); err != nil {
		return err
	}

	return nil
}

func CreateLogger(config LoggerConfig) (*FileLogger, error) {
	if config.logDirectoryPath == "" {
		return nil, &MissingLogDirectoryPathError{}
	}

	_, err := os.Stat(config.logDirectoryPath)

	if errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(config.logDirectoryPath, 0770)
		log.Printf("Path %s does not existing, creating directory\n", config.logDirectoryPath)
		if err != nil {
			return nil, err
		}
	}

	if config.logName == "" {
		return nil, &MissingLogNameError{}
	}

	logPath := path.Join(config.logDirectoryPath, config.logName)

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(file)

	logger := &FileLogger{file: file, writer: writer}

	return logger, nil
}
