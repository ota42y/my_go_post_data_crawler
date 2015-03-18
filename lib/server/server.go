package server

import (
	"../../command/command"
	"./../../lib/database"
	"./../../lib/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// curl -X POST -d "{\"Command\": \"status\",\"Data\":\"d\"}" http://localhost:8080/post

var commandRegExp, _ = regexp.Compile("^(.+)( (.*))")

type PostData struct {
	Message string
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
	var post PostData

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(rw, "read error %s\n", err)
		return
	}
	fmt.Println("get ", string(b))

	err = json.Unmarshal(b, &post)
	if err != nil {
		fmt.Fprintf(rw, "unmarshal error %s\n", err)
		return
	}

	match := commandRegExp.FindSubmatch([]byte(post.Message))
	if len(match) == 0 {
		fmt.Fprintf(rw, "regexp error %s\n", post.Message)
		return
	}

	cmd := string(match[1])

	data := ""
	if 2 < len(match) {
		data = string(match[3])
	}

	s.executeCommand(cmd, data, rw)
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
		fmt.Fprintf(rw, "resopnse error %s\n", err)
	}
}

func (s *Server) Start() {
	s.logger.LogPrint("server", "start")
	http.HandleFunc("/", s.receivePost)
	http.ListenAndServe(":8080", nil)
}

func (s *Server) SendPost(post *database.Post) {
	s.postDatabase.AddNewPost(post)
}

func (s *Server) LogPrint(tag string, message string) {
	s.logger.LogPrint(tag, message)
}
