package periodic

import (
	"../command"
	"./../../lib/database"
	"./../../lib/server"
	"fmt"
	"github.com/robfig/cron"
	"time"
)

type Command struct {
	server       *server.Server
	sendRoomName string
	cron         *cron.Cron
}

func New(server *server.Server, sendRoomName string) (c *Command) {
	c = &Command{
		server:       server,
		sendRoomName: sendRoomName,
		cron:         cron.New(),
	}

	c.cron.AddFunc("0 0 7-23 * * *", func() { c.sendMessage() })
	c.cron.Start()

	return c
}

func (c *Command) IsExecute(order command.Order) bool {
	return order.Name == "periodic"
}

func (c *Command) Execute(order command.Order) string {
	return "status: periodic exist"
}

func (c *Command) sendMessage() {
	now := fmt.Sprintf("%d", time.Now().Unix())
	s := c.server
	s.SendPost(database.NewPost(c.sendRoomName, "periodic: I'm alive", "periodiccommand:"+now))
}
