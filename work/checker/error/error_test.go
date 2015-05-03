package error

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"

	"../../../config"
	"../../../lib/database"
)

var testCollections = []string{"test1", "test2", "test3"}

type TestLogger struct {
	B bytes.Buffer
}

func (l *TestLogger) Error(tag string, format string, args ...interface{}) {
	fmt.Fprintf(&l.B, format, args...)
}

func (l *TestLogger) Info(tag string, format string, args ...interface{}) {
	fmt.Fprintf(&l.B, format, args...)
}

func deleteAllData(mongo *config.MongodbDatabase) bool {
	session, err := mgo.Dial(mongo.GetDialURL())
	if err != nil {
		return false
	}
	defer session.Close()

	d := session.DB(mongo.Database)
	names, err := d.CollectionNames()
	if err != nil {
		return false
	}

	for _, name := range names {
		d.C(name).DropCollection()
	}
	return true
}
func createTestData(mongo *config.MongodbDatabase) bool {
	session, err := mgo.Dial(mongo.GetDialURL())
	if err != nil {
		return false
	}
	defer session.Close()

	d := session.DB(mongo.Database)
	for _, name := range testCollections {
		var data = struct {
			Info string
		}{
			Info: "info data",
		}
		d.C(name).Insert(data)

		var er = struct {
			Error string
		}{
			Error: "error data",
		}
		d.C(name).Insert(er)
	}

	return true
}

type TestSender struct {
	IsSaveSuccess bool

	P        []*database.Post
	Complete []*database.Post
}

func (s *TestSender) AddPost(p *database.Post) bool {
	s.P = append(s.P, p)
	return s.IsSaveSuccess
}

func (s *TestSender) AddPosts(posts []*database.Post) int {
	count := 0
	for _, post := range posts {
		if s.AddPost(post) {
			count++
		}
	}
	return count
}

func (s *TestSender) SendComplete(p *database.Post) bool {
	s.Complete = append(s.Complete, p)
	return s.IsSaveSuccess
}

func (s *TestSender) GetLogRoomName() string {
	return "LogRoomName"
}

func (s *TestSender) GetMessageRoomName() string {
	return "MessageRoomName"
}

func (s *TestSender) Reset() {
	s.P = make([]*database.Post, 0)
	s.Complete = make([]*database.Post, 0)
	s.IsSaveSuccess = true
}

func TestErrorChecker(t *testing.T) {
	mongo := &config.MongodbDatabase{
		URL:      "localhost",
		User:     "",
		Pass:     "",
		Database: "work_backup_mongodb_test",
	}

	testLogger := &TestLogger{}
	sender := &TestSender{}
	sender.Reset()

	checker := NewChecker(mongo, sender, testLogger)

	Convey("FindError", t, func() {

		Convey("ErrorExist", func() {
			deleteAllData(mongo)
			createTestData(mongo)
			So(len(checker.getAllErrorLog()), ShouldEqual, 3)
		})

		Convey("ErrorNotExist", func() {
			deleteAllData(mongo)
			So(checker.getAllErrorLog(), ShouldBeEmpty)
		})

	})

	Convey("SendErrorLog", t, func() {

		Convey("ErrorExist", func() {
			deleteAllData(mongo)
			createTestData(mongo)
			logs := checker.getAllErrorLog()
			So(checker.sendErrorLog(logs), ShouldEqual, 3)
			So(len(sender.P), ShouldEqual, 3)
		})

		Convey("ErrorNotExist", func() {
			deleteAllData(mongo)
			sender.Reset()
			logs := checker.getAllErrorLog()
			So(checker.sendErrorLog(logs), ShouldEqual, 0)
			So(len(sender.P), ShouldEqual, 0)
		})

		Convey("DuplicateSendLog", func() {
			deleteAllData(mongo)
			createTestData(mongo)
			sender.Reset()
			logs := checker.getAllErrorLog()
			sender.IsSaveSuccess = false
			So(checker.sendErrorLog(logs), ShouldEqual, 0)
			So(len(sender.P), ShouldEqual, 3)
		})
	})
}
