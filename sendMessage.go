package main

import (
	"./database"
	"./sendHubot"
	"fmt"
)

type SendData struct {
	Database *database.Database
	Server *sendHubot.Server
}

func NewSendDataFromMap(sendDataConfig map[interface{}]interface{}) (sendData *SendData) {
	dsn := sendDataConfig["dsn"].(string)
	postPath := sendDataConfig["postPath"].(string)
	defaultRoomName := sendDataConfig["defaultRoomName"].(string)
	return NewSendData(dsn, postPath, defaultRoomName)
}

func NewSendData(dsn string, postPath string, defaultRoomName string) (sendData *SendData) {
	return &SendData{
		Database: database.NewDatabase(dsn, defaultRoomName),
		Server: sendHubot.NewServer(postPath),
	}
}

func (sendData *SendData) SendData(limit int) {
	noSendPosts := sendData.Database.GetNoSendPosts(limit)
	if(len(noSendPosts) != 0){
		for _, post := range noSendPosts{
			if(sendData.Server.SendData(post.GetUrlValue())){
				sendData.Database.SendPost(post)
			}
		}
	}
}


func (sendData *SendData) TestData(){
	noSendPosts := sendData.Database.GetNoSendPosts(100)
	fmt.Println("get post ", len(noSendPosts))

	if(len(noSendPosts) != 0){
		fmt.Println(noSendPosts[0].MessageId)

		if(sendData.Server.SendData(noSendPosts[0].GetUrlValue())){
			sendData.Database.SendPost(noSendPosts[0])
		}
	}
}
