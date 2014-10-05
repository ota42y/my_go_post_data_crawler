package main

import (
  "fmt"
  "time"
  "os"
  "github.com/robfig/cron"
  "./work/twitter"
  "io/ioutil"
  "gopkg.in/yaml.v2"
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
	sendData := NewSendDataFromMap(configData)

	twitterWorker := twitter.NewWorkerFromMap(configData, loadYaml(setting_home + "/twitter.yml"), sendData.Database)

	c := cron.New()
	c.AddFunc("0 */10 * * * *", func() { twitterWorker.Work() })
	c.AddFunc("0 */1 * * * *", func() { sendData.SendData(100) })
	c.Start()

	for {
		_, err := os.Stat(setting_home + "/go_crawler_setting.yml")
		if err != nil{
			return
		}
		time.Sleep(1 * time.Minute)
		fmt.Println("sleep")
	}
}
