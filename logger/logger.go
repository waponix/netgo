package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/waponix/netgo/utils/sliceUtil"
)

// log type constants
const (
	INFO   = "INFO"
	DEBUG  = "DEBUG"
	NOTICE = "NOTICE"
	ERROR  = "ERROR"
	FATAL  = "FATAL"
)

// format constants
const (
	FORMAT          = "[%s] (%d) %s.%s: %s"
	DATETIME_FORMAT = "2006-01-02 15:04:05"
)

type LogInterface interface {
	Info(string) error
	Debug(string) error
	Notice(string) error
	Error(string) error
	Fatal(string) error
}

type Log struct {
	Filename  string
	LogLevels []string
}

// public getter for the logger struct
func New() *Log {
	// set the default values for the filename
	return &Log{
		Filename:  "",
		LogLevels: []string{INFO, DEBUG, NOTICE, ERROR, FATAL},
	}
}

// public: write info log
func (l *Log) Info(message string) error {
	var err error = nil
	// only log when log level is present in the LogLevels
	if sliceUtil.Use(l.LogLevels).InItems(INFO) {
		err = l.writeLog(composeLogMessage(message, INFO))
	}
	return err
}

// public: write debug log
func (l *Log) Debug(message string) error {
	var err error = nil
	// only log when log level is present in the LogLevels
	if sliceUtil.Use(l.LogLevels).InItems(DEBUG) {
		err = l.writeLog(composeLogMessage(message, DEBUG))
	}
	return err
}

// public: write notice log
func (l *Log) Notice(message string) error {
	var err error = nil
	// only log when log level is present in the LogLevels
	if sliceUtil.Use(l.LogLevels).InItems(NOTICE) {
		err = l.writeLog(composeLogMessage(message, NOTICE))
	}
	return err
}

// public: write error log
func (l *Log) Error(message string) error {
	var err error = nil
	// only log when log level is present in the LogLevels
	if sliceUtil.Use(l.LogLevels).InItems(ERROR) {
		err = l.writeLog(composeLogMessage(message, ERROR))
	}
	return err
}

// public: write fatal log
func (l *Log) Fatal(message string) error {
	var err error = nil
	// only log when log level is present in the LogLevels
	if sliceUtil.Use(l.LogLevels).InItems(FATAL) {
		err = l.writeLog(composeLogMessage(message, FATAL))
	}
	return err
}

// does the actual writing to the log
func (l *Log) writeLog(line string) error {
	// Open the file for appending (or create it if it doesn't exist)
	file, err := os.OpenFile(l.Filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Line to append to the file
	line += "\n"

	// Write the line to the file
	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	return nil
}

func composeLogMessage(message string, messageType string) string {
	file, line := getCallerInfo(3)
	currentTime := time.Now()
	return fmt.Sprintf(FORMAT, currentTime.Format(DATETIME_FORMAT), line, file, messageType, message)
}

func getCallerInfo(level int) (string, int) {
	// Get information about the caller at depth 1 (the immediate caller)
	_, file, line, ok := runtime.Caller(level)
	if ok {
		return file, line
	} else {
		// this will most likely never be reached
		return "Unknown", 0
	}
}
