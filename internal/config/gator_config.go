package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DB   string `json:"db_url"`
	User string `json:"current_user_name"`
}

func TestDir() {
	test, _ := getConfigFilePath()

	fmt.Println(test)
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	config := Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, errors.New("issue when decoding json file")
	}

	return config, nil
}

func (cfg *Config) SetUser(name string) error {
	cfg.User = name
	return write(*cfg)
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return errors.New("issue opening file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return errors.New("issue writing to file")
	}

	return nil
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("Issue reading file")
	}

	fullPath := dir + "/.gatorconfig.json"

	return fullPath, nil
}
