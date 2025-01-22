package logger

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

var (
	Logger *log.Logger
)

// InitLogger initializes the logger with the given log file path
func InitLogger(logFilePath string) error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		return err
	}

	// Open the log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// Create a new logger
	Logger = log.NewWithOptions(logFile, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Level:           log.DebugLevel,
		Prefix:          "nexus",
	})

	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *log.Logger {
	if Logger == nil {
		// If logger hasn't been initialized, create a default one that logs to a file in ~/.nexus/logs
		homeDir, err := os.UserHomeDir()
		if err != nil {
			Logger = log.NewWithOptions(os.Stderr, log.Options{
				ReportCaller:    true,
				ReportTimestamp: true,
				Level:           log.DebugLevel,
				Prefix:          "nexus",
			})
			return Logger
		}
		logDir := filepath.Join(homeDir, ".nexus", "logs")
		defaultLogPath := filepath.Join(logDir, "nexus.log")
		if err := InitLogger(defaultLogPath); err != nil {
			// If we can't create the file logger, fall back to stderr
			Logger = log.NewWithOptions(os.Stderr, log.Options{
				ReportCaller:    true,
				ReportTimestamp: true,
				Level:           log.DebugLevel,
				Prefix:          "nexus",
			})
		}
	}
	return Logger
}
