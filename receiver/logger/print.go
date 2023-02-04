package logger

import "fmt"

type printLogger struct {
}

// NewPrintLogger :
func NewPrintLogger() ILogger {
	return &printLogger{}
}

// Debug :
func (l *printLogger) Debug(a ...interface{}) {
	fmt.Println(a...)
}

// Debugf :
func (l *printLogger) Debugf(format string, prm ...interface{}) {
	fmt.Printf(format, prm...)
}

// Print :
func (l *printLogger) Print(a ...interface{}) {
	fmt.Println(a...)
}

// Print :
func (l *printLogger) Println(a ...interface{}) {
	fmt.Println(a...)
}

// Printf :
func (l *printLogger) Printf(format string, prm ...interface{}) {
	fmt.Printf(format, prm...)
}

// Info :
func (l *printLogger) Info(a ...interface{}) {
	fmt.Println(a...)
}

// Infof :
func (l *printLogger) Infof(format string, prm ...interface{}) {
	fmt.Printf(format, prm...)
}

// Warn :
func (l *printLogger) Warn(a ...interface{}) {
	fmt.Println(a...)
}

// Warnf :
func (l *printLogger) Warnf(format string, prm ...interface{}) {
	fmt.Printf(format, prm...)
}

// Error :
func (l *printLogger) Error(a ...interface{}) {
	fmt.Println(a...)
}

// Errorf :
func (l *printLogger) Errorf(format string, prm ...interface{}) {
	fmt.Printf(format, prm...)
}

// Fatal :
func (l *printLogger) Fatal(a ...interface{}) {
	fmt.Println(a...)
}

// Fatalf :
func (l *printLogger) Fatalf(format string, prm ...interface{}) {
	fmt.Printf(format, prm...)
}

// Sync :
func (l *printLogger) Close() {
}
