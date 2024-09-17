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

	Scan struct {
		Netmask  string `json:"netmask"`
		Interval int    `json:"interval"`
		Timeout  int    `json:"timeout"`
	} `json:"scan"`

	path string
}

func defaultConfig() *Config {
	return &Config{
		Port: 5101,
		Scan: struct {
			Netmask  string `json:"netmask"`
			Interval int    `json:"interval"`
			Timeout  int    `json:"timeout"`
		}{
			Netmask:  "255.255.255.0",
			Interval: 120,
			Timeout:  1,
		},
	}
}

func NewConfig(path string) (*Config, error) {
	config := new(Config)
	config.path = path

	file, err := os.ReadFile(path)
	if err == nil {
		if err = json.Unmarshal(file, config); err != nil {
			return nil, err
		}

		return config, nil
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	config = defaultConfig()
	config.path = path
	newConfig, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Println("C")
		return nil, err
	}
	if err = os.WriteFile(path, newConfig, 0600); err != nil {
		return nil, err
	}

	return config, nil
}
