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

func TestErrorChecker(t *testing.T) {
	mongo := &config.MongodbDatabase{
		URL:      "localhost",
		User:     "",
		Pass:     "",
		Database: "work_backup_mongodb_test",
	}

	db := &database.Database{}
	testLogger := &TestLogger{}

	checker := NewChecker(mongo, db, testLogger)

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

		Convey("ErrorExist", nil)
		/* need refactoring
		   func() {
		       deleteAllData(mongo)
		       createTestData(mongo)
		       logs := checker.getAllErrorLog()
		       So(checker.sendErrorLog(logs), ShouldEqual, 3)
		   })
		*/

		Convey("ErrorNotExist", nil)
		/* need refactoring
		   func() {
		       deleteAllData(mongo)
		       logs := checker.getAllErrorLog()
		       So(checker.sendErrorLog(logs), ShouldEqual, 0)
		   })
		*/

		Convey("DuplicateSendLog", nil)
		/* need refactoring
		   func() {
		       deleteAllData(mongo)
		       createTestData(mongo)
		       logs := checker.getAllErrorLog()
		       So(logs, ShouldNotBeEmpty)

		       checker.sendErrorLog(logs[1:])
		       So(checker.sendErrorLog(logs), ShouldEqual, 2)
		   })
		*/
	})
}
