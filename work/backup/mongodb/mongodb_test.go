package mongodb

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"

	"../../../config"
	"../../../test"
)

var testCollections = []string{"test1", "test2", "test3"}

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
			Data string
		}{
			Data: "test",
		}

		d.C(name).Insert(data)
	}

	return true
}

func checkExpireIndexExist(mongo *config.MongodbDatabase) bool {
	session, err := mgo.Dial(mongo.GetDialURL())
	if err != nil {
		return false
	}
	defer session.Close()

	d := session.DB(mongo.Database)
	for _, name := range testCollections {
		indexes, err := d.C(name).Indexes()
		if err != nil {
			return false
		}

		So(indexes, ShouldNotBeEmpty)
		for _, index := range indexes {
			// skip id index
			if index.Key[0] == "_id" {
				continue
			}

			So(index.ExpireAfter, ShouldEqual, time.Hour*24*7)
		}
	}

	return true
}

func TestSomething(t *testing.T) {
	c := &config.MongodbBackup{
		ArchivePath: "./backup",
		Workspace:   "./dump",
	}

	mongo := &config.MongodbDatabase{
		URL:      "localhost",
		User:     "",
		Pass:     "",
		Database: "work_backup_mongodb_test",
	}

	testLogger := &test.Logger{}

	m := NewMongodb(c, mongo, testLogger)
	Convey("get collection name from mongodb", t, func() {
		Convey("data exist", func() {
			Convey("return collections", func() {
				So(deleteAllData(mongo), ShouldBeTrue)
				So(createTestData(mongo), ShouldBeTrue)
				So(m.getAllCollectionNames(), ShouldResemble, testCollections)
			})
		})

		Convey("data no exist", func() {

			Convey("return empty list", func() {
				So(deleteAllData(mongo), ShouldBeTrue)
				So(m.getAllCollectionNames(), ShouldBeEmpty)
			})
		})

	})

	Convey("make index", t, func() {
		Convey("new index", func() {
			So(deleteAllData(mongo), ShouldBeTrue)
			So(createTestData(mongo), ShouldBeTrue)
			m.createExpireIndex()
			So(checkExpireIndexExist(mongo), ShouldBeTrue)
		})

		Convey("already index", func() {
			So(deleteAllData(mongo), ShouldBeTrue)
			So(createTestData(mongo), ShouldBeTrue)
			m.createExpireIndex()
			m.createExpireIndex()
			So(checkExpireIndexExist(mongo), ShouldBeTrue)
		})
	})
}
