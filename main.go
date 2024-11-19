package main

import (
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
		log.Fatalf("error reading config: %v", err)
	}

	s := &state.State{
		ConfigPointer: &cfg,
	}

	cmds := commands.Commands{
		Handlers: make(map[string]func(*state.State, commands.Command) error),
	}

	cmds.Register("login", handlers.HandlerLogin)

	if len(os.Args) < 3 {
		fmt.Println("Provide a command name and arguments")
		return
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	err = cmds.Run(s, commands.Command{Name: commandName, Args: commandArgs})
	if err != nil {
		log.Fatal(err)
	}
}
