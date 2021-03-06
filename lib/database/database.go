package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"net/url"
	"time"
)

type Post struct {
	Id        int64
	RoomName  string
	Message   string
	MessageId string
	CreatedAt time.Time
	IsSend    bool
}

type Database struct {
	LogRoomName          string
	DefaultRoomName      string
	sendMessageTableName string
	db                   gorm.DB
}

func NewDatabase(dataSourceName string, defaultRoomName string, logRoomName string) *Database {
	db, err := gorm.Open("mysql", dataSourceName)

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	db.DB()
	db.AutoMigrate(&Post{})
	db.Model(&Post{}).AddUniqueIndex("idx_message_id", "message_id")

	return &Database{
		LogRoomName:          logRoomName,
		DefaultRoomName:      defaultRoomName,
		sendMessageTableName: "send_message",
		db:                   db,
	}
}

func (post *Post) GetUrlValue() (postData *url.Values) {
	return &url.Values{"room_name": []string{post.RoomName}, "message": []string{post.Message}}
}

// message_idで登録済みかどうかをチェックできる
// 使わなくても問題は無い
func NewPost(room_name string, message string, message_id string) *Post {
	return &Post{
		RoomName:  room_name,
		Message:   message,
		MessageId: message_id,
		CreatedAt: time.Now(),
		IsSend:    false,
	}
}

// message_idを持ったpostが存在するかどうか
func (database *Database) GetPost(message_id string) *Post {
	var post Post
	err := database.db.Where("message_id = ?", message_id).First(&post).Error
	if err == gorm.RecordNotFound {
		return nil
	}

	return &post
}

func (database *Database) SendPost(post *Post) (success bool) {
	post.IsSend = true
	err := database.db.Save(post)

	return err != nil
}

func (database *Database) AddNewPost(post *Post) (is_success bool) {
	if database.GetPost(post.MessageId) == nil {
		err := database.db.Save(post)
		if err != nil {
			return false
		}
		fmt.Println("add post %v", post)
	}
	return true
}

func (database *Database) AddNewPosts(posts []*Post) (is_success bool) {
	for _, post := range posts {
		if !database.AddNewPost(post) {
			return false
		}
	}

	return true
}

func (database *Database) GetNoSendPosts(limit int) (posts []*Post) {
	database.db.Where("is_send = false").Limit(limit).Find(&posts)
	return posts
}
