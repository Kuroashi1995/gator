package state

import (
	"github.com/Kuroashi1995/rss-go/internal/config"
	"github.com/Kuroashi1995/rss-go/internal/database"
)


type State struct {
	Config			*config.Config
	Db				*database.Queries
}
