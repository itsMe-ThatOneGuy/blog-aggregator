package commands

import (
	"errors"

	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*state.State, Command) error
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.Handlers[name] = f
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	value, ok := c.Handlers[cmd.Name]
	if !ok {
		return errors.New("command is not registered")
	}

	return value(s, cmd)
}
