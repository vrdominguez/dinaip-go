package cmd

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/vrdominguez/dinaip-go/config"
	"github.com/vrdominguez/dinaip-go/dh_api"
	"github.com/vrdominguez/dinaip-go/iputils"
)

var configFile string
var configuration *config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dinaip-go",
	Short: "DDNS system for Dinahosting SL API",
	Long: `DDNS system for Dinahosting SL API.

Keeps multiple subdomains IP address updated.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initializeLogger(configFile); err != nil {
			return err
		}

		client := dh_api.NewClient(configuration.Username, configuration.Password)

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
							"domain":    domain,
							"subdomain": subdomain,
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
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "/etc/dinaip.yaml", "Path to the configuration YAML file.")
}
