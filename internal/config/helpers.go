package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func getConfigFilePath() (string, error) {
	// Get home dir of OS
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("An error ocurred while getting home directory: %v\n", err)
		return "", err
	}
	
	// Return config file path
	return homeDir + "/" + configFileName, nil
}

func write (config Config) error {
	byteData, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("An error ocurred while marshalling data: %v\n", err)
		return err
	}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("An error ocurred while getting config file path: %v\n", err)
		return err
	}
	if err := os.WriteFile(configFilePath, byteData, 0666); err != nil {
		fmt.Printf("An error ocurred while writing config data: %v\n",  err)
		return err
	}
	return nil
}
