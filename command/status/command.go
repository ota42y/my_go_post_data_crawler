package status

type Status struct{
}

func New() (s *Status){
	return &Status{}
}

func (s Status) IsExecute(command string) bool{
	return command == "status"
}

func (s Status) Execute(data string) string{
	return "status: server exist"
}

