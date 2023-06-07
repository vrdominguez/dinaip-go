package main

import (
	"flag"
	"fmt"
	"time"

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
	// Load config
	var configFile string
	flag.StringVar(&configFile, "c", "", "path to config file")
	flag.Parse()

	if configFile == "" {
		configFile = defaultConfigPath
	}

	if config, err := config.ReadConfig(configFile); err != nil {
		panic(err)
	} else {
		configuration = config
	}

	// Set client
	client = dh_api.NewClient(configuration.Username, configuration.Password)
}

// main is a proof of concept where various functions are tested until sufficient libraries are implemented.
func main() {

	for {
		currentIp, err := iputils.GetPublicIp()
		if err != nil {
			panic(err)
		}

		for domain, subdomains := range configuration.Zones {
			// Get domain records
			recordsResponse, err := client.GetTypeIpRecords(domain, dh_api.TypeA)
			if err != nil {
				fmt.Printf("Cannot check domain %s: %s\n", domain, err)
				continue // Go for next domain if zones cannot be checked
			}

			for _, subdomain := range subdomains {
				subdomainIp := ""

				for _, dnsRecord := range recordsResponse.Data {
					if dnsRecord.Hostname == subdomain {
						subdomainIp = dnsRecord.Ip
						break
					}
				}

				if subdomainIp == "" {
					fmt.Printf("Domain: %s, subdomain: %s, Create new Zone\n", domain, subdomain)
					if err := client.SetTypeIpRecord(domain, subdomain, currentIp, "", dh_api.TypeA); err != nil {
						fmt.Printf("ERROR: Cannot create %s.%s with ip %s: %s\n", domain, subdomain, currentIp, err)
					} else {
						fmt.Printf("Updated %s.%s with ip %s\n", domain, subdomain, currentIp)
					}

				} else if currentIp != subdomainIp {
					fmt.Printf("Domain: %s, subdomain: %s, Update IP\n", domain, subdomain)
					if err := client.SetTypeIpRecord(domain, subdomain, currentIp, subdomainIp, dh_api.TypeA); err != nil {
						fmt.Printf("ERROR: Cannot update %s.%s from ip %s to %s: %s\n", domain, subdomain, subdomainIp, currentIp, err)
					} else {
						fmt.Printf("Updated %s.%s with ip %s\n", subdomain, domain, currentIp)
					}

				} else {
					fmt.Printf("%s.%s already has the ip %s\n", subdomain, domain, currentIp)
				}
			}
		}

		// Wait 5 minutes...
		time.Sleep(5 * time.Minute)
	}
}
