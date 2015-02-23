package logger

import (
	"fmt"
	"github.com/t-k/fluent-logger-golang/fluent"
	"gopkg.in/yaml.v2"
	"time"
)

type MyLogger struct {
	fluent      *fluent.Fluent
	tagBasename string
}

func (logger *MyLogger) Print(tag string, key string, message string) {
	tagname := fmt.Sprintf("%s.%s", logger.tagBasename, tag)
	data := map[string]string{key: message}

	fmt.Printf("[%s][%s] %s : %s\n", time.Now(), tagname, key, message)

	if logger.fluent != nil {
		logger.fluent.PostWithTime(tagname, time.Now(), data)
	}
}

func (logger *MyLogger) ErrorPrint(tag string, message string) {
	logger.Print(tag, "error", message)
}

func (logger *MyLogger) LogPrint(tag string, message string) {
	logger.Print(tag, "log", message)
}

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

func (logger *MyLogger) Close() {
	if logger.fluent != nil {
		logger.fluent.Close()
		logger.fluent = nil
	}
}
