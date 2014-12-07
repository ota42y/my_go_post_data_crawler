package sendHubot

import (
	"net/http"
	"net/url"
)

type Server struct {
	PostPath string
}

func NewServer(postPath string) *Server {
	return &Server{
		PostPath: postPath,
	}
}

func (server *Server) SendData(post_data *url.Values) (success bool) {
	resp, _ := http.PostForm(
		server.PostPath,
		*post_data,
	)

	if resp == nil {
		return false
	}

	return resp.StatusCode == 200
}
