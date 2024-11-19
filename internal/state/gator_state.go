package state

import (
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/config"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/database"
)

type State struct {
	DB            *database.Queries
	ConfigPointer *config.Config
}
