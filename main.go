package main

import (
  "time"
  "os"
  "github.com/robfig/cron"
  "./work/twitter"
  "./work/chatLog"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "./logger"
)

func loadYaml(path string) (map[interface{}]interface{}){
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

func main(){
	setting_home := os.Args[1]
	configData := loadYaml(setting_home + "/go_crawler_setting.yml")

  logger := logger.NewFromMap("go_cron", configData)
  defer logger.Close()
  logger.LogPrint("main", "start")

	sendData := NewSendDataFromMap(configData)

	twitterWorker := twitter.NewWorkerFromMap(configData, loadYaml(setting_home + "/twitter.yml"), sendData.Database)
  chatLogWorker := chatLog.NewWorkerFromMap(configData)

	c := cron.New()
	c.AddFunc("0 */10 * * * *", func() { twitterWorker.Work() })
	c.AddFunc("0 */1 * * * *", func() { sendData.SendData(100) })
  c.AddFunc("0 */10 * * * *", func() { chatLogWorker.Work() })
	c.Start()

	for {
		_, err := os.Stat(setting_home + "/go_crawler_setting.yml")
		if err != nil{
			return
		}
		time.Sleep(1 * time.Minute)
		logger.LogPrint("main", "sleep")
	}
}
