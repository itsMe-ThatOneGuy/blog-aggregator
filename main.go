package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/commands"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/config"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/database"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/handlers"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/state"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DB)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	s := &state.State{
		DB:            dbQueries,
		ConfigPointer: &cfg,
	}

	cmds := commands.Commands{
		Handlers: make(map[string]func(*state.State, commands.Command) error),
	}

	cmds.Register("login", handlers.HandlerLogin)
	cmds.Register("register", handlers.HandlerRegister)
	cmds.Register("reset", handlers.Reset)
	cmds.Register("users", handlers.HandlerGetUsers)
	cmds.Register("agg", handlers.Agg)
	cmds.Register("addfeed", handlers.AddFeed)

	if len(os.Args) < 2 {
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
