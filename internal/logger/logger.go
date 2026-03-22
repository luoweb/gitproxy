package logger

import (
	"log"
	"strings"
	"sync"
)

type Level int

const (
	DEBUG Level = 0
	INFO  Level = 1
	WARN  Level = 2
	ERROR Level = 3
)

func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	IsDebugEnabled() bool
}

type SimpleLogger struct {
	level  Level
	format string
	mu     sync.Mutex
}

func New(level string, format string) *SimpleLogger {
	return &SimpleLogger{
		level:  ParseLevel(level),
		format: format,
	}
}

func (l *SimpleLogger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.log("DEBUG", format, v...)
	}
}

func (l *SimpleLogger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.log("INFO", format, v...)
	}
}

func (l *SimpleLogger) Warn(format string, v ...interface{}) {
	if l.level <= WARN {
		l.log("WARN", format, v...)
	}
}

func (l *SimpleLogger) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.log("ERROR", format, v...)
	}
}

func (l *SimpleLogger) IsDebugEnabled() bool {
	return l.level <= DEBUG
}

func (l *SimpleLogger) log(level, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log.Printf("[%s] "+format, append([]interface{}{level}, v...)...)
}
