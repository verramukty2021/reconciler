package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
)

// Initialize loggers
func init() {
	logFile := os.Stdout // You can set this to a file if needed

	infoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
	warnLogger = log.New(logFile, "WARN: ", log.Ldate|log.Ltime)
	errorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info logs an informational message with variadic arguments
func Info(msg string, args ...interface{}) {
	infoLogger.Println(formatMessage(msg, args...))
}

// Warn logs a warning message with variadic arguments
func Warn(msg string, args ...interface{}) {
	warnLogger.Println(formatMessage(msg, args...))
}

// Error logs an error message with variadic arguments
func Error(msg string, args ...interface{}) {
	errorLogger.Println(formatMessage(msg, args...))
}

// formatMessage formats the message with optional args
func formatMessage(msg string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}
