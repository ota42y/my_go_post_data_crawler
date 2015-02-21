package serverWorker

import (
	"./../../lib/database"
	"./../../lib/logger"
	"./../../lib/server"
	"../../command/status"
	"../../command/periodic"
)

type Worker struct {
	s *server.Server
	logger *logger.MyLogger
	postDatabase       *database.Database
}

func New(logger *logger.MyLogger, postDatabase       *database.Database) *Worker{
	s := server.New(logger, postDatabase)
	return &Worker{
		s: s,
		logger: logger,
		postDatabase: postDatabase,
	}
}

func (w *Worker) Work() {
	s := w.s
	s.AddCommand(status.New())
	s.AddCommand(periodic.New(w.s, w.postDatabase.DefaultRoomName))

	s.Start()
}
