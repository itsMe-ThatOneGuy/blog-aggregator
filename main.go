package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/commands"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/config"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/handlers"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/state"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	s := &state.State{
		ConfigPointer: &cfg,
	}

	cmds := commands.Commands{
		Handlers: make(map[string]func(*state.State, commands.Command) error),
	}

	cmds.Register("login", handlers.HandlerLogin)

	if len(os.Args) < 3 {
		err := errors.New("Provide a command name and arguments")
		fmt.Printf("%v\n", err)
		return
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	cmds.Run(s, commands.Command{Name: commandName, Args: commandArgs})
}
