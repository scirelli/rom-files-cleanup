package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/scirelli/roms-files-cleanup/internal/pkg/log"
)

// AppConfig configuration data for entire application.
type AppConfig struct {
	Debug          bool     `json:"debug"`
	LogLevel       string   `json:"logLevel"`
	ROMDatabaseURI []string `json:"romDatabaseURI"`
}

type option func(*AppConfig)

func (ac *AppConfig) Option(opts ...option) {
	for _, opt := range opts {
		opt(ac)
	}
}

func WithLogLevel(ll string) option {
	return func(ac *AppConfig) {
		ac.LogLevel = ll
	}
}

func WithDebug(db bool) option {
	return func(ac *AppConfig) option {
		ac.Debug = db
	}

}

func NewConfig() *AppConfig {
	var config AppConfig

	return &config
}

// LoadConfig a config file.
func LoadConfig(fileName string) (*AppConfig, error) {
	var config *AppConfig = NewConfig()

	jsonFile, err := os.Open(fileName)
	if err != nil {
		return &config, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return config, err
	}

	json.Unmarshal(byteValue, config)

	if config.LogLevel == "" {
		config.WithLogLevel(log.DEFAULT_LOG_LEVEL)
	} else {
		config.WithLogLevel(log.GetLevel(config.LogLevel))
	}

	if len(config.ROMDatabaseFiles) <= 0 {
		return config, errors.New("Must have at least one ROM database file.")
	}

	return config, nil
}
