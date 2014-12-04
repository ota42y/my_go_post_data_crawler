package logger

import (
	"github.com/t-k/fluent-logger-golang/fluent"
	"gopkg.in/yaml.v2"
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

func (logger *MyLogger) ErrorPrint(tag string, message string) {
	logger.Print(tag, "error", message)
}

func (logger *MyLogger) LogPrint(tag string, message string) {
	logger.Print(tag, "log", message)
}

func NewFromData(tagBasename string, buf []byte) (logger *MyLogger) {
	c := fluent.Config{}
	err := yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil
	}

	flu, err := fluent.New(c)
	if err != nil {
		return nil
	}

	logger = &MyLogger{
		fluent: flu,
		tagBasename: tagBasename,
	}

	return
}

func (logger *MyLogger) Close() {
	logger.fluent.Close()
	logger.fluent = nil
}
