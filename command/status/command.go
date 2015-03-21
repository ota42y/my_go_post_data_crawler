package status

import (
	"../command"
)

type Status struct {
}

func New() (s *Status) {
	return &Status{}
}

func (s Status) IsExecute(order command.Order) bool {
	return order.Name == "status"
}

func (s Status) Execute(order command.Order) string {
	return "status: server exist"
}
