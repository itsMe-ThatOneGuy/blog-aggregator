package handlers

import (
	"errors"
	"fmt"

	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/commands"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/state"
)

func HandlerLogin(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("login handler expects a username argument")
	}
	username := cmd.Args[0]

	s.ConfigPointer.SetUser(username)
	fmt.Printf("User has been set: %s\n", username)

	return nil
}
