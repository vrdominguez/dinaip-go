/*
Package config provides functionality for reading and working with configuration files.

The package includes a Config struct that represents the configuration settings, as well as a function for reading the configuration from a YAML file.
*/
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration settings.
type Config struct {
	Username string              `yaml:"username"`
	Password string              `yaml:"password"`
	Logs     LogConfig           `yaml:"logs"`
	Zones    map[string][]string `yaml:"zones"`
}

// LogConfig represents the log configuration settings.
type LogConfig struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

// ReadConfig reads the configuration from a YAML file.
// It takes the filename as input and returns a pointer to the Config struct and an error if any.
func ReadConfig(filename string) (*Config, error) {
	// Read YAML file
	data, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML data into Config struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
