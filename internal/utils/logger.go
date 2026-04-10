package utils

import (
	"fmt"
	"log"
	"os"
)

// Logger is a simple logger
type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

// NewLogger creates a new logger
func NewLogger() *Logger {
	return &Logger{
		InfoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile),
		ErrorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile),
	}
}

// Info logs an info message
func (l *Logger) Info(message string, args ...interface{}) {
	l.InfoLogger.Println(fmt.Sprintf(message, args...))
}

// Error logs an error message
func (l *Logger) Error(message string, args ...interface{}) {
	l.ErrorLogger.Println(fmt.Sprintf(message, args...))
}

// Fatal logs a fatal error and exits
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.ErrorLogger.Fatalln(fmt.Sprintf(message, args...))
}
