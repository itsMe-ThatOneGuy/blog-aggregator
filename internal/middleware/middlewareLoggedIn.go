package middleware

import (
	"context"

	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/commands"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/database"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/state"
)

func MiddlewareLoggedIn(handler func(s *state.State, cmd commands.Command, user database.User) error) func(*state.State, commands.Command) error {
	return func(s *state.State, cmd commands.Command) error {
		user, err := s.DB.GetUser(context.Background(), s.ConfigPointer.User)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
