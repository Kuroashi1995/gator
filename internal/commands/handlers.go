package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Kuroashi1995/rss-go/internal/database"
	"github.com/Kuroashi1995/rss-go/internal/rss"
	"github.com/Kuroashi1995/rss-go/internal/state"
	"github.com/google/uuid"
)

func HandlerLogin(s *state.State, cmd Command) error {
	//check arguments
	if len(cmd.Arguments) == 0 {
		return errors.New("Login command expects an username")
	}
	//check if user exists in database
	ctx := context.Background()
	_, err := s.Db.GetUser(ctx, cmd.Arguments[0])	
	if err != nil {
		fmt.Println("User not found: ", err.Error())
		os.Exit(1)
	}

	s.Config.SetUser(cmd.Arguments[0])
	fmt.Printf("Username set!\n")
	return nil
}

func HandlerRegister(s *state.State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("Register command needs at least one argument")
	}
	ctx := context.Background()
	user, err := s.Db.CreateUser(ctx, database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Arguments[0],
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	s.Config.SetUser(user.Name)
	fmt.Println("User created!")
	fmt.Println(user)
	return nil
}

func HandlerReset(s *state.State, cmd Command) error {
	ctx := context.Background()
	if err := s.Db.ResetDB(ctx); err != nil {
		return fmt.Errorf("Error resetting database: %v", err.Error())
	}
	fmt.Println("User database reset")
	return nil
}

func HandlerUsers(s *state.State, cmd Command) error {
	ctx := context.Background()
	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("Error getting users: %v", err.Error())
	}
	for _, user := range users {
		
		if user.Name == s.Config.CurrentUserName{
			fmt.Println("*", user.Name, "(current)")
		} else {
			fmt.Println("*", user.Name)
		}

	}
	return nil
}

func HandlerAgg(s *state.State, cmd Command) error {
	ctx := context.Background()
	result, err := rss.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Error fecthing rrs feed: %v\n", err.Error())
	}
	fmt.Println(result)
	return nil
}

func HandlerAddFeed(s *state.State, cmd Command, user database.User) error {
	//check args for failure
	if len(cmd.Arguments) < 2 {
		fmt.Println("addfeed expects 2 arguments")
		os.Exit(1)
	}
	// get parent context
	ctx := context.Background()
	
	//create the feed
	createdFeed, err := s.Db.CreateFeed(ctx, database.CreateFeedParams{
		ID: uuid.New(),
		Name: cmd.Arguments[0],
		UserID: user.ID,
		Url: cmd.Arguments[1],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("An error ocurred while creating feed: %v\n", err.Error())
	}

	_, err = s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: createdFeed.ID,
	})

	if err != nil {
		return fmt.Errorf("an error ocurred while creating feed_follow: %v\n", err.Error())
	}
	fmt.Println(createdFeed)
	return nil

}

func HandlerFeeds(s *state.State, cmd Command) error {
	//set context for db calls
	ctx := context.Background()

	//get the feeds
	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("an error ocurres while getting feeds: %v\n", err.Error())
	}
	
	//show feeds on stdout
	for _, feed := range feeds {
		fmt.Printf("- %v:\n\turl: %v\n\tcreated by: %v\n", feed.Name, feed.Url, feed.UserName)
	}

	return nil
}

func HandlerFollow(s *state.State, cmd Command, currentUser database.User) error {
	// set context for db call
	ctx := context.Background()
	
	//get feed if it exists already
	feed, err := s.Db.GetFeedByUrl(ctx, cmd.Arguments[0])
	if err != nil {
		fmt.Printf("Error getting feed from db: %v\n", err.Error())
	}
	// if feed does not exist, create it
	if feed == (database.Feed{}) {
		feed, err = s.Db.CreateFeed(ctx, database.CreateFeedParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID: currentUser.ID,
			Url: cmd.Arguments[0],
		})
		if err != nil {
			return fmt.Errorf("An error ocurred creating the feed: %v\n", err.Error())
		}
	}

	//create feed_follow
	feedFollow, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: currentUser.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("an error ocurred while creating feed follow: %v\n", err.Error())
	}

	//print succesful result
	fmt.Printf("Created: %v, with user: %v\n", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func HandlerFollowing(s *state.State, cmd Command) error {
	//get context for db calls
	ctx := context.Background()

	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("an error ocurred while getting feeds for user: %v\n", err.Error())
	}
	fmt.Println("The current user follows:")
	for _, feedFollow := range feedFollows {
		fmt.Printf("- %v\n", feedFollow.FeedName)
	}
	return nil
}

func HandlerUnfollow(s *state.State, cmd Command, user database.User) error {
	ctx := context.Background()
	// get feed
	feed, err := s.Db.GetFeedByUrl(ctx, cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("an error ocurred while getting feed: %v\n", err.Error())
	}
	if err := s.Db.DeleteFeedFollowByUserFeed(context.Background(), database.DeleteFeedFollowByUserFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("an error occurred while deleting feed follow: %v\n", err.Error())
	}
	return nil

}
