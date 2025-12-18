package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)


type RSSFeed struct {
	Channel			struct {
		Title		string		`xml:"title"`
		Link		string		`xml:"link"`
		Description	string		`xml:"description"`
		Item		[]RSSItem	`xml:"item"`
	}							`xml:"channel"`
}

type RSSItem struct {
	Title			string		`xml:"title"`
	Link			string		`xml:"link"`
	Description		string		`xml:"description"`
	PubDate			string		`xml:"pubDate"`
}


func FetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	// configuring request
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v\n", err.Error())
	}
	req.Header.Set("User-Agent", "gator")

	// making request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %v\n", err.Error())
	}
	defer res.Body.Close()
	
	// process body data
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading body: %v", err.Error())
	}
	var rssData RSSFeed
	if err := xml.Unmarshal(data, &rssData); err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %v", err.Error())
	}

	// clean data
	rssData.Channel.Title = html.UnescapeString(rssData.Channel.Title)
	rssData.Channel.Description = html.UnescapeString(rssData.Channel.Description)
	for i, rssItem := range rssData.Channel.Item {
		rssData.Channel.Item[i].Title = html.UnescapeString(rssItem.Title)
		rssData.Channel.Item[i].Description = html.UnescapeString(rssItem.Description)
	}

	return &rssData, nil
}
