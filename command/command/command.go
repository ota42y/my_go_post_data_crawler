package command

type Command interface {
	// command string equal it's command name
	// if return true, call Execute
	IsExecute(command string) bool

	// execute command
	Execute(data string) string
}
