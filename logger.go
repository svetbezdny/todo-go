package main

import (
	"fmt"
	"time"
)

const (
	LevelInfo    = "INFO"
	LevelWarning = "WARNING"
	LevelDebug   = "DEBUG"
	LevelError   = "ERROR"
)

type Logger struct {
	Timestamp string
	Level     string
	Message   string
}

func Log(message string) Logger {
	return Logger{
		Timestamp: time.Now().Format("2006.01.02 15:04:05"),
		Level:     LevelInfo,
		Message:   message,
	}
}

func (l Logger) WithLevel(level string) Logger {
	l.Level = level
	return l
}

func (l Logger) String() string {
	return fmt.Sprintf("%s | [%s] | %s", l.Timestamp, l.Level, l.Message)
}
