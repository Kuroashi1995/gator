package config

import (
	"errors"
	"fmt"
)


type Config struct {
	DBUrl				string		`json:"db_url"`
	CurrentUserName		string		`json:"current_user_name"`
}


const configFileName = ".gatorconfig.json"


func (c *Config) SetUser(username string) error {
	if username == "" {
		fmt.Printf("User name can't be empty")
		return errors.New("Username empty")
	}
	c.CurrentUserName = username
	if err := write(*c); err != nil {
		return err
	}
	return nil
}
