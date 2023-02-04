package logger

type emptyLogger struct {
}

// NewEmptyLogger :
func NewEmptyLogger() ILogger {
	return &emptyLogger{}
}

// Debug :
func (l *emptyLogger) Debug(a ...interface{}) {

}

// Debugf :
func (l *emptyLogger) Debugf(format string, prm ...interface{}) {

}

// Print :
func (l *emptyLogger) Print(a ...interface{}) {

}

// Print :
func (l *emptyLogger) Println(a ...interface{}) {

}

// Printf :
func (l *emptyLogger) Printf(format string, prm ...interface{}) {

}

// Info :
func (l *emptyLogger) Info(a ...interface{}) {

}

// Infof :
func (l *emptyLogger) Infof(format string, prm ...interface{}) {

}

// Warn :
func (l *emptyLogger) Warn(a ...interface{}) {

}

// Warnf :
func (l *emptyLogger) Warnf(format string, prm ...interface{}) {

}

// Error :
func (l *emptyLogger) Error(a ...interface{}) {

}

// Errorf :
func (l *emptyLogger) Errorf(format string, prm ...interface{}) {

}

// Fatal :
func (l *emptyLogger) Fatal(a ...interface{}) {

}

// Fatalf :
func (l *emptyLogger) Fatalf(format string, prm ...interface{}) {

}

// Sync :
func (l *emptyLogger) Close() {
}
