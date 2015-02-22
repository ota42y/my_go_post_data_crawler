package periodic

import (
	"./../../lib/database"
	"./../../lib/server"
	"github.com/robfig/cron"
	"time"
	"fmt"
)

type Command struct{
	server *server.Server
	sendRoomName string
	cron *cron.Cron
}

func New(server *server.Server, sendRoomName string) (c *Command){
	c = &Command{
		server: server,
		sendRoomName: sendRoomName,
		cron: cron.New(),
	}

	c.cron.AddFunc("0 */30 7-23 * * *", func() { c.sendMessage() })
	c.cron.Start()

	return c
}

func (c *Command) IsExecute(command string) bool{
	return command == "periodic"
}

func (c *Command) Execute(data string) string{
	return "status: periodic exist"
}

func (c *Command) sendMessage() {
	now := fmt.Sprintf("%d", time.Now().Unix())
	s := c.server
	s.SendPost(database.NewPost(c.sendRoomName, "periodic: I'm alive", "periodiccommand:" + now))
}
