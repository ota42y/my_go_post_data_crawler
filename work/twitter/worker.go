package twitter

import (
	"./../../lib/database"
	"./../../lib/logger"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/yaml.v2"
)

func makePostDataFromTweet(roomName string, tweet *anaconda.Tweet) (postData *database.Post) {
	account_name := tweet.User.ScreenName
	message := tweet.Text

	return database.NewPost(roomName, "twetter: "+account_name+" "+message, "twetter:"+tweet.IdStr)
}

type Setting struct {
	MongodbUrl     string
	DatabaseName   string
	CollectionName string
	ScreenNames    []string
}

type TwitterAuth struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

type Worker struct {
	checkTwitterIdList []string
	postDatabase       *database.Database
	mongodbData        *MongoDBData
	consumerKey        string
	consumerSecret     string
	accessToken        string
	accessTokenSecret  string
	logger             *logger.MyLogger
}

func New(settingBuf []byte, twitterAuthBuf []byte, db *database.Database, logger *logger.MyLogger) *Worker {
	s := Setting{}
	err := yaml.Unmarshal(settingBuf, &s)
	if err != nil {
		return nil
	}

	auth := TwitterAuth{}
	err = yaml.Unmarshal(twitterAuthBuf, &auth)
	if err != nil {
		return nil
	}

	return &Worker{
		checkTwitterIdList: s.ScreenNames,
		mongodbData:        NewMongoDBData(s.MongodbUrl, s.DatabaseName, s.CollectionName),
		consumerKey:        auth.ConsumerKey,
		consumerSecret:     auth.ConsumerSecret,
		accessToken:        auth.AccessToken,
		accessTokenSecret:  auth.AccessTokenSecret,
		logger:             logger,
		postDatabase:       db,
	}

}

func (worker *Worker) Work() {
	// get unregister tweet
	anaconda.SetConsumerKey(worker.consumerKey)
	anaconda.SetConsumerSecret(worker.consumerSecret)
	api := anaconda.NewTwitterApi(worker.accessToken, worker.accessTokenSecret)

	worker.logger.LogPrint("twitter", "work")

	for _, twitterId := range worker.checkTwitterIdList {
		tweets := getUnRegisterTweet(twitterId, api, worker.mongodbData)

		var posts []*(database.Post)
		for _, tweet := range tweets {
			posts = append(posts, makePostDataFromTweet(worker.postDatabase.DefaultRoomName, &tweet))
		}

		worker.logger.LogPrint("twitter", fmt.Sprintf("add posts %d", len(posts)))
		worker.postDatabase.AddNewPosts(posts)
		registerTweets(tweets, worker.mongodbData)
	}
}
