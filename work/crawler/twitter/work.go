package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/yaml.v2"

	"./../../../lib/logger"
	"./../../../lib/post"
)

// Twitter is twitter crawling struct
type Twitter struct {
	setting *Setting

	sender post.Sender
	l      logger.Logger
	api    *anaconda.TwitterApi
}

// Setting is Twitter setting struct
// ScreenNames is crawling user screen name list
type Setting struct {
	MongodbURL     string
	DatabaseName   string
	CollectionName string
	ScreenNames    []string
}

// Auth is twitter oauth setting struct
type Auth struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// NewTwitter return Twitter struct
func NewTwitter(settingBuf []byte, authBuf []byte, sender post.Sender, l logger.Logger) *Twitter {
	s := Setting{}
	err := yaml.Unmarshal(settingBuf, &s)
	if err != nil {
		return nil
	}

	auth := Auth{}
	err = yaml.Unmarshal(authBuf, &auth)
	if err != nil {
		return nil
	}

	anaconda.SetConsumerKey(auth.ConsumerKey)
	anaconda.SetConsumerSecret(auth.ConsumerSecret)
	api := anaconda.NewTwitterApi(auth.AccessToken, auth.AccessTokenSecret)

	t := &Twitter{
		setting: &s,
		sender:  sender,
		l:       l,
		api:     api,
	}

	return t
}

// Execute is starting crawler
func (t *Twitter) Execute() {
	t.l.Debug(logName, "execute")
	t.twitterCrawling()
	t.l.Debug(logName, "execute end")
}
