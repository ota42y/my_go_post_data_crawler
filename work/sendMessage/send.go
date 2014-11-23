package sendMessage

import (
	"../../lib/database"
	"../../lib/sendHubot"
)

type SendData struct {
	Database *database.Database
	Server   *sendHubot.Server
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
		Server:   sendHubot.NewServer(postPath),
	}
}

func (sendData *SendData) Work() {
	noSendPosts := sendData.Database.GetNoSendPosts(100)
	if len(noSendPosts) != 0 {
		for _, post := range noSendPosts {
			if sendData.Server.SendData(post.GetUrlValue()) {
				sendData.Database.SendPost(post)
			}
		}
	}
}
