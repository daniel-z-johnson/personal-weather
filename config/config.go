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
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	Database struct {
		Path string `json:"path"`
	} `json:"database"`
}

func (c *Config) String() string {
	return fmt.Sprintf("conf loaded key size: '%d', port: '%s', db: '%s'", len(c.WeatherAPI.Key), c.Server.Port, c.Database.Path)
}

func LoadConfig(fileLocation string) (*Config, error) {
	f1, err := os.Open(fileLocation)
	if err != nil {
		return nil, err
	}
	conf, err := LoadConfigFile(f1)
	if err != nil {
		return nil, err
	}
	
	// Override with environment variables if present
	if apiKey := os.Getenv("WEATHER_API_KEY"); apiKey != "" {
		conf.WeatherAPI.Key = apiKey
	}
	if port := os.Getenv("PORT"); port != "" {
		conf.Server.Port = port
	} else if conf.Server.Port == "" {
		conf.Server.Port = "1117" // default port
	}
	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		conf.Database.Path = dbPath
	} else if conf.Database.Path == "" {
		conf.Database.Path = "w.db" // default database path
	}
	
	return conf, nil
}

func LoadConfigFile(f1 *os.File) (*Config, error) {
	conf := &Config{}
	err := json.NewDecoder(f1).Decode(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
