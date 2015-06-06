package notification

import (
	"time"

	"../../../../config"
	"../../../../lib/database"
	"../../../../lib/logger"
	"../../../../lib/post"

	"gopkg.in/mgo.v2"
)

// LogBase return log message and log time interface
type LogBase interface {
	GetMessage() string
	GetTime() time.Time
}

// LogCollector collect log in mongodb interface
type LogCollector interface {
	GetLogs(collectionNames []string, database *mgo.Database) []LogBase
	CollectorName() string
}

// Notification collect logs and send notification
type Notification struct {
	col    LogCollector
	logger logger.Logger
	mongo  *config.MongodbDatabase
	sender post.Sender
}

// NewNotification return Notification
func NewNotification(col LogCollector, mongo *config.MongodbDatabase, sender post.Sender, logger logger.Logger) *Notification {
	return &Notification{
		col:    col,
		logger: logger,
		mongo:  mongo,
		sender: sender,
	}
}

// Execute is check error log in database and send log
func (n *Notification) Execute() {
	l := n.GetLog()
	n.SendLogs(l)
}

// GetLog return logs from database
func (n *Notification) GetLog() []LogBase {
	n.logger.Debug(n.col.CollectorName(), "getLog")

	session, err := mgo.Dial(n.mongo.GetDialURL())
	if err != nil {
		n.logger.Error(n.col.CollectorName(), err.Error())
		return make([]LogBase, 0)
	}
	defer session.Close()

	d := session.DB(n.mongo.Database)
	names, err := d.CollectionNames()
	if err != nil {
		n.logger.Error(n.col.CollectorName(), err.Error())
		return make([]LogBase, 0)
	}

	logs := n.col.GetLogs(names, d)
	n.logger.Debug(n.col.CollectorName(), "get logs %d", len(logs))
	return logs
}

// SendLogs return send logs, if already exist, not send
func (n *Notification) SendLogs(logs []LogBase) int {
	n.logger.Debug(n.col.CollectorName(), "send %d logs", len(logs))
	count := 0
	name := n.col.CollectorName()

	for _, l := range logs {
		p := database.NewPost(n.sender.GetMessageRoomName(), l.GetMessage(), name+l.GetTime().String())
		if n.sender.AddPost(p) {
			count++
		}
	}

	n.logger.Debug(n.col.CollectorName(), "send %d logs end", count)
	return count
}
