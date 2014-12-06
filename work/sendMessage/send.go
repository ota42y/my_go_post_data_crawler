package sendMessage

import (
	"../../lib/database"
	"../../lib/sendHubot"
	"gopkg.in/yaml.v2"
)

type Setting struct {
	DataBaseName string
	HubotPostPath string
	DefaultChatRoomName string
}

type SendData struct {
	Database *database.Database
	Server   *sendHubot.Server
}

func New(buf []byte) (sendData *SendData) {
	s := &Setting{}
	err := yaml.Unmarshal(buf, s)
	if err != nil {
		return nil
	}

	return &SendData{
		Database: database.NewDatabase(s.DataBaseName, s.DefaultChatRoomName),
		Server:   sendHubot.NewServer(s.HubotPostPath),
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
