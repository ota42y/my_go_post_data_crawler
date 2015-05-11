package test

import (
	"gopkg.in/mgo.v2"

	"../config"
)

// NewTestMongodbConfig return test mongodb data
func NewTestMongodbConfig() *config.MongodbDatabase {
	return &config.MongodbDatabase{
		URL:      "localhost",
		User:     "",
		Pass:     "",
		Database: "work_backup_mongodb_test",
	}
}

// DeleteAllData delete all data in database
func DeleteAllData(mongo *config.MongodbDatabase) bool {
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
