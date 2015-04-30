package main

import (
	"./config"
	"./lib/evernote"
	"./lib/logger"
	"./util"
	"./work/backup/mongodb"
	"./work/chatLog"
	"./work/sendMessage"
	"./work/serverWorker"
	"./work/twitter"
	"./worker"
	"fmt"
	"github.com/robfig/cron"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	setting_home := os.Args[1]
	fmt.Println(setting_home)

	// evernote送信用
	evernote := evernote.NewSenderFromData(util.LoadFile(setting_home + "/evernote.yml"))

	logger, err := logger.NewFromData("go_cron", util.LoadFile(setting_home+"/fluent.yml"))
	if err != nil {
		logger.LogPrint("main", "----------")
		logger.LogPrint("main", "no fluentd")
		logger.LogPrint("main", "----------")
	}
	defer logger.Close()
	logger.LogPrint("main", "start")

	// hubotへのポスト用
	sendData := sendMessage.New(util.LoadFile(setting_home + "/send_message.yml"))

	// twitter収集用
	twitterWorker := twitter.New(util.LoadFile(setting_home+"/crawler.yml"),
		util.LoadFile(setting_home+"/twitter.yml"), sendData.Database, logger)

	// チャットログ収集用
	chatLogWorker := chatLog.New(util.LoadFile(setting_home+"/chatlog.yml"), logger, evernote)

	// mongodbバックアップ用
	dailyWorker := worker.NewWorker()
	dailyWorker.AddWork(mongodb.NewMongodb(
		config.NewMongodbBackupFromData(util.LoadFile(setting_home+"/mongodb_backup.yml")),
		config.NewMongodbDatabaseFromData(util.LoadFile(setting_home+"/mongodb_logserver.yml")),
		logger))

	c := cron.New()
	c.AddFunc("0 */10 * * * *", func() { twitterWorker.Work() })
	c.AddFunc("0 */1 * * * *", func() { sendData.Work() })
	c.AddFunc("0 */10 * * * *", func() { chatLogWorker.Work() })
	c.AddFunc("0 2 * * * *", func() { dailyWorker.Work() })
	c.Start()

	w := serverWorker.New(logger, sendData.Database, setting_home)
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
