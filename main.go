package main

import (
	"./lib/logger"
	"./lib/evernote"
	"./work/sendMessage"
	"./work/chatLog"
	"./work/twitter"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

func loadFile(path string) (buf []byte) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return
}

func loadYaml(path string) map[interface{}]interface{} {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	config := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func main() {
	setting_home := os.Args[1]
	configData := loadYaml(setting_home + "/go_crawler_setting.yml")

	// evernote送信用
	evernote := evernote.NewSenderFromData(loadFile(setting_home+"/evernote.yml"))

	logger := logger.NewFromData("go_cron", loadFile(setting_home + "/fluent.yml"))
	defer logger.Close()
	logger.LogPrint("main", "start")

	// hubotへのポスト用
	sendData := sendMessage.New(loadFile(setting_home + "/send_message.yml"))

	// twitter収集用
	twitterWorker := twitter.New(loadFile(setting_home + "/crawler.yml"),
		loadFile(setting_home + "/twitter.yml"), sendData.Database, logger)

	// チャットログ収集用
	chatLogWorker := chatLog.NewWorkerFromMap(configData, logger, evernote)

	c := cron.New()
	c.AddFunc("0 */10 * * * *", func() { twitterWorker.Work() })
	c.AddFunc("0 */1 * * * *", func() { sendData.Work() })
	c.AddFunc("0 */10 * * * *", func() { chatLogWorker.Work() })
	c.Start()

	for {
		_, err := os.Stat(setting_home + "/go_crawler_setting.yml")
		if err != nil {
			return
		}
		logger.LogPrint("main", "sleep")
		time.Sleep(1 * time.Minute)
	}
}
