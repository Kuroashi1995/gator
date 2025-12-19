package commands

import (
	"errors"
	"fmt"

	"github.com/Kuroashi1995/gator/internal/state"
)

type Command struct {
	Name		string
	Arguments	[]string
}

type Commands struct {
	commandList		map[string]func(*state.State, Command) error
}

func InitializeCommands() *Commands {
	CommandList := make(map[string]func(*state.State, Command)error)
	newCommands := Commands{
		commandList: CommandList,
	}
	return &newCommands
}

func (c *Commands) Run (s *state.State, cmd Command) error {
	if s == nil {
		fmt.Printf("State is nil\n")
		return errors.New("State is nil")
	}
	_, ok := c.commandList[cmd.Name]
	if !ok {
		fmt.Printf("command entry doesn't exists in commandlist: %v\n", cmd.Name)
	}
	if err := c.commandList[cmd.Name](s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.commandList[name] = f
}
