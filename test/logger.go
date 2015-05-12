package test

import (
	"bytes"
	"fmt"
)

// Logger is stub logger class
type Logger struct {
	B bytes.Buffer
}

// Error write TestLogger.B to error log
func (l *Logger) Error(tag string, format string, args ...interface{}) {
	fmt.Fprintf(&l.B, format, args...)
}

// Info write TestLogger.B to info log
func (l *Logger) Info(tag string, format string, args ...interface{}) {
	fmt.Fprintf(&l.B, format, args...)
}

// Debug write TestLogger.B to info log
func (l *Logger) Debug(tag string, format string, args ...interface{}) {
	fmt.Fprintf(&l.B, format, args...)
}
