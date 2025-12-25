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

	tests := []struct {
		Name        string
		LogName     string
		LogBody     string
		ExpectedLog string
	}{
		{
			"Valid info log",
			"test-info-log.log",
			"Test info log",
			"[INFO] Test info log\n",
		},
		{
			"Valid debug log",
			"test-debug-log.log",
			"Test debug log",
			"[INFO] Test debug log\n",
		},
		{
			"Valid error log",
			"test-error-log.log",
			"Test error log",
			"[ERROR] Test error log\n",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			loggerConfig := LoggerConfig{
				cwd,
				test.LogName,
			}

			logger, err := CreateLogger(loggerConfig)

			if err != nil {
				t.Fatalf("Error with setup, failed to create logger: %s\n", err)
			}

			if test.Name == "Valid info log" {
				logger.CreateLog(test.LogBody)
			} else if test.Name == "Valid debug log" {
				logger.CreateDebugLog(test.LogBody)
			} else if test.Name == "Valid error log" {
				logger.CreateErrorLog(test.LogBody)
			} else {
				CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
				t.Fatalf("Error with test, no test with name %s", test.Name)
			}

			logFile, err := os.ReadFile(
				path.Join(cwd, loggerConfig.LogName))

			if err != nil {
				CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
				t.Fatalf("Error with test, failed to open log file: %s", err)
			}

			logContent := string(logFile)

			regex := `^\[(.*?)\] \[(.*?)\] *(.*)\n$`
			re := regexp.MustCompile(regex)

			matches := re.FindStringSubmatch(logContent)
			if len(matches) != 4 {
				t.Errorf("Expected 3 matches but found %d\nMatches: %v\nLog:\n%s\n", len(matches), matches, logContent)
				CleanUpTest(t, path.Join(cwd, loggerConfig.LogDirectoryPath))
				t.FailNow()
			}

			logTime, err := time.Parse(time.RFC3339, matches[1])
			if err != nil {
				t.Errorf("Expected timestamp to be a valid UTC\nReceived: %s\n", logTime)
				CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
				t.FailNow()
			}

			if logTime.Location() != time.UTC {
				t.Errorf("Expected timestamp to be a valid UTC\nReceived: %s\n", logTime)
				CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
				t.FailNow()
			}

			var expectedLog string
			var expectedLogType string

			if test.Name == "Valid info log" {
				expectedLog = "Test info log"
				expectedLogType = "INFO"
			} else if test.Name == "Valid debug log" {
				expectedLog = "Test debug log"
				expectedLogType = "DEBUG"
			} else if test.Name == "Valid error log" {
				expectedLog = "Test error log"
				expectedLogType = "ERROR"
			} else {
				CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
				t.Fatalf("Error with test, no test with name %s", test.Name)
			}

			if strings.Trim(matches[2], "\t\r\n") != expectedLogType {
				t.Errorf("Expected log type \"%s\"\nReceived: \"%s\"\n", expectedLog, matches[2])
				CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
				t.FailNow()
			}

			if strings.Trim(matches[3], "\t\r\n") != expectedLog {
				t.Errorf("Expected log to \"%s\"\nReceived: \"%s\"\n", expectedLog, matches[3])
				CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
				t.FailNow()
			}

			CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
		})
	}
}

func TestWritingDebugLogs(t *testing.T) {
	cwd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Error with setup, failed to get cwd: %s\n", err)
	}
	
	loggerConfig := LoggerConfig{
		cwd,
		"debug-mode-test.log",
	}

	if err := os.Setenv("MODE", "DEVELOPMENT"); err != nil {
		t.Fatalf("Error with setup, failed to set env: %s\n", err)
	}

	logger, err := CreateLogger(loggerConfig)
	if err != nil {
		t.Fatalf("Error with test, failed to create logger: %s\n", err)
	}

	defer func() {
		os.Unsetenv("MODE")
	}()

	logger.CreateDebugLog("test debug mode log")

	logFile, err := os.ReadFile(path.Join(cwd, loggerConfig.LogName))
	if err != nil {
		t.Fatalf("Error with test, failed to read log: %s\n", err)
	}

	if string(logFile) == "" {
		CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
		t.Fatalf("Expected log file to be empty be received %s", string(logFile))
	}

	CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))

	if err := os.Setenv("MODE", "PRODUCTION"); err != nil {
		t.Fatalf("Error with test, failed to set env: %s\n", err)
	}

	logger, err = CreateLogger(loggerConfig)
	if err != nil {
		t.Fatalf("Error with test, failed to create logger: %s\n", err)
	}

	logger.CreateDebugLog("test debug mode log")

	logFile, err = os.ReadFile(path.Join(cwd, loggerConfig.LogName))
	if err != nil {
		t.Fatalf("Error with test, failed to read log: %s\n", err)
	}

	if string(logFile) != "" {
		CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
		t.Fatalf("Expected log file to be empty be received %s", string(logFile))
	}

	CleanUpTest(t, path.Join(cwd, loggerConfig.LogName))
}
