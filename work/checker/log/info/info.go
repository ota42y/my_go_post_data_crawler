package info

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

// LogCollector collect info logs
type LogCollector struct {
}

// GetLogs return all info logs from mongodb
func (c *LogCollector) GetLogs(collectionNames []string, database *mgo.Database) []notification.LogBase {
	var infoLogs []notification.LogBase

	for _, name := range collectionNames {
		var logs []Log

		q := database.C(name).Find(bson.M{"info": bson.M{"$exists": true}})
		q.All(&logs)

		for _, l := range logs {
			infoLogs = append(infoLogs, &l)
		}
	}
	return infoLogs
}

// CollectorName return "LogCollector"
func (c *LogCollector) CollectorName() string {
	return "LogCollector"
}

// Log is log data from mongodb
type Log struct {
	Info string
	Time time.Time
}

// GetMessage return info log^s message
func (l *Log) GetMessage() string {
	return l.Info
}

// GetTime return log time
func (l *Log) GetTime() time.Time {
	return l.Time
}
