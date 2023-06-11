package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	// Get directory of current file
	dir := filepath.Dir(t.Name())

	t.Run("Configuration file can be readed", func(t *testing.T) {
		// Unmarshal YAML data into Config struct
		config, err := ReadConfig(dir + "/../__fixtures__/dinaIP.yaml")
		assert.NoError(t, err)

		// Check if values are correct
		assert.Equal(t, "exampleUser", config.Username)
		assert.Equal(t, "examplePassword", config.Password)
		assert.Equal(t, []string{"www", "subdomain"}, config.Zones["example.com"])
		assert.Equal(t, []string{"zone"}, config.Zones["example.org"])
		assert.Equal(t, "/var/log/dinaip.log", config.Logs.Path)
		assert.Equal(t, "DEBUG", config.Logs.Level)
	})
}
