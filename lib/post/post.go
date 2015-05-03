package post

import (
	"../database"
)

// Sender send message to chat
type Sender interface {
	AddPost(p *database.Post) bool
	AddPosts(posts []*database.Post) int
	SendComplete(p *database.Post) bool
	GetLogRoomName() string
	GetMessageRoomName() string
}

/* need refactoring
type Config struct {
    LogRoomName string
    MessageRoomName string
    DataSourceName string
}

func NewSender(c *Config) *Sender {
    d := database.NewDatabase(c.DataSourceName, c.LogRoomName, c.LogRoomName)

    return &MysqlSender{
        *database.Database : d,
    }
}
*/

// NewSender return Sender
func NewSender(d *database.Database) Sender {
	return &MysqlSender{
		db: d,
	}
}

// MysqlSender use mysql for checking duplicate management
type MysqlSender struct {
	db *database.Database
}

// AddPost register new send message
// if already exist or error, return false
func (s *MysqlSender) AddPost(p *database.Post) bool {
	return s.db.AddNewPost(p)
}

// AddPosts register new send messages
// return register success posts num
func (s *MysqlSender) AddPosts(posts []*database.Post) int {
	count := 0
	for _, post := range posts {
		if s.AddPost(post) {
			count++
		}
	}
	return count
}

// SendComplete update post^s send flag
func (s *MysqlSender) SendComplete(p *database.Post) bool {
	return s.db.SendPost(p)
}

// GetLogRoomName return room name which not notification room
func (s *MysqlSender) GetLogRoomName() string {
	return s.db.LogRoomName
}

// GetMessageRoomName return room name which notification room
func (s *MysqlSender) GetMessageRoomName() string {
	return s.db.DefaultRoomName
}
