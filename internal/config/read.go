package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read() (Config, error) {
	// Read config file from homedir
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("An error ocurred while reading config file: %v\n", err)
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		fmt.Printf("An error ocurred while unmarshalling data: %v\n", err)
		return Config{}, err
	}
	return config, nil
}
