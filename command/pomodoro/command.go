package pomodoro

import (
	"../command"
	"./../../lib/database"
	"./../../lib/server"
	"fmt"
	"github.com/mrjones/oauth"
	"github.com/ota42y/go-tumblr/tumblr"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
	"math/rand"
	"strconv"
	"time"
)

type Setting struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	BlogUrl           string
}

type Command struct {
	server       *server.Server
	sendRoomName string
	cron         *cron.Cron
	isStart      bool
	blog         *tumblr.BlogApi
}

func New(server *server.Server, sendRoomName string, setting []byte) (c *Command) {
	s := Setting{}
	err := yaml.Unmarshal(setting, &s)
	if err != nil {
		return nil
	}

	t := tumblr.New(s.ConsumerKey, s.ConsumerSecret)
	token := &oauth.AccessToken{
		Token:  s.AccessToken,
		Secret: s.AccessTokenSecret,
	}

	blogApi := t.NewBlogApi(s.BlogUrl, token)

	c = &Command{
		server:       server,
		sendRoomName: sendRoomName,
		cron:         cron.New(),
		isStart:      false,
		blog:         blogApi,
	}

	c.cron.AddFunc("0 */30 * * * *", func() { c.sendMessage() })
	return c
}

func (c *Command) IsExecute(order command.Order) bool {
	return order.Name == "pomodoro"
}

func (c *Command) Execute(order command.Order) string {
	if c.isStart {
		c.cron.Stop()
		c.isStart = false
		c.sendRoomName = order.Room
		return "pomodoro: stop"
	} else {
		c.cron.Start()
		c.isStart = true
		c.sendRoomName = order.Room
		return "pomodoro: start"
	}
}

func (c *Command) sendMessage() {
	// 投稿数をとってくる
	_, b, err := c.blog.Info()
	if err != nil {
		message := fmt.Sprintf("blog.Info error %v", err)
		c.server.LogPrint("pomodoro", message)
		return
	}

	// 投稿数の取得
	postNum := rand.Int() % b.Posts

	// ランダムに投稿を取ってくる
	params := make(map[string]string)
	params["offset"] = strconv.Itoa(postNum)
	params["limit"] = "1"
	_, posts, err := c.blog.Photo(&params)
	if err != nil {
		message := fmt.Sprintf("blog.Photos error %v", err)
		c.server.LogPrint("pomodoro", message)
		return
	}

	if len(*posts) == 0 {
		c.server.LogPrint("pomodoro", "no posts")
		return
	}

	// 画像のURLを取り出す
	post := (*posts)[0]
	if len(post.Photos) == 0 {
		c.server.LogPrint("pomodoro", "no photos")
		return
	}
	if len(post.Photos[0].AltSizes) == 0 {
		c.server.LogPrint("pomodoro", "no sizes")
		return
	}
	url := post.Photos[0].AltSizes[0].Url

	// 画像付きで進捗を聞く
	now := fmt.Sprintf("%d", time.Now().Unix())
	s := c.server
	s.SendPost(database.NewPost(c.sendRoomName, "pomodoro: "+url, "pomodorocommand:img"+now))
	s.SendPost(database.NewPost(c.sendRoomName, "pomodoro: 進捗どうですか？", "pomodorocommand:"+now))
}
