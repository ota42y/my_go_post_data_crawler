package test

import (
	"../lib/database"
)

// Sender is stub for lib/post.Sender
type Sender struct {
	IsSaveSuccess bool

	P        []*database.Post
	Complete []*database.Post
}

// AddPost add TestSender.P and return IsSaveSuccess
func (s *Sender) AddPost(p *database.Post) bool {
	s.P = append(s.P, p)
	return s.IsSaveSuccess
}

// AddPosts add TestSender.P
func (s *Sender) AddPosts(posts []*database.Post) int {
	count := 0
	for _, post := range posts {
		if s.AddPost(post) {
			count++
		}
	}
	return count
}

// SendComplete add TestSender.Complete and return IsSaveSuccess
func (s *Sender) SendComplete(p *database.Post) bool {
	s.Complete = append(s.Complete, p)
	return s.IsSaveSuccess
}

// GetLogRoomName return LogRoomName
func (s *Sender) GetLogRoomName() string {
	return "LogRoomName"
}

// GetMessageRoomName return MessageRoomName
func (s *Sender) GetMessageRoomName() string {
	return "MessageRoomName"
}

// Reset reset P, Complete, IsSaveSuccess
func (s *Sender) Reset() {
	s.P = make([]*database.Post, 0)
	s.Complete = make([]*database.Post, 0)
	s.IsSaveSuccess = true
}
