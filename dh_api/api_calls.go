/*
Package dh_api provides a client for interacting with the Dinahosting API.

The package includes types and functions for making requests to the Dinahosting API and handling the responses.
*/
package dh_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

/*
GetTypeIpRecords retrieves IP records of a specific type for a domain.

It takes the domain name and record type as input and returns a DomainGetZonesResponse and an error if any.

Example usage:

	ipRecords, err := client.GetTypeIpRecords("example.com", dh_api.TypeA)
	if err != nil {
		// Handle error
	}

	// Access the IP records...
*/
func (c *DinahostingClient) GetTypeIpRecords(domain string, rtype RecordType) (DomainGetZonesResponse, error) {
	var apiResponse DomainGetZonesResponse
	params := url.Values{}

	switch rtype {
	case TypeA:
		params.Set("command", "Domain_Zone_GetTypeA")
	case TypeAAAA:
		params.Set("command", "Domain_Zone_GetTypeAAAA")
	default:
		return apiResponse, fmt.Errorf("invalid type: %s", rtype)

	}

	params.Set("domain", domain)

	resp, err := c.DoRequest(params)
	if err != nil {
		return apiResponse, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiResponse, err
	}

	if err = json.Unmarshal(body, &apiResponse); err != nil {
		return apiResponse, err
	}

	responseCode := apiResponse.ResponseCode

	if responseCode == 1000 {
		return apiResponse, nil
	}

	errApi := fmt.Errorf("API error %d, %s", responseCode, apiResponse.Errors[0].(map[string]interface{})["message"])
	return apiResponse, errApi
}

/*
SetTypeIpRecord sets an IP record for a domain.

It takes the domain name, hostname, IP address, old IP address (optional), and record type as input.
It returns an error if any.

Example usage:

	err := client.SetTypeIpRecord("example.com", "www", "192.168.0.1", "", dh_api.TypeA)
	if err != nil {
		// Handle error
	}

	// ...
*/
func (c *DinahostingClient) SetTypeIpRecord(domain, hostname, ip, oldIp string, rtype RecordType) error {
	var apiResponse GenericResponse
	params := url.Values{}
	params.Set("domain", domain)
	params.Set("hostname", hostname)
	params.Set("ip", ip)

	if oldIp != "" {
		params.Set("oldip", ip)

		switch rtype {
		case TypeA:
			params.Set("command", "Domain_Zone_UpdateTypeA")
		case TypeAAAA:
			params.Set("command", "Domain_Zone_UpdateTypeAAAA")
		default:
			return fmt.Errorf("invalid type: %s", rtype)
		}

	} else {
		switch rtype {
		case TypeA:
			params.Set("command", "Domain_Zone_AddTypeA")
		case TypeAAAA:
			params.Set("command", "Domain_Zone_AddTypeAAAA")
		default:
			return fmt.Errorf("invalid type: %s", rtype)
		}
	}

	resp, err := c.DoRequest(params)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}

	responseCode := apiResponse.ResponseCode

	if responseCode == 1000 {
		return nil
	}

	return fmt.Errorf("API error %d, %s", responseCode, apiResponse.Errors[0].(map[string]interface{})["message"])
}
