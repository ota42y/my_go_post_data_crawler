package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mattn/go-shellwords"

	"../command"
)

type Setting struct {
	Shells []ShellCommand
}

type ShellCommand struct {
	Name    string
	Command string
	WorkDir string
}

type Shell struct {
	commands map[string]ShellCommand
}

func New(settingText string) (c *Shell) {
	s := Setting{}
	_, err := toml.Decode(settingText, &s)
	if err != nil {
		return nil
	}

	shell := &Shell{}
	shell.commands = make(map[string]ShellCommand)

	for _, commandData := range s.Shells {
		shell.commands[commandData.Name] = commandData
	}

	return shell
}

func (c *Shell) IsExecute(order command.Order) bool {
	return order.Name == "shell"
}

func (c *Shell) Execute(order command.Order) string {
	inputs := strings.Split(order.Data, " ")
	if 0 < len(inputs) {
		name := inputs[0]
		cmd, ok := c.commands[name]
		if ok {
			prev, err := filepath.Abs(".")
			if err != nil {
				return "prev directory get error"
			}
			defer os.Chdir(prev)
			os.Chdir(cmd.WorkDir)

			var line string
			if 1 < len(inputs) {
				line = fmt.Sprintf(cmd.Command, strings.Join(inputs[1:], " "))
			} else {
				line = cmd.Command
			}

			args, err := shellwords.Parse(line)
			if err != nil {
				return err.Error()
			}

			out, err := exec.Command(args[0], args[1:]...).Output()
			if err != nil {
				return err.Error()
			}
			return strings.Trim(string(out), "\n")
		}
	}

	return "no shell command"
}
