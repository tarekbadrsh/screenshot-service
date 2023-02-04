package logger

import (
	"bytes"
	"fmt"
	"sync"
)

var (
	internalLogger *ILogger
	initLoggerOnce sync.Once
)

// ILogger :
type ILogger interface {
	Debug(a ...interface{})
	Debugf(format string, prm ...interface{})
	Print(a ...interface{})
	Printf(format string, prm ...interface{})
	Println(a ...interface{})
	Info(a ...interface{})
	Infof(format string, prm ...interface{})
	Warn(a ...interface{})
	Warnf(format string, prm ...interface{})
	Error(a ...interface{})
	Errorf(format string, prm ...interface{})
	Fatal(a ...interface{})
	Fatalf(format string, prm ...interface{})
	Close()
}

// WithFields creates an string from the string map and adds multiple to return string.
func WithFields(f map[string]interface{}) string {
	b := new(bytes.Buffer)
	for key, value := range f {
		fmt.Fprintf(b, "%s=\"%s\" ", key, value)
	}
	return b.String()
}

// InitializeLogger : initialize the main logger
func InitializeLogger(log *ILogger) {
	initLoggerOnce.Do(func() {
		internalLogger = log
		(*internalLogger).Info("logger has been initialized")
	})
}

// GetLogger :
func GetLogger() *ILogger {
	return internalLogger
}

// Debug :
func Debug(a ...interface{}) {
	(*internalLogger).Debug(a...)
}

// Debugf :
func Debugf(format string, prm ...interface{}) {
	(*internalLogger).Debugf(format, prm...)
}

// Print :
func Print(a ...interface{}) {
	(*internalLogger).Print(a...)
}

// Println :
func Println(a ...interface{}) {
	(*internalLogger).Println(a...)
}

// Printf :
func Printf(format string, prm ...interface{}) {
	(*internalLogger).Printf(format, prm...)
}

// Info :
func Info(a ...interface{}) {
	(*internalLogger).Info(a...)
}

// Infof :
func Infof(format string, prm ...interface{}) {
	(*internalLogger).Infof(format, prm...)
}

// Warn :
func Warn(a ...interface{}) {
	(*internalLogger).Warn(a...)
}

// Warnf :
func Warnf(format string, prm ...interface{}) {
	(*internalLogger).Warnf(format, prm...)
}

// Error :
func Error(a ...interface{}) {
	(*internalLogger).Error(a...)
}

// Errorf :
func Errorf(format string, prm ...interface{}) {
	(*internalLogger).Errorf(format, prm...)
}

// Fatal :
func Fatal(a ...interface{}) {
	(*internalLogger).Fatal(a...)
}

// Fatalf :
func Fatalf(format string, prm ...interface{}) {
	(*internalLogger).Fatalf(format, prm...)
}

// Close :
func Close() {
	(*internalLogger).Close()
}
