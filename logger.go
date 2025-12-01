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
	LogDirectoryPath string
	LogName          string
}

type FileLogger struct {
	File   *os.File
	Writer *bufio.Writer
}

func (fl *FileLogger) CreateLog(body string) error {
	currDateTime := time.Now().UTC()

	formattedLog := fmt.Sprintf("[%s]\n%s\n", currDateTime.Format(time.RFC3339), body)

	if _, err := fl.Writer.WriteString(formattedLog); err != nil {
		return err
	}

	if err := fl.Writer.Flush(); err != nil {
		return err
	}

	return nil
}

func CreateLogger(config LoggerConfig) (*FileLogger, error) {
	if config.LogDirectoryPath == "" {
		return nil, &MissingLogDirectoryPathError{}
	}

	_, err := os.Stat(config.LogDirectoryPath)

	if errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(config.LogDirectoryPath, 0770)
		log.Printf("Path %s does not existing, creating directory\n", config.LogDirectoryPath)
		if err != nil {
			return nil, err
		}
	}

	if config.LogName == "" {
		return nil, &MissingLogNameError{}
	}

	logPath := path.Join(config.LogDirectoryPath, config.LogName)

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(file)

	logger := &FileLogger{File: file, Writer: writer}

	return logger, nil
}
