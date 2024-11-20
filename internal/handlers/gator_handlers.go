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
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/rss"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/state"
)

func Agg(s *state.State, cmd commands.Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%+v", feed)

	return nil
}

func Reset(s *state.State, cmd commands.Command) error {
	err := s.DB.ResetDB(context.Background())
	if err != nil {
		return fmt.Errorf("issue deleting users: %w", err)
	}
	fmt.Println("Database reset")

	return nil
}

func HandlerLogin(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login handler expects a username argument")
	}
	name := cmd.Args[0]

	_, err := s.DB.GetUser(context.Background(), name)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("can't login user: %w\n", err)
		}

		return err
	}

	s.ConfigPointer.SetUser(name)
	fmt.Printf("User has been set: %s\n", name)

	return nil
}

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

func HandlerGetUsers(s *state.State, cmd commands.Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("issue retrieving users")
	}

	for _, user := range users {
		current := s.ConfigPointer.User
		if user.Name == current {
			fmt.Printf("* %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %v\n", user.Name)
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf("* ID:   %v\n", user.ID)
	fmt.Printf("* Name: %v\n", user.Name)
}
