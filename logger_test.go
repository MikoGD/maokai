package maokai

import (
	"errors"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"
	"time"
)

func CleanUpTest(t *testing.T, logFilePath string) {
	err := os.Remove(logFilePath)
	if err != nil {
		t.Fatalf("Test clean up failed: %s\n", err)
	}
}

func TestCreatingInvalidLoggerConfig(t *testing.T) {
	cwd, err := os.Getwd()

	if err != nil {
		log.Fatalln(err)
	}

	tests := []struct {
		Name          string
		Config        LoggerConfig
		ExpectedError error
	}{
		{
			"Invalid config - missing path to log directory",
			LoggerConfig{
				"",
				"invalid-path-to-log-test-logs.log",
			},
			&MissingLogDirectoryPathError{},
		},
		{
			"Invalid config - missing log file name",
			LoggerConfig{
				cwd,
				"",
			},
			&MissingLogNameError{},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := CreateLogger(test.Config)
			if !errors.Is(err, test.ExpectedError) {
				if test.Config.LogDirectoryPath != "" {
					CleanUpTest(t, test.Config.LogDirectoryPath)
				}
				t.Errorf("\nExpected err: \"%s\"Received err: \"%s\"\n", test.ExpectedError, err)
			}
		})
	}
}

func TestWritingLogs(t *testing.T) {
	cwd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Error with setup, failed to get cwd: %s\n", err)
	}

	loggerConfig := LoggerConfig{
		cwd,
		"writing-test-log.log",
	}

	logger, err := CreateLogger(loggerConfig)

	if err != nil {
		t.Fatalf("Error with setup, failed to create logger: %s\n", err)
	}

	logger.CreateLog("test log")

	logFile, err := os.ReadFile(
		path.Join(cwd, loggerConfig.LogName))

	if err != nil {
		CleanUpTest(t, path.Join(cwd, loggerConfig.LogDirectoryPath))
		t.Fatalf("Error with test, failed to open log file: %s", err)
	}

	logContent := string(logFile)

	regex := `^\[(.*?)\]\s*(.*)\n$`
	re := regexp.MustCompile(regex)

	matches := re.FindStringSubmatch(logContent)
	if len(matches) != 3 {
		t.Errorf("Expected 3 matches but found %d\nMatches: %v\nLog:\n%s\n", len(matches), matches, logContent)
		CleanUpTest(t, path.Join(cwd, loggerConfig.LogDirectoryPath))
		t.FailNow()
	}

	logTime, err := time.Parse(time.RFC3339, matches[1])
	if err != nil {
		t.Errorf("Expected timestamp to be a valid UTC\nReceived: %s\n", logTime)
		CleanUpTest(t, path.Join(cwd, loggerConfig.LogDirectoryPath))
		t.FailNow()
	}

	if logTime.Location() != time.UTC {
		t.Errorf("Expected timestamp to be a valid UTC\nReceived: %s\n", logTime)
		CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
		t.FailNow()
	}

	if strings.Trim(matches[2], "\t\r\n") != "test log" {
		t.Errorf("Expected log to \"test log\"\nReceived: \"%s\"\n", matches[2])
		CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
		t.FailNow()
	}

	CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
}
