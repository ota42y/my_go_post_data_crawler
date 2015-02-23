package server

import (
	"../../command/command"
	"./../../lib/database"
	"./../../lib/logger"
	"encoding/json"
	"fmt"
	"net/http"
)

// curl -X POST -d "{\"Command\": \"status\",\"Data\":\"d\"}" http://localhost:8080/post

type Post struct {
	Command string
	Data    string
}

type Response struct {
	Result []string
}

type Server struct {
	commands     []command.Command
	logger       *logger.MyLogger
	postDatabase *database.Database
}

func New(logger *logger.MyLogger, postDatabase *database.Database) *Server {
	return &Server{
		logger:       logger,
		postDatabase: postDatabase,
	}
}

func (s *Server) AddCommand(c command.Command) {
	s.commands = append(s.commands, c)
}

func (s *Server) receivePost(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var t Post
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}

	s.executeCommand(t.Command, t.Data, rw)
}

func (s *Server) executeCommand(command string, data string, rw http.ResponseWriter) {
	var res Response

	for _, listener := range s.commands {
		if listener.IsExecute(command) {
			res.Result = append(res.Result, listener.Execute(data))
		}
	}

	enc := json.NewEncoder(rw)
	err := enc.Encode(&res)
	if err != nil {
		fmt.Printf("%d\n", err)
	}
}

func (s *Server) Start() {
	s.logger.LogPrint("server", "start")
	http.HandleFunc("/", s.receivePost)
	http.ListenAndServe(":8080", nil)
}

func (s *Server) SendPost(post *database.Post) {
	s.postDatabase.AddNewPost(post)
	fmt.Println(post)
}
