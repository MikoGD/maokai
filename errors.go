package maokai

import (
	"fmt"
)

type MissingLogDirectoryPathError struct{}
type MissingLogNameError struct{}
type DirectoryDoesNotExistError struct{ directoryPath string }

func (e *MissingLogDirectoryPathError) Error() string {
	return "Missing logDirecotoryPath in LoggerConfig"
}

func (e *MissingLogNameError) Error() string {
	return "Missing logName in LoggerConfig"
}

func (e *DirectoryDoesNotExistError) Error() string {
	return fmt.Sprintf("Directory %s does not exist", e.directoryPath)
}
