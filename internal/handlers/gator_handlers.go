package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/commands"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/database"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/state"
)

func HandlerLogin(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login handler expects a username argument")
	}
	username := cmd.Args[0]

	s.ConfigPointer.SetUser(username)
	fmt.Printf("User has been set: %s\n", username)
func HandlerRegister(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("register handler expects a username")
	}
	name := cmd.Args[0]

	_, err := s.DB.GetUser(context.Background(), name)
	if err == nil {
		fmt.Printf("user '%s' already exists\n", name)
		os.Exit(1)
	} else if err != sql.ErrNoRows {
		return err
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	newuser, err := s.DB.CreateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("Issue creating user")
	}
	s.ConfigPointer.SetUser(name)

	fmt.Println("User has been created")
	printUser(newuser)

	return nil
}
