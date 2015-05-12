package logger

import (
	"fmt"
	"github.com/t-k/fluent-logger-golang/fluent"
	"gopkg.in/yaml.v2"
	"time"
)

// Logger is log interface
type Logger interface {
	Debug(tag string, format string, args ...interface{})
	Info(tag string, format string, args ...interface{})
	Error(tag string, format string, args ...interface{})
}

// MyLogger send log fluent and Stdout
type MyLogger struct {
	fluent      *fluent.Fluent
	tagBasename string
}

// Debug write error log
func (logger *MyLogger) Debug(tag string, format string, args ...interface{}) {
	logger.Print(tag, "debug", fmt.Sprintf(format, args...))
}

// Info write info log
func (logger *MyLogger) Info(tag string, format string, args ...interface{}) {
	logger.Print(tag, "info", fmt.Sprintf(format, args...))
}

// Error write error log
func (logger *MyLogger) Error(tag string, format string, args ...interface{}) {
	logger.ErrorPrint(tag, fmt.Sprintf(format, args...))
}

// Print write log to fluentd and Stdout
// Deprecate
func (logger *MyLogger) Print(tag string, key string, message string) {
	tagname := fmt.Sprintf("%s.%s", logger.tagBasename, tag)
	data := map[string]string{key: message}

	fmt.Printf("[%s][%s] %s : %s\n", time.Now(), tagname, key, message)

	if logger.fluent != nil {
		logger.fluent.PostWithTime(tagname, time.Now(), data)
	}
}

// ErrorPrint is deprecate, use Error
func (logger *MyLogger) ErrorPrint(tag string, message string) {
	logger.Print(tag, "error", message)
}

// LogPrint is deprecate, use Info
func (logger *MyLogger) LogPrint(tag string, message string) {
	logger.Print(tag, "log", message)
}

// NewFromData create logger from yaml data
func NewFromData(tagBasename string, buf []byte) (logger *MyLogger, err error) {
	c := fluent.Config{}
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, nil
	}

	flu, err := fluent.New(c)
	logger = &MyLogger{
		fluent:      flu,
		tagBasename: tagBasename,
	}

	return
}

// Close close fluent socket if it was open.
func (logger *MyLogger) Close() {
	if logger.fluent != nil {
		logger.fluent.Close()
		logger.fluent = nil
	}
}
