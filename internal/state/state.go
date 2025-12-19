package state

import (
	"github.com/Kuroashi1995/gator/internal/config"
	"github.com/Kuroashi1995/gator/internal/database"
)


type State struct {
	Config			*config.Config
	Db				*database.Queries
}
