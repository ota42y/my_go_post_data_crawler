package error

import (
	"time"

	"../../../config"
	"../../../lib/database"
	"../../../lib/logger"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Checker is error log checker
type Checker struct {
	logger logger.Logger
	mongo  *config.MongodbDatabase
	db     *database.Database
}

// NewChecker return Checker
func NewChecker(mongo *config.MongodbDatabase, db *database.Database, logger logger.Logger) *Checker {
	return &Checker{
		logger: logger,
		mongo:  mongo,
		db:     db,
	}
}

// Log is struct for error log by fluentd
type Log struct {
	Error string
	Time  time.Time
}

// Execute is check error log in database and send log
func (c *Checker) Execute() {
	eLogs := c.getAllErrorLog()
	c.sendErrorLog(eLogs)
}

func (c *Checker) getAllErrorLog() []Log {
	c.logger.Info("errorLogChecker", "getAllErrorLog")

	var eLogs []Log

	session, err := mgo.Dial(c.mongo.GetDialURL())
	if err != nil {
		return eLogs
	}
	defer session.Close()

	d := session.DB(c.mongo.Database)
	names, err := d.CollectionNames()
	if err != nil {
		return eLogs
	}

	for _, name := range names {
		var logs []Log

		q := d.C(name).Find(bson.M{"error": bson.M{"$exists": true}})
		q.All(&logs)

		c.logger.Info("errorLogChecker", "get %d logs in %d", name, len(logs))

		eLogs = append(eLogs, logs...)
	}

	c.logger.Info("errorLogChecker", "get error logs %d", len(eLogs))
	return eLogs
}

// return send error logs, if already exist, not send
func (c *Checker) sendErrorLog(eLogs []Log) int {
	c.logger.Info("errorLogChecker", "sendErrorLog %d", len(eLogs))
	count := 0

	for _, l := range eLogs {
		p := database.NewPost(c.db.LogRoomName, l.Error, "errorLog:"+l.Time.String())
		if c.db.AddNewPost(p) {
			count++
		}
	}

	c.logger.Info("errorLogChecker", "sendErrorLog %d errors send", count)
	return count
}
