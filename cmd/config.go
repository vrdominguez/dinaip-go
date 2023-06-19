package cmd

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/vrdominguez/dinaip-go/config"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate configuration for dinaip-go",
	Long: `Generate configuration for dinaip-go

It will ask for all the necesary data and generate (or overwrite) the indicated
configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := os.Stat(configFile)

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Warnf("Log file (%s) does not exists", configFile)
			} else {
				log.Errorf("Cannot determine if config file (%s) already exists: %s", configFile, err)
			}
		} else {
			if err := initializeLogger(configFile); err != nil {
				log.Errorf("An invalid file already exists at %s, it will be overwritten.", configFile)
			}
		}

		fmt.Println("Enter data for the config file:")

		var config config.Config

		fmt.Print("Enter username: ")
		if _, err := fmt.Scan(&config.Username); err != nil {
			return fmt.Errorf("cannot read username: %s", err)
		}

		fmt.Print("Enter password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("cannot read password: %s", err)
		}
		config.Password = string(password)

		fmt.Print("\nEnter log level: ")
		if _, err := fmt.Scan(&config.Logs.Level); err != nil {
			return fmt.Errorf("cannot read log level: %s", err)
		}

		fmt.Print("Enter log path: ")
		if _, err := fmt.Scan(&config.Logs.Path); err != nil {
			return fmt.Errorf("cannot tead log path: %s", err)
		}

		config.Zones = make(map[string][]string)

		for {
			var domain string
			var subdomains []string

			fmt.Print("Enter domain (or 'done' to finish): ")
			if _, err := fmt.Scan(&domain); err != nil {
				return fmt.Errorf("cannot read domain: %e", err)
			}

			if domain == "done" {
				break
			}

			for {
				var subdomain string

				fmt.Print("Enter subdomain (or 'done' to finish): ")
				if _, err := fmt.Scan(&subdomain); err != nil {
					return fmt.Errorf("cannot read subdomain: %e", err)
				}

				if subdomain == "done" {
					break
				}

				subdomains = append(subdomains, subdomain)
			}

			config.Zones[domain] = subdomains
		}

		data, err := yaml.Marshal(&config)
		if err != nil {
			return fmt.Errorf("cannot generate yaml for config: %s", err)
		}

		err = os.WriteFile(configFile, data, 0600)
		if err != nil {
			return fmt.Errorf("cannot write config file at %s: %s", configFile, err)
		}

		fmt.Printf("Config file written at %s\n", configFile)

		// Reload log config to write the log of config creation according to the new config
		if err := initializeLogger(configFile); err != nil {
			return err
		}

		log.Infof("Generated new configuration file at %s", configFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
