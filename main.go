package main

import (
	"fmt"

	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Config: %v\n", cfg)

	err = cfg.SetUser("Matthew")

	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Config: %v\n", cfg)
}
