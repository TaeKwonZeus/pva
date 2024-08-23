package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
)

type Config struct {
	Port int `json:"port"`

	path string
}

func defaultConfig(path string) *Config {
	return &Config{Port: 5101, path: path}
}

func NewConfig(path string) (*Config, error) {
	config := new(Config)
	config.path = path

	file, err := os.ReadFile(path)
	if err == nil {
		if err = json.Unmarshal(file, config); err != nil {
			log.Println("A")
			return nil, err
		}

		return config, nil
	}

	if !errors.Is(err, fs.ErrNotExist) {
		log.Println("B")
		return nil, err
	}

	newConfig, err := json.MarshalIndent(defaultConfig(path), "", "    ")
	if err != nil {
		log.Println("C")
		return nil, err
	}

	if err = os.WriteFile(path, newConfig, 0600); err != nil {
		log.Println("D")
		return nil, err
	}

	return config, nil
}

func (c *Config) Save() error {
	newC, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, newC, 0700)
}
