package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Username string              `yaml:"username"`
	Password string              `yaml:"password"`
	Zones    map[string][]string `yaml:"zones"`
}

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
