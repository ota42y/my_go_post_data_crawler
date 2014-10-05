package database

import (
	"time"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"net/url"
)

type Post struct {
	RoomName string
	Message string
	MessageId string
	CreatedAt time.Time
	IsSend bool
}

type Database struct {
	DefaultRoomName string
	driverName string
	dataSourceName string
	sendMessageTableName string
}

func NewDatabase(dataSourceName string, defaultRoomName string) *Database {
	return &Database{
		DefaultRoomName: defaultRoomName,
		driverName: "mysql",
		dataSourceName: dataSourceName,
		sendMessageTableName: "send_message",
	}
}

func (post *Post) GetUrlValue() (postData *url.Values){
	return &url.Values{"room_name": []string{post.RoomName}, "message": []string{post.Message}}
}

// message_idで登録済みかどうかをチェックできる
// 使わなくても問題は無い
func NewPost(room_name string, message string, message_id string) *Post{
	return &Post{
		RoomName: room_name,
		Message: message,
		MessageId: message_id,
		CreatedAt: time.Now(),
		IsSend: false,
	}
}

// message_idを持ったpostが存在するかどうか
func (database *Database) GetPost(message_id string) *Post{
	db, err := sql.Open(database.driverName, database.dataSourceName)

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Prepare statement for reading data
	stmtOut, err := db.Prepare(fmt.Sprintf("SELECT room_name, message, message_id, created_at, is_send FROM %s WHERE message_id = ? LIMIT 1", database.sendMessageTableName))
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	no_send_posts := getPostFromDatabase(stmtOut.QueryRow(message_id))
	if no_send_posts == nil {
		return nil
	}
	return no_send_posts
}

func (database *Database) SendPost(post *Post) (success bool) {
	db, err := sql.Open(database.driverName, database.dataSourceName)

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare(fmt.Sprintf("UPDATE %s SET is_send = 1 WHERE ( message_id = ?)", database.sendMessageTableName))
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(post.MessageId) // Insert tuples (i, i^2)
	return err != nil
}

func (database *Database) AddNewPosts(posts []*Post) (is_success bool){
	db, err := sql.Open(database.driverName, database.dataSourceName)

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (room_name, message, message_id, created_at, is_send) VALUES (?, ?, ?, ?, ?)", database.sendMessageTableName))
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	for _, post := range posts{
		if(database.GetPost(post.MessageId) == nil){
			fmt.Println("save %s", post.Message)
			_, err = stmtIns.Exec(post.RoomName, post.Message, post.MessageId, post.CreatedAt, post.IsSend)
			if err != nil{
				return false

			}
		}
	}

	return true
}

func (database *Database) GetNoSendPosts(limit int) (posts []*Post){
	db, err := sql.Open(database.driverName, database.dataSourceName)

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT room_name, message, message_id, created_at, is_send FROM %s WHERE is_send = false LIMIT %d", database.sendMessageTableName, limit))
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	no_send_posts := make([]*Post, 0)

	for rows.Next() {
		var room_name string
		var message string
		var message_id string
		var created_at time.Time
		var is_send bool
		if err := rows.Scan(&room_name, &message, &message_id, &created_at, &is_send); err != nil {
		}

		post := NewPost(room_name, message, message_id)
		post.IsSend = is_send
		post.CreatedAt = created_at
		no_send_posts = append(no_send_posts, post)
	}

	return no_send_posts
}

func getPostFromDatabase(row *sql.Row) (post *Post){
	var room_name string
	var message string
	var message_id string
	var created_at time.Time
	var is_send bool
	if err := row.Scan(&room_name, &message, &message_id, &created_at, &is_send); err != nil {
		return nil
	}

	postData := NewPost(room_name, message, message_id)
	postData.IsSend = is_send
	postData.CreatedAt = created_at
	return postData
}
