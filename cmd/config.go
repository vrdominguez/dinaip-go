package cmd

import (
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
		if err := initializeLogger(configFile); err != nil {
			return err
		}

		fmt.Println("Enter data for the config file:")

		var config config.Config

		fmt.Print("Enter username: ")
		fmt.Scan(&config.Username)

		fmt.Print("Enter password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("cannot read password: %e", err)
		}
		config.Password = string(password)

		fmt.Print("\nEnter log level: ")
		fmt.Scan(&config.Logs.Level)

		fmt.Print("Enter log path: ")
		fmt.Scan(&config.Logs.Path)

		config.Zones = make(map[string][]string)

		for {
			var domain string
			var subdomains []string

			fmt.Print("Enter domain (or 'done' to finish): ")
			fmt.Scan(&domain)

			if domain == "done" {
				break
			}

			for {
				var subdomain string

				fmt.Print("Enter subdomain (or 'done' to finish): ")
				fmt.Scan(&subdomain)

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

		log.Infof("Generated new configuration file at %s", configFile)
		fmt.Printf("Config file written at %s\n", configFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
