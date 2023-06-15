package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/vrdominguez/dinaip-go/config"
)

func initializeLogger(configFile string) error {
	if config, err := config.ReadConfig(configFile); err != nil {
		return fmt.Errorf("cannot read config file at %s: %s", configFile, err)
	} else {
		configuration = config
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Update log level with config value
	if logLevel, err := convertLogLevel(configuration.Logs.Level); err != nil {
		return fmt.Errorf("cannot set log level: %s", err)
	} else {
		log.SetLevel(logLevel)
	}

	// Set log output to file (if not STDOUT as file)
	if configuration.Logs.Path != "STDOUT" {
		if file, err := os.OpenFile(configuration.Logs.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err != nil {
			return fmt.Errorf("cannot create logs file at %s: %s", configuration.Logs.Path, err)
		} else {
			log.SetOutput(file)
		}
	}

	return nil
}

func convertLogLevel(logLevel string) (log.Level, error) {
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		return log.DebugLevel, nil
	case "INFO":
		return log.InfoLevel, nil
	case "WARN", "WARNING":
		return log.WarnLevel, nil
	case "ERROR":
		return log.ErrorLevel, nil
	case "FATAL":
		return log.FatalLevel, nil
	case "PANIC":
		return log.PanicLevel, nil
	default:
		return log.InfoLevel, fmt.Errorf("invalid log level: %s", logLevel)
	}
}

func convertToLogrusFields(fields map[string]string) log.Fields {
	logFields := log.Fields{}

	for key, value := range fields {
		logFields[key] = value
	}

	return logFields
}
