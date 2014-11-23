package twitter

import (
	"fmt"
	"./../../lib/database"
	"github.com/ChimeraCoder/anaconda"
	"./../../lib/logger"
)

func makePostDataFromTweet(roomName string, tweet *anaconda.Tweet) (postData *database.Post){
	account_name := tweet.User.ScreenName
	message := tweet.Text

	return database.NewPost(roomName, "twetter: " + account_name + " " + message , "twetter:" + tweet.IdStr)
}

type Worker struct {
	checkTwitterIdList []string
	postDatabase *database.Database
	mongodbData *MongoDBData
	consumerKey string
	consumerSecret string
	accessToken string
	accessTokenSecret string
	logger *logger.MyLogger
}

func NewWorkerFromMap(sendDataConfig map[interface{}]interface{}, twitterAuthConfig map[interface{}]interface{}, postDatabase *database.Database, logger *logger.MyLogger) (* Worker) {
	mongodbUrl := sendDataConfig["mongodbUrl"].(string)


	twitterData := sendDataConfig["twitter"].(map[interface{}]interface{})
	databaseName := twitterData["databaseName"].(string)
	collectionName := twitterData["collectionName"].(string)

	checkTwitterIdList := make([]string, 0)
	screenNames := twitterData["screenNames"].([]interface{})

	for _, screenName := range screenNames{
		checkTwitterIdList = append(checkTwitterIdList, screenName.(string))
	}

	consumerKey := twitterAuthConfig["consumerKey"].(string)
	consumerSecret := twitterAuthConfig["consumerSecret"].(string)
	accessToken := twitterAuthConfig["accessToken"].(string)
	accessTokenSecret := twitterAuthConfig["accessTokenSecret"].(string)

	return &Worker{
		checkTwitterIdList: checkTwitterIdList,
		postDatabase: postDatabase,
		mongodbData: NewMongoDBData(mongodbUrl, databaseName, collectionName),
		consumerKey: consumerKey,
		consumerSecret: consumerSecret,
		accessToken: accessToken,
		accessTokenSecret: accessTokenSecret,
		logger: logger,
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
