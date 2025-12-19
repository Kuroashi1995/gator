package rss

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Kuroashi1995/gator/internal/database"
	"github.com/Kuroashi1995/gator/internal/state"
	"github.com/google/uuid"
)

func ScrapeFeeds(s *state.State) error {
	// get context for queries
	ctx := context.Background()
		// get the next feed
		feed, err := s.Db.GetNextFeedToFetch(ctx, s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("an error ocurred while fetching feed: %v\n", err.Error())
		}
		
		// mark feed as fetched
		if err := s.Db.MarkFeedFetched(ctx, feed.ID); err != nil {
			return fmt.Errorf("an error ocurred while setting as fetched: %v\n", err.Error())
		}

		//fetch the feed
		fmt.Printf("Fetching feed from: %v\n", feed.Url)
		fetchedData, err := FetchFeed(ctx, feed.Url)
		if err != nil {
			return fmt.Errorf("an error ocurred while fetching the feed: %v\n", err.Error())
		}

		for _, item := range fetchedData.Channel.Item {
			pubTime, err := time.Parse(time.Layout, item.PubDate)
			if err != nil {
				fmt.Printf("error parsing pub date: %v\n", err.Error())
			}
			_, err = s.Db.CreatePost(ctx, database.CreatePostParams{
				ID: uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Title: item.Title,
				Url: feed.Url,
				Description: sql.NullString{
					String: item.Description,
					Valid: true,
				},
				PublishedAt: sql.NullTime{
					Time: pubTime,
					Valid: true,
				},
				FeedID: feed.ID,
			})
			if err != nil {
				fmt.Printf("error creating post: %v\n", err.Error())
			}
		}
	return nil
}
