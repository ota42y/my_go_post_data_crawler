package logger

import (
	"github.com/t-k/fluent-logger-golang/fluent"
	"time"
  "fmt"
)

type MyLogger struct {
	fluent *fluent.Fluent
  tagBasename string
}

func (logger *MyLogger) Print(tag string, key string, message string) {
  tagname := fmt.Sprintf("%s.%s", logger.tagBasename, tag)
	data := map[string]string{key: message}
	logger.fluent.PostWithTime(tagname, time.Now(), data)
	fmt.Printf("[%s][%s] %s : %s\n", time.Now(), tagname, key, message)
}

func (logger *MyLogger) LogPrint(tag string, message string) {
	logger.Print(tag, "log", message)
}

func NewFromMap(tagBasename string, config map[interface{}]interface{}) (logger *MyLogger) {
	fluentData := config["fluentd"].(map[interface{}]interface{})

	host := fluentData["host"].(string)
	port := fluentData["port"].(int)
	logger, err := New(tagBasename, fluent.Config{FluentPort: port, FluentHost: host})
	if err != nil {
		panic(err)
	}

	return
}


func New(tagBasename string, config fluent.Config) (logger *MyLogger, err error) {
	flu, err := fluent.New(config)
	if err != nil {
		return nil, err
	}
	logger = &MyLogger{
		fluent: flu,
    tagBasename: tagBasename,
	}
	return logger, nil
}

func (logger *MyLogger) Close() {
	logger.fluent.Close()
	logger.fluent = nil
}
