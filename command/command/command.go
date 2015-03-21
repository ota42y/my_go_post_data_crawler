package command

type Command interface {
	// command string equal it's command name
	// if return true, call Execute
	IsExecute(order Order) bool

	// execute command
	Execute(order Order) string
}

type Order struct {
	Name string
	Data string
	Room string
	User string
}
