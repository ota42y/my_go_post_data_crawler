package twitter

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"../../../test"
	"github.com/ChimeraCoder/anaconda"
)

func TestSomething(t *testing.T) {
	logger := &test.Logger{}
	sender := &test.Sender{}

	mongoConfig := test.NewTestMongodbConfig()
	setting := &Setting{
		MongodbURL:     mongoConfig.URL,
		DatabaseName:   mongoConfig.Database,
		CollectionName: "twitter_test",
	}

	twitter := &Twitter{
		setting: setting,
		l:       logger,
		sender:  sender,
		api:     nil,
	}

	testTweet := anaconda.Tweet{
		IdStr: "123456",
		Text:  "test text",
	}

	Convey("filterUnRegistTweet", t, func() {

		Convey("tweet exist", func() {
			test.DeleteAllData(mongoConfig)

			var tw []anaconda.Tweet
			tw = append(tw, testTweet)

			// regist
			twitter.registTweet(&tw)
			So(twitter.filterUnRegistTweet(&tw), ShouldBeEmpty)
		})

		Convey("tweet not exit", func() {
			test.DeleteAllData(mongoConfig)

			var tw []anaconda.Tweet
			tw = append(tw, testTweet)

			So(twitter.filterUnRegistTweet(&tw), ShouldResemble, &tw)
		})

		Convey("no tweet", func() {
			test.DeleteAllData(mongoConfig)

			var tw []anaconda.Tweet
			So(twitter.filterUnRegistTweet(&tw), ShouldBeEmpty)
		})

	})

	Convey("sendTweets", t, func() {

		Convey("tweet exist", func() {
			test.DeleteAllData(mongoConfig)
			sender.Reset()
			sender.IsSaveSuccess = false

			var tw []anaconda.Tweet
			tw = append(tw, testTweet)

			So(twitter.sendTweets(&tw), ShouldBeEmpty)
			So(sender.P, ShouldNotBeEmpty)
		})

		Convey("tweet not exit", func() {
			test.DeleteAllData(mongoConfig)
			sender.Reset()

			var tw []anaconda.Tweet
			tw = append(tw, testTweet)

			So(twitter.sendTweets(&tw), ShouldNotBeEmpty)
			So(sender.P, ShouldNotBeEmpty)
		})

		Convey("no tweet", func() {
			test.DeleteAllData(mongoConfig)
			sender.Reset()

			var tw []anaconda.Tweet
			So(twitter.sendTweets(&tw), ShouldBeEmpty)
			So(sender.P, ShouldBeEmpty)
		})

	})

	Convey("registTweet", t, func() {

		Convey("regist", func() {
			test.DeleteAllData(mongoConfig)
			var tw []anaconda.Tweet
			tw = append(tw, testTweet)

			So(twitter.filterUnRegistTweet(&tw), ShouldNotBeEmpty)
			twitter.registTweet(&tw)
			So(twitter.filterUnRegistTweet(&tw), ShouldBeEmpty)
		})

		Convey("no tweet", func() {
			test.DeleteAllData(mongoConfig)
			var tw []anaconda.Tweet
			tw = append(tw, testTweet)

			twitter.registTweet(&tw)
		})

	})
}
