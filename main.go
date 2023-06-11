package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/vrdominguez/dinaip-go/config"
	"github.com/vrdominguez/dinaip-go/dh_api"
	"github.com/vrdominguez/dinaip-go/iputils"
)

const (
	defaultConfigPath = "/etc/dianaIP.yaml"
)

var configuration *config.Config
var client *dh_api.DinahostingClient

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Before config set logs to info via STDOUT
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// Load config
	var configFile string
	flag.StringVar(&configFile, "c", "", "path to config file")
	flag.Parse()

	if configFile == "" {
		configFile = defaultConfigPath
	}

	if config, err := config.ReadConfig(configFile); err != nil {
		log.Panicf("Cannot read config file at %s", configFile)
	} else {
		configuration = config
	}

	// Set file log
	if file, err := os.OpenFile(configuration.Logs.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		log.Panicf("Cannot create logs file at %s: %s", configuration.Logs.Path, err)
	} else {
		log.Infof("Logs changed to file: %s", configuration.Logs.Path)
		log.SetOutput(file)
	}

	// Update loglevel with config value
	if logLevel, err := convertLogLevel(configuration.Logs.Level); err != nil {
		log.Panicf("Cannot set log level: %s", err)
	} else {
		log.SetLevel(logLevel)
	}

	client = dh_api.NewClient(configuration.Username, configuration.Password)
}

// main is a proof of concept where various functions are tested until sufficient libraries are implemented.
func main() {

	for {
		currentIp, err := iputils.GetPublicIp()
		if err != nil {
			log.Errorf("Cannot get current IP: %s", err)
		} else {
			for domain, subdomains := range configuration.Zones {
				// Get domain records
				recordsResponse, err := client.GetTypeIpRecords(domain, dh_api.TypeA)
				if err != nil {
					log.Errorf("Cannot check domain %s: %s", domain, err)
					continue // Go for next domain if zones cannot be checked
				}

				for _, subdomain := range subdomains {
					logFields := convertToLogrusFields(map[string]string{
						"domain":  domain,
						subdomain: subdomain,
					})
					subdomainLogger := log.WithFields(logFields)

					subdomainLogger.Debugf("Working with %s.%s", subdomain, domain)

					subdomainIp := ""
					for _, dnsRecord := range recordsResponse.Data {
						if dnsRecord.Hostname == subdomain {
							subdomainIp = dnsRecord.Ip
							break
						}
					}

					subdomainLogger.Debugf("Current IP %s", subdomainIp)
					if subdomainIp == "" {
						subdomainLogger.Infof("Creating new dns record")
						if err := client.SetTypeIpRecord(domain, subdomain, currentIp, "", dh_api.TypeA); err != nil {
							subdomainLogger.Errorf("Cannot create zone: %s", err)
						} else {
							subdomainLogger.Info("Create new zone")
							subdomainLogger.Debugf("New ip: %s", currentIp)
						}

					} else if currentIp != subdomainIp {
						subdomainLogger.Infof("Updating dns record")
						if err := client.SetTypeIpRecord(domain, subdomain, currentIp, subdomainIp, dh_api.TypeA); err != nil {
							subdomainLogger.Errorf("Cannot update zone record: %s", err)
						} else {
							subdomainLogger.Info("Updated zone recod")
							subdomainLogger.Debugf("New ip: %s", currentIp)
						}

					} else {
						subdomainLogger.Info("No updates needed")
					}
				}
			}
		}

		// Wait 5 minutes...
		log.Debug("Waiting for 5 minutes until next review")
		time.Sleep(5 * time.Minute)
	}
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
