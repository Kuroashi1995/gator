package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/Kuroashi1995/rss-go/internal/commands"
	"github.com/Kuroashi1995/rss-go/internal/config"
	"github.com/Kuroashi1995/rss-go/internal/database"
	"github.com/Kuroashi1995/rss-go/internal/state"
	_ "github.com/lib/pq"
)
func middlewareLoggedIn(handler func(s *state.State, cmd commands.Command, user database.User) error) func(*state.State, commands.Command) error {
	return func (s *state.State, cmd commands.Command) error {
		currentUser, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			fmt.Printf("error while logging user: %v\n", err.Error())
			return nil
		}
		return handler(s, cmd, currentUser)
	}
}

func main() {
	// Get config
	gatorConfig, err := config.Read()
	if err != nil {
		return
	}
	state := state.State{
		Config : &gatorConfig,
	}

	// store state db connection
	db, err := sql.Open("postgres", gatorConfig.DBUrl)
	if err != nil {
		fmt.Println("Error connecting to the db: ", err.Error())
	}
	dbQueries := database.New(db)
	state.Db = dbQueries

	//initialize commands
	cliCommands := commands.InitializeCommands()
	cliCommands.Register("login", commands.HandlerLogin)
	cliCommands.Register("register", commands.HandlerRegister)
	cliCommands.Register("reset", commands.HandlerReset)
	cliCommands.Register("users", commands.HandlerUsers)
	cliCommands.Register("agg", commands.HandlerAgg)
	cliCommands.Register("addfeed", middlewareLoggedIn(commands.HandlerAddFeed))
	cliCommands.Register("feeds", commands.HandlerFeeds)
	cliCommands.Register("follow", middlewareLoggedIn(commands.HandlerFollow))
	cliCommands.Register("following", commands.HandlerFollowing)
	cliCommands.Register("unfollow", middlewareLoggedIn(commands.HandlerUnfollow))


	//check arguments
	if len(os.Args) < 2 {
		fmt.Println("Too few arguments")
		os.Exit(1)
	}
	commandName := os.Args[1]
	var cliArguments []string
	if len(os.Args) > 2 {
		cliArguments = os.Args[2:]
	}
	if err := cliCommands.Run(&state, commands.Command{
		Name: commandName,
		Arguments: cliArguments,
	}); err != nil {
		fmt.Printf("An error ocurred running command: %v\n", err.Error())
		os.Exit(1)
	}
}
