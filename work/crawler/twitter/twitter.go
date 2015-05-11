package twitter

import (
	"fmt"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../../../lib/database"
)

var logName = "twitter"

func (t *Twitter) twitterCrawling() {
	for _, screenName := range t.setting.ScreenNames {
		t.checkNewTweet(screenName)
	}
}

func (t *Twitter) checkNewTweet(screenName string) {
	tweets := t.getTweet(screenName)
	unRegisterTweets := t.filterUnRegistTweet(tweets)
	complete := t.sendTweets(unRegisterTweets)
	t.registTweet(complete)
}

func (t *Twitter) getTweet(screenName string) *[]anaconda.Tweet {
	v := url.Values{}
	v.Set("screen_name", screenName)

	tweets, err := t.api.GetUserTimeline(v)
	if err != nil {
		t.l.Error(logName, "api.GetUserTimenile(%s) error : %s", screenName, err.Error())
		var tw []anaconda.Tweet
		return &tw
	}
	t.l.Info(logName, "get %d tweets", len(tweets))

	return &tweets
}

func (t *Twitter) filterUnRegistTweet(tweets *[]anaconda.Tweet) *[]anaconda.Tweet {
	var unregister []anaconda.Tweet

	session, err := mgo.Dial(t.setting.MongodbURL)
	if err != nil {
		t.l.Error(logName, "filterUnRegistTweet mongodb dail %s but error", t.setting.MongodbURL, err.Error())
		return &unregister
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(t.setting.DatabaseName).C(t.setting.CollectionName)
	for _, tweet := range *(tweets) {
		count, _ := c.Find(bson.M{"id": tweet.Id}).Count()
		if count == 0 {
			unregister = append(unregister, tweet)
		}
	}

	t.l.Info(logName, "unregist %d tweets", len(unregister))
	return &unregister
}

func (t *Twitter) sendTweets(tweets *[]anaconda.Tweet) *[]anaconda.Tweet {
	var complete []anaconda.Tweet

	for _, tweet := range *(tweets) {
		p := makePostDataFromTweet(t.sender.GetMessageRoomName(), &tweet)

		if t.sender.AddPost(p) {
			t.l.Info(logName, "send tweet : %s", tweet.Text)
			complete = append(complete, tweet)
		}
	}

	t.l.Info(logName, "send %s tweet", len(complete))
	return &complete
}

func makePostDataFromTweet(roomName string, tweet *anaconda.Tweet) (postData *database.Post) {
	message := fmt.Sprintf("twitter: %s %s", tweet.User.ScreenName, tweet.Text)
	id := fmt.Sprintf("twitter:%s", tweet.IdStr)

	return database.NewPost(roomName, message, id)
}

func (t *Twitter) registTweet(tweets *[]anaconda.Tweet) {
	session, err := mgo.Dial(t.setting.MongodbURL)
	if err != nil {
		t.l.Error(logName, "registTweet mongodb dail %s but error", t.setting.MongodbURL, err.Error())
		return
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(t.setting.DatabaseName).C(t.setting.CollectionName)
	for _, tweet := range *(tweets) {
		err = c.Insert(&tweet)
		if err != nil {
			t.l.Error(logName, "registTweet insert %s error %s", tweet.Text, err.Error())
		}
	}
}
