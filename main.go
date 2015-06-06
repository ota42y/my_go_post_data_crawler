package main

import (
	"./config"
	"./lib/evernote"
	"./lib/logger"
	"./lib/post"
	"./util"
	"./work/backup/mongodb"
	"./work/chatLog"
	"./work/crawler/twitter"
	"./work/sendMessage"
	"./work/serverWorker"
	"./worker"
	"fmt"
	"github.com/robfig/cron"
	"math/rand"
	"os"
	"time"

	"./work/checker/log/error"
	"./work/checker/log/info"
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
	sender := post.NewSender(sendData.Database)

	// チャットログ収集用
	chatLogWorker := chatLog.New(util.LoadFile(setting_home+"/chatlog.yml"), logger, evernote)

	// mongodbバックアップ用
	hourlyWorker := worker.NewWorker()
	hourlyWorker.AddWork(mongodb.NewMongodb(
		config.NewMongodbBackupFromData(util.LoadFile(setting_home+"/mongodb_backup.yml")),
		config.NewMongodbDatabaseFromData(util.LoadFile(setting_home+"/mongodb_logserver.yml")),
		logger))

	// check error log
	hourlyWorker.AddWork(error.NewLogCollector(
		config.NewMongodbDatabaseFromData(util.LoadFile(setting_home+"/mongodb_logserver.yml")),
		sender,
		logger))

	// check info log
	hourlyWorker.AddWork(info.NewLogCollector(
		config.NewMongodbDatabaseFromData(util.LoadFile(setting_home+"/mongodb_logserver.yml")),
		sender,
		logger))

	// 10 minutes worker
	tenMinutesWorker := worker.NewWorker()
	tenMinutesWorker.AddWork(twitter.NewTwitter(
		util.LoadFile(setting_home+"/crawler.yml"),
		util.LoadFile(setting_home+"/twitter.yml"),
		sender, logger))

	c := cron.New()
	c.AddFunc("0 */1 * * * *", func() { sendData.Work() })
	c.AddFunc("0 */10 * * * *", func() { chatLogWorker.Work() })

	c.AddFunc("0 2 * * * *", func() { hourlyWorker.Work() })
	c.AddFunc("0 */10 * * * *", func() { tenMinutesWorker.Work() })
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
