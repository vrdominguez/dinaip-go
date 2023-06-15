package cmd

import (
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestInitializeLogger(t *testing.T) {
	// Get the absolute path of the current file
	absPath, _ := filepath.Abs("./" + t.Name())

	// Calculate the parent directory
	parentDir := filepath.Dir(filepath.Dir(absPath))

	configFile := filepath.Join(parentDir, "__fixtures__", "dinaIP.yaml")

	// Call the function being tested
	err := initializeLogger(configFile)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert logrus formatter
	assert.IsType(t, &log.JSONFormatter{}, log.StandardLogger().Formatter)

	// Assert log level
	assert.Equal(t, log.DebugLevel, log.StandardLogger().Level)

	// Assert log output file
	logOutput, ok := log.StandardLogger().Out.(*os.File)
	assert.True(t, ok)
	assert.Equal(t, "/tmp/test-dinaip.log", logOutput.Name())
}

func TestConvertLogLevel(t *testing.T) {
	// Test valid log levels
	validLevels := map[string]log.Level{
		"DEBUG":   log.DebugLevel,
		"INFO":    log.InfoLevel,
		"WARN":    log.WarnLevel,
		"WARNING": log.WarnLevel,
		"ERROR":   log.ErrorLevel,
		"FATAL":   log.FatalLevel,
		"PANIC":   log.PanicLevel,
	}

	for levelStr, expectedLevel := range validLevels {
		level, err := convertLogLevel(levelStr)
		assert.NoError(t, err)
		assert.Equal(t, expectedLevel, level)
	}

	// Test invalid log level
	level, err := convertLogLevel("INVALID")
	assert.Error(t, err)
	assert.Equal(t, log.InfoLevel, level)
}

func TestConvertToLogrusFields(t *testing.T) {
	fields := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	expectedFields := log.Fields{}
	expectedFields["key1"] = "value1"
	expectedFields["key2"] = "value2"

	logFields := convertToLogrusFields(fields)

	assert.Equal(t, expectedFields, logFields)
}
