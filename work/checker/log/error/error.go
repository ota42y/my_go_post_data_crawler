package error

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../../../../config"
	"../../../../lib/logger"
	"../../../../lib/post"

	"../notification"
)

// NewLogCollector return Notification with LogCollector
func NewLogCollector(mongo *config.MongodbDatabase, sender post.Sender, logger logger.Logger) *notification.Notification {
	return notification.NewNotification(&LogCollector{}, mongo, sender, logger)
}

// LogCollector collect error logs
type LogCollector struct {
}

// GetLogs return all error logs from mongodb
func (c *LogCollector) GetLogs(collectionNames []string, database *mgo.Database) []notification.LogBase {
	var errorLogs []notification.LogBase

	for _, name := range collectionNames {
		var logs []Log

		q := database.C(name).Find(bson.M{"error": bson.M{"$exists": true}})
		q.All(&logs)

		for _, l := range logs {
			errorLogs = append(errorLogs, &l)
		}
	}
	return errorLogs
}

// CollectorName return "LogCollector"
func (c *LogCollector) CollectorName() string {
	return "LogCollector"
}

// Log is log data from mongodb
type Log struct {
	Error string
	Time  time.Time
}

// GetMessage return info log^s message
func (l *Log) GetMessage() string {
	return l.Error
}

// GetTime return log time
func (l *Log) GetTime() time.Time {
	return l.Time
}
