# RSS CLI TOOL
This tool is meant to subscribe to a bunch of RSS feeds and explore them via CLI

## Requisites:
- Go
- Postgres

## Installation:
To install the tool, run `go install www.github.com/Kuroashi1995/rss-go`
Then, set up a config file in your home directory, as follows:
.gatorconfig.json:
`
{
    "db_url": "example://user:password@host:port/gator?sslmode=disable",
    "current_user_name": "your username"
}
`
## Commands:
- `gator login username`            - Log-in to a different user
- `gator register username`         - Registers an user
- `gator users`                     - Logs all users
- `gator agg frequency`             - Fetches posts from feeds every {frequency}
- `gator addfeed "name" "url"`      - Add feed source with {name} and {url}
- `gator feeds`                     - Logs feeds sources
- `gator follow "url"`              - Follows the given {url} feed
- `gator following`                 - Logs all the feeds the current user follows
- `gator unfollow "url"`            - Unfollows the given {url} feed
- `gator browse [limit]`            - Logs the latest [limit defaults to 2] posts
