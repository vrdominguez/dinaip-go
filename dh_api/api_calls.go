package dh_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

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
