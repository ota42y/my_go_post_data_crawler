package main

import (
	"./lib/evernote"
	"./lib/logger"
	"./work/chatLog"
	"./work/sendMessage"
	"./work/twitter"
	"./work/serverWorker"
	"fmt"
	"github.com/robfig/cron"
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

func main() {
	setting_home := os.Args[1]
	fmt.Println(setting_home)

	// evernote送信用
	evernote := evernote.NewSenderFromData(loadFile(setting_home + "/evernote.yml"))

	logger, err := logger.NewFromData("go_cron", loadFile(setting_home+"/fluent.yml"))
	if err != nil{
		logger.LogPrint("main", "----------")
		logger.LogPrint("main", "no fluentd")
		logger.LogPrint("main", "----------")
	}
	defer logger.Close()
	logger.LogPrint("main", "start")

	// hubotへのポスト用
	sendData := sendMessage.New(loadFile(setting_home + "/send_message.yml"))

	// twitter収集用
	twitterWorker := twitter.New(loadFile(setting_home+"/crawler.yml"),
		loadFile(setting_home+"/twitter.yml"), sendData.Database, logger)

	// チャットログ収集用
	chatLogWorker := chatLog.New(loadFile(setting_home+"/chatlog.yml"), logger, evernote)

	c := cron.New()
	c.AddFunc("0 */10 * * * *", func() { twitterWorker.Work() })
	c.AddFunc("0 */1 * * * *", func() { sendData.Work() })
	c.AddFunc("0 */10 * * * *", func() { chatLogWorker.Work() })
	c.Start()

	w := serverWorker.New(logger, sendData.Database)
	go w.Work()

	for {
		_, err := os.Stat(setting_home + "/go_crawler_setting.yml")
		if err != nil {
			return
		}
		logger.LogPrint("main", "sleep")
		time.Sleep(1 * time.Minute)
	}
}
