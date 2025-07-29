package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	WeatherAPI struct {
		Key string `json:"key"`
	} `json:"weatherAPI"`
}

func (c *Config) String() string {
	return fmt.Sprintf("conf loaded key size: '%d'", len(c.WeatherAPI.Key))
}

func LoadConfig(fileLocation string) (*Config, error) {
	f1, err := os.Open(fileLocation)
	if err != nil {
		return nil, err
	}
	return LoadConfigFile(f1)
}

func LoadConfigFile(f1 *os.File) (*Config, error) {
	conf := &Config{}
	err := json.NewDecoder(f1).Decode(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
