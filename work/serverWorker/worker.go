package serverWorker

import (
	"../../command/periodic"
	"../../command/pomodoro"
	"../../util"
	"./../../command/status"
	"./../../lib/database"
	"./../../lib/logger"
	"./../../lib/server"
)

type Worker struct {
	s            *server.Server
	logger       *logger.MyLogger
	postDatabase *database.Database
	settingHome  string
}

func New(logger *logger.MyLogger, postDatabase *database.Database, settingHome string) *Worker {
	s := server.New(logger, postDatabase)
	return &Worker{
		s:            s,
		logger:       logger,
		postDatabase: postDatabase,
		settingHome:  settingHome,
	}
}

func (w *Worker) Work() {
	s := w.s
	s.AddCommand(status.New())
	s.AddCommand(periodic.New(w.s, w.postDatabase.LogRoomName))
	s.AddCommand(pomodoro.New(w.s, w.postDatabase.DefaultRoomName, util.LoadFile(w.settingHome+"/tumblr.yml")))

	s.Start()
}
