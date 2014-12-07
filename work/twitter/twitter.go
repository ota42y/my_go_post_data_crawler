package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/url"
)

type MongoDBData struct {
	url            string
	databaseName   string
	collectionName string
}

func NewMongoDBData(url string, databaseName string, collectionName string) *MongoDBData {
	return &MongoDBData{
		url:            url,
		databaseName:   databaseName,
		collectionName: collectionName,
	}
}

func registerTweets(tweets []anaconda.Tweet, mongodb *MongoDBData) {
	session, err := mgo.Dial(mongodb.url)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(mongodb.databaseName).C(mongodb.collectionName)

	for _, tweet := range tweets {
		err = c.Insert(&tweet)
		if err != nil {
			panic(err)
		}
	}
}

func filterUnRegisterTweets(tweets []anaconda.Tweet, mongodb *MongoDBData) []anaconda.Tweet {
	filteredTweets := []anaconda.Tweet{}

	session, err := mgo.Dial(mongodb.url)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(mongodb.databaseName).C(mongodb.collectionName)

	for _, tweet := range tweets {
		count, _ := c.Find(bson.M{"id": tweet.Id}).Count()
		if count == 0 {
			filteredTweets = append(filteredTweets, tweet)
		}
	}

	return filteredTweets
}

func getTweets(screen_name string, api *anaconda.TwitterApi) []anaconda.Tweet {
	v := url.Values{}
	v.Set("screen_name", screen_name)

	tweets, err := api.GetUserTimeline(v)
	if err != nil {
		panic(err)
	}
	return tweets
}

func getUnRegisterTweet(screen_name string, api *anaconda.TwitterApi, mongodb *MongoDBData) []anaconda.Tweet {
	tweets := getTweets(screen_name, api)
	return filterUnRegisterTweets(tweets, mongodb)
}
