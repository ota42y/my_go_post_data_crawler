package pomodoro

import (
	"./../../lib/database"
	"./../../lib/server"
	"github.com/robfig/cron"
	"time"
	"fmt"
	"github.com/ota42y/go-tumblr/tumblr"
	"gopkg.in/yaml.v2"
)

type Setting struct {
	ConsumerKey string
	ConsumerSecret string
	AccessToken string
	AccessTokenSecret string
}

type Command struct{
	server *server.Server
	sendRoomName string
	cron *cron.Cron
	isStart bool
	blog *tumblr.BlogApi
}

func New(server *server.Server, sendRoomName string, setting []byte) (c *Command){
	s := Setting{}
	err := yaml.Unmarshal(setting, &s)
	if err != nil {
		return nil
	}

	c = &Command{
		server: server,
		sendRoomName: sendRoomName,
		cron: cron.New(),
		isStart: false,
		blog: nil,
	}

	c.cron.AddFunc("*/30 * * * * *", func() { c.sendMessage() })

	return c
}

func (c *Command) IsExecute(command string) bool{
	return command == "pomodoro"
}

func (c *Command) Execute(data string) string{
	if c.isStart {
		c.cron.Stop()
		c.isStart = false
		return "pomodoro: stop"
	}else{
		c.cron.Start()
		c.isStart = true
		return "pomodoro: start"
	}
}

func (c *Command) sendMessage() {
	now := fmt.Sprintf("%d", time.Now().Unix())
	s := c.server
	s.SendPost(database.NewPost(c.sendRoomName, "pomodoro: 進捗どうですか？", "pomodorocommand:" + now))
}
